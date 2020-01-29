package qredochain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/kv"
)

func (app *QredoChain) processTX(tx []byte, lightWeight bool) (uint32, []abcitypes.Event) {
	// 	The workhorse of the application - non-optional.
	// 	Execute the transaction in full.
	// 	ResponseDeliverTx.Code == 0 only if the transaction is fully valid.

	//Decode the Asset
	signedAsset, err := decodeTX(tx)
	if err != nil {
		return code.CodeTypeEncodingError, nil
	}
	//Make TX Hash
	txHashA := sha256.Sum256(tx)
	txHash := txHashA[:]

	fmt.Println("txHash = ", txHash[:])

	//Retrieve the Asset ID
	assetID := signedAsset.Asset.GetID()
	if assetID == nil {
		return code.CodeTypeEncodingError, nil
	}

	//Process the Transaction
	switch signedAsset.Asset.GetType() {
	case protobuffer.PBAssetType_wallet:
		wallet, err := assets.ReBuildWallet(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError, nil
		}
		code, events := app.processWallet(wallet, tx, txHash, lightWeight)
		return uint32(code), events
	case protobuffer.PBAssetType_iddoc:
		iddoc, err := assets.ReBuildIDDoc(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError, nil
		}
		code, events := app.processIDDoc(iddoc, tx, txHash, lightWeight)
		return uint32(code), events
	case protobuffer.PBAssetType_Group:
		group, err := assets.ReBuildGroup(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError, nil
		}
		code, events := app.processGroup(group, lightWeight)
		return uint32(code), events
	}
	return code.CodeTypeEncodingError, nil
}

func (app *QredoChain) processIDDoc(iddoc *assets.IDDoc, rawAsset []byte, txHash []byte, lightWeight bool) (TransactionCode, []abcitypes.Event) {
	if app.exists(txHash) {
		//Usually the tx cache takes care of this, but once its full, we need to stop duplicates of very old transactions
		dumpMessage(2, "Fail to add IDDoc - tx already in chain")
		return CodeAlreadyExists, nil
	}

	//IDDoc is immutable so if this AssetID already has a value we can't update it.
	if app.exists(iddoc.Key()) == true {
		dumpMessage(2, "Fail to add IDDoc - already exists")
		return CodeAlreadyExists, nil
	}

	//Check the IDDoc is valid
	if app.VerifyIDDoc(iddoc) == false {
		return CodeFailVerfication, nil
	}

	var events []types.Event

	//Add pointer from AssetID to the txHash of the Object
	if lightWeight == false {
		//Set the Tags

		err := app.Set(txHash, rawAsset)
		if err != nil {
			return CodeDatabaseFail, nil
		}

		events = []types.Event{
			{
				Type: "tag",
				Attributes: []kv.Pair{
					{Key: []byte("myname"), Value: []byte("chris")},
					//{Key: []byte("assetid"), Value: iddoc.Key()},
					//{Key: []byte("txid"), Value: []byte(hex.EncodeToString(txHash))}, //txid is hex string  but all tags need to be byte array
				},
			},
		}

		// err := app.Set(txHash, rawAsset)
		// if err != nil {
		// 	return CodeDatabaseFail, nil
		// }
		// err = app.Set(iddoc.Key(), txHash)
		// if err != nil {
		// 	return CodeDatabaseFail, nil
		// }

	}

	print("----Events---------------------------------------\n")
	print(events)
	print("---- End Events---------------------------------------\n")

	return CodeTypeOK, events
}

