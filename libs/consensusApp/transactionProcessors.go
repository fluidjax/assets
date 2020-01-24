package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/tendermint/abci/example/code"
)

func (app *KVStoreApplication) processTX(tx []byte, lightWeight bool) uint32 {
	// 	The workhorse of the application - non-optional.
	// 	Execute the transaction in full.
	// 	ResponseDeliverTx.Code == 0 only if the transaction is fully valid.

	//Decode the Asset
	signedAsset, err := decodeTX(tx)
	if err != nil {
		return code.CodeTypeEncodingError
	}
	//Make TX Hash
	txHashA := sha256.Sum256(tx)
	txHash := txHashA[:]

	//Retrieve the Asset ID
	assetID := signedAsset.Asset.GetID()
	if assetID == nil {
		return code.CodeTypeEncodingError
	}

	//Process the Transaction
	switch signedAsset.Asset.GetType() {
	case protobuffer.PBAssetType_wallet:
		wallet, err := assets.ReBuildWallet(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError
		}
		return app.processWallet(wallet, tx, txHash, lightWeight)
	case protobuffer.PBAssetType_iddoc:
		iddoc, err := assets.ReBuildIDDoc(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError
		}
		return app.processIDDoc(iddoc, tx, txHash, lightWeight)
	case protobuffer.PBAssetType_Group:
		group, err := assets.ReBuildGroup(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError
		}
		return app.processGroup(group, lightWeight)
	}
	return code.CodeTypeEncodingError
}

func (app *KVStoreApplication) processIDDoc(iddoc *assets.IDDoc, rawAsset []byte, txHash []byte, lightWeight bool) uint32 {
	if app.exists(txHash) {
		//Usually the tx cache takes care of this, but once its full, we need to stop duplicates of very old transactions
		dumpMessage(2, "Fail to add IDDoc - tx already in chain")
		return CodeAlreadyExists
	}

	//IDDoc is immutable so if this AssetID already has a value we can't update it.
	if app.exists(iddoc.Key()) == true {
		dumpMessage(2, "Fail to add IDDoc - already exists")
		return CodeAlreadyExists
	}

	//Check the IDDoc is valid
	if app.VerifyIDDoc(iddoc) == false {
		return CodeFailVerfication
	}

	//Add pointer from AssetID to the txHash of the Object
	if lightWeight == false {
		err := app.Set(iddoc.Key(), txHash)
		if err != nil {
			return CodeDatabaseFail
		}
	}
	return CodeTypeOK
}

func (app *KVStoreApplication) processWallet(wallet *assets.Wallet, rawAsset []byte, txHash []byte, lightWeight bool) uint32 {
	if app.exists(txHash) {
		//Usually the tx cache takes care of this, but once its full, we need to stop duplicates of very old transactions
		dumpMessage(2, "Fail to add wallet - tx already in chain\n")
		return CodeAlreadyExists
	}

	currentIndex, err := app.Get(wallet.Key())
	if err != nil {
		return CodeDatabaseFail
	}
	var newAssetIndexString string

	if currentIndex == nil {
		//New Wallet
		if app.VerifyNewWallet(wallet) == false {
			dumpMessage(4, "Wallet failed verification")
			return CodeFailVerfication
		}
		newAssetIndexString = IndexFormater(0)

	} else {
		//Check we are correctly incrementing the index
		currentIndexInteger, err := strconv.ParseInt(string(currentIndex), 10, 64)
		if err != nil {
			dumpMessage(4, "Failed to Parse Current Index")
			return CodeFailVerfication
		}
		newIndex := wallet.CurrentAsset.Asset.Index
		if newIndex != currentIndexInteger+1 {
			dumpMessage(2, "Invalid Wallet Index\n")
			return CodeFailVerfication
		}
		newAssetIndexString = IndexFormater(newIndex)
		//Wallet update
		if app.VerifyWalletUpdate(wallet) == false {
			dumpMessage(4, "Wallet failed verification")
			return CodeFailVerfication
		}
	}

	if lightWeight == false {
		//Write the Pointer Key
		// ABCDE.0 = txHash
		pointerKey := KeySuffix(wallet.Key(), newAssetIndexString)
		msg := fmt.Sprintf("Wallet set (assetid.index:tx)     %v:%v", hex.EncodeToString(pointerKey), hex.EncodeToString(txHash))
		dumpMessage(4, msg)
		err = app.Set(pointerKey, txHash)
		if err != nil {
			return CodeDatabaseFail
		}

		//Write the lastest index to the asset key
		// ABCDE = 0
		msg = fmt.Sprintf("Wallet set (assetid:latest_index) %v:%v", hex.EncodeToString(wallet.Key()), newAssetIndexString)
		dumpMessage(4, msg)
		err = app.Set(wallet.Key(), []byte(newAssetIndexString))
		if err != nil {
			return CodeDatabaseFail
		}
	}
	return CodeTypeOK
}

func (app *KVStoreApplication) processGroup(group *assets.Group, lightWeight bool) uint32 {
	fmt.Printf("Process an Group\n")
	return CodeFailVerfication
}

func (app *KVStoreApplication) VerifyIDDoc(iddoc *assets.IDDoc) bool {
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

	if payload.Descriptor == nil ||
		payload.AuthenticationReference == "" ||
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

func (app *KVStoreApplication) VerifyWalletUpdate(wallet *assets.Wallet) bool {
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

func (app *KVStoreApplication) VerifyNewWallet(wallet *assets.Wallet) bool {

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

func (app *KVStoreApplication) VerifyGroupUpdate(iddoc *assets.Group) bool {
	return true
}

func (app *KVStoreApplication) VerifyNewGroup(iddoc *assets.Group) bool {
	return true
}