func (app *QredoChain) processWallet(wallet *assets.Wallet, rawAsset []byte, txHash []byte, lightWeight bool) (TransactionCode, []abcitypes.Event) {
	if app.exists(txHash) {
		//Usually the tx cache takes care of this, but once its full, we need to stop duplicates of very old transactions
		dumpMessage(2, "Fail to add wallet - tx already in chain\n")
		return CodeAlreadyExists, nil
	}

	currentIndex, err := app.Get(wallet.Key())
	if err != nil {
		return CodeDatabaseFail, nil
	}
	var newAssetIndexString string

	if currentIndex == nil {
		//New Wallet
		if app.VerifyNewWallet(wallet) == false {
			dumpMessage(4, "Wallet failed verification")
			return CodeFailVerfication, nil
		}
		newAssetIndexString = IndexFormater(0)

	} else {
		//Check we are correctly incrementing the index
		currentIndexInteger, err := strconv.ParseInt(string(currentIndex), 10, 64)
		if err != nil {
			dumpMessage(4, "Failed to Parse Current Index")
			return CodeFailVerfication, nil
		}
		newIndex := wallet.CurrentAsset.Asset.Index
		if newIndex != currentIndexInteger+1 {
			dumpMessage(2, "Invalid Wallet Index\n")
			return CodeFailVerfication, nil
		}
		newAssetIndexString = IndexFormater(newIndex)
		//Wallet update
		if app.VerifyWalletUpdate(wallet) == false {
			dumpMessage(4, "Wallet failed verification")
			return CodeFailVerfication, nil
		}
	}

	if lightWeight == false {

		err := app.Set(txHash, rawAsset)
		if err != nil {
			return CodeDatabaseFail, nil
		}

		//Write the Pointer Key
		// ABCDE.0 = txHash
		pointerKey := KeySuffix(wallet.Key(), newAssetIndexString)
		msg := fmt.Sprintf("Wallet set (assetid.index:tx)     %v:%v", hex.EncodeToString(pointerKey), hex.EncodeToString(txHash))
		dumpMessage(4, msg)
		err = app.Set(pointerKey, txHash)
		if err != nil {
			return CodeDatabaseFail, nil
		}

		//Write the lastest index to the asset key
		// ABCDE = 0
		msg = fmt.Sprintf("Wallet set (assetid:latest_index) %v:%v", hex.EncodeToString(wallet.Key()), newAssetIndexString)
		dumpMessage(4, msg)
		err = app.Set(wallet.Key(), []byte(newAssetIndexString))
		if err != nil {
			return CodeDatabaseFail, nil
		}
	}
	return CodeTypeOK, nil
}

func (app *QredoChain) processGroup(group *assets.Group, lightWeight bool) (TransactionCode, []abcitypes.Event) {
	fmt.Printf("Process an Group\n")
	return CodeFailVerfication, nil
}

func (app *QredoChain) VerifyIDDoc(iddoc *assets.IDDoc) bool {
	//Check signature
	verify, err := iddoc.Verify(iddoc)
	if err != nil {
		return false
	}
	if verify == false {
		return false
	}

	//Check Payload fields
	payload, err := iddoc.Payload()
	if err != nil {
		return false
	}
	if payload == nil {
		return false
	}

	if payload.AuthenticationReference == "" ||
		payload.BeneficiaryECPublicKey == nil ||
		payload.SikePublicKey == nil ||
		payload.BLSPublicKey == nil {
		return false
	}

	if iddoc.CurrentAsset.Asset.Index != 0 {
		return false
	}

	return true
}

func (app *QredoChain) VerifyWalletUpdate(wallet *assets.Wallet) bool {
	return true
	// verify, err := wallet.OnChainFullVerify(app)
	// if err != nil {
	// 	return false
	// }
	// if verify == false {
	// 	return false
	// }
	//	return true
}

func (app *QredoChain) VerifyNewWallet(wallet *assets.Wallet) bool {

	// signers := wallet.CurrentAsset.Signers
	// if signers == nil || len(signers) == 0 || len(signers) > 1 {
	// 	return false
	// }

	// idkey := signers[0]

	// verify, err := wallet.Verify()
	// if err != nil {
	// 	return false
	// }
	// if verify == false {
	// 	return false
	// }

	// //Check Payload fields
	// payload, err := wallet.Payload()
	// if err != nil {
	// 	return false
	// }
	// if payload == nil {
	// 	return false
	// }

	return true
}

func (app *QredoChain) VerifyGroupUpdate(iddoc *assets.Group) bool {
	return true
}

func (app *QredoChain) VerifyNewGroup(iddoc *assets.Group) bool {
	return true
}

// func (app *QredoChain) processQuery(iddoc *assets.Group) bool {
// }