package qredochain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

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

	//fmt.Println("txHash = ", txHash[:])

	//Retrieve the Asset ID
	assetID := signedAsset.Asset.GetID()
	if assetID == nil {
		return code.CodeTypeEncodingError, nil
	}

	//Process the Transaction
	switch signedAsset.Asset.GetType() {
	case protobuffer.PBAssetType_Wallet:
		wallet, err := assets.ReBuildWallet(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError, nil
		}
		code, events := app.processWallet(wallet, tx, txHash, lightWeight)
		return uint32(code), events
	case protobuffer.PBAssetType_Iddoc:
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
	case protobuffer.PBAssetType_Underlying:
		underlying, err := assets.ReBuildUnderlying(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError, nil
		}
		code, events := app.processUnderlying(underlying, lightWeight)
		return uint32(code), events
	case protobuffer.PBAssetType_MPC:
		mpc, err := assets.ReBuildMPC(signedAsset, assetID)
		if err != nil {
			return code.CodeTypeEncodingError, nil
		}
		code, events := app.processMPC(mpc, lightWeight)
		return uint32(code), events
	}
	return code.CodeTypeEncodingError, nil
}

func (app *QredoChain) processMPC(mpc *assets.MPC, lightWeight bool) (code TransactionCode, events []abcitypes.Event) {
	//fmt.Printf("Process an MPC\n")

	payload, err := mpc.Payload()
	if err != nil {
		return CodeTypeEncodingError, nil
	}
	address := payload.Address
	assetID := payload.AssetID

	if lightWeight == false {
		//Set the KV to link between asset & address
		app.SetWithSuffix(assetID, ".as2ad", address)
		app.SetWithSuffix(address, ".ad2as", assetID)
	}

	return CodeTypeOK, events
}

func (app *QredoChain) processUnderlying(underlying *assets.Underlying, lightWeight bool) (code TransactionCode, events []abcitypes.Event) {
	//fmt.Printf("Process an Underlying\n")

	payload, err := underlying.Payload()
	if err != nil {
		return CodeTypeEncodingError, nil
	}
	address := []byte(payload.Address)
	amount := payload.Amount
	underlyingTxID := []byte(payload.TxID)

	exists, err := app.GetWithSuffix(underlyingTxID, ".UTxID")
	if err != nil || exists != nil {
		return CodeConsensusError, nil
	}

	if lightWeight == false {
		//underlying has Crypto Address - get AssetID from KV Store
		assetID, err := app.GetWithSuffix(address, ".ad2as")
		if err != nil {
			return CodeTypeEncodingError, nil
		}
		code := app.addToBalanceKey(assetID, amount)
		if code != 0 {
			return code, nil
		}

	}
	return CodeTypeOK, events
}

func (app *QredoChain) subtractFromBalanceKey(assetID []byte, amount int64) (code TransactionCode) {
	currentBalance, code := app.getBalanceKey(assetID)
	newBalance := currentBalance - amount
	if newBalance < 0 {
		return CodeConsensusBalanceError
	}
	return app.setBalanceKey(assetID, newBalance)
}

func (app *QredoChain) addToBalanceKey(assetID []byte, amount int64) (code TransactionCode) {
	currentBalance, code := app.getBalanceKey(assetID)
	newBalance := currentBalance + amount
	return app.setBalanceKey(assetID, newBalance)
}

func (app *QredoChain) setBalanceKey(assetID []byte, newBalance int64) (code TransactionCode) {
	//Convert new balance to bytes and save for AssetID
	newBalanceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(newBalanceBytes, uint64(newBalance))
	err := app.SetWithSuffix(assetID, ".balance", newBalanceBytes)
	if err != nil {
		return CodeConsensusBalanceError
	}
	return 0
}

func (app *QredoChain) getBalanceKey(assetID []byte) (amount int64, code TransactionCode) {
	currentBalanceBytes, err := app.GetWithSuffix(assetID, ".balance")
	if currentBalanceBytes == nil || err != nil {
		return 0, CodeConsensusBalanceError
	}
	currentBalance := int64(binary.LittleEndian.Uint64(currentBalanceBytes))
	return currentBalance, 0
}

func (app *QredoChain) processIDDoc(iddoc *assets.IDDoc, rawAsset []byte, txHash []byte, lightWeight bool) (code TransactionCode, events []abcitypes.Event) {

	if app.exists(txHash) {
		//Usually the tx cache takes care of this, but once its full, we need to stop duplicates of very old transactions
		dumpMessage(2, "Fail to add IDDoc - tx already in Consensus Database")
		return CodeAlreadyExists, nil
	}

	//IDDoc is immutable so if this AssetID already has a value we can't update it.
	if app.exists(iddoc.Key()) == true {
		dumpMessage(2, "Fail to add IDDoc - in Consensus Database")
		return CodeAlreadyExists, nil
	}

	//Check the IDDoc is valid
	if app.VerifyIDDoc(iddoc) == false {
		return CodeFailVerfication, nil
	}

	//Add pointer from AssetID to the txHash of the Object
	if lightWeight == false {
		err1 := app.Set(txHash, rawAsset)
		if err1 != nil {
			return CodeDatabaseFail, nil
		}
		err2 := app.Set(iddoc.Key(), txHash)
		if err2 != nil {
			return CodeDatabaseFail, nil
		}
		events = processTags(iddoc.CurrentAsset.Asset.Tags)
	}
	return CodeTypeOK, events
}

func (app *QredoChain) processWallet(wallet *assets.Wallet, rawAsset []byte, txHash []byte, lightWeight bool) (code TransactionCode, events []abcitypes.Event) {
	if app.exists(txHash) {
		//Usually the tx cache takes care of this, but once its full, we need to stop duplicates of very old transactions
		dumpMessage(2, "Fail to add wallet - tx already in chain\n")
		return CodeAlreadyExists, nil
	}

	exists, err := app.Get(wallet.Key())
	if err != nil {
		return CodeDatabaseFail, nil
	}

	if exists == nil {
		return app.processWalletCreate(wallet, rawAsset, txHash, lightWeight)
	} else {
		return app.processWalletUpdate(wallet, rawAsset, txHash, lightWeight)
	}

}

func (app *QredoChain) processWalletUpdate(wallet *assets.Wallet, rawAsset []byte, txHash []byte, lightWeight bool) (code TransactionCode, events []abcitypes.Event) {
	if lightWeight == false {
		err := app.Set(wallet.Key(), txHash)
		if err != nil {
			return CodeDatabaseFail, nil
		}
		events = processTags(wallet.CurrentAsset.Asset.Tags)

		//Loop through all the transfers out and update their destinations
		payload, err := wallet.Payload()
		if err != nil {
			return CodeFailVerfication, nil
		}

		currentBalance, code := app.getBalanceKey(wallet.Key())
		if code != 0 {
			return code, nil
		}

		//Check we have enough - Pass 1
		var totalOutgoing int64
		for _, wt := range payload.WalletTransfers {
			res := bytes.Compare(wt.AssetID, wallet.Key())
			if res == 0 {
				//this is money coming back to self, just ignore it
				continue
			}
			totalOutgoing = totalOutgoing + wt.Amount
		}

		if totalOutgoing > currentBalance {
			return CodeInsufficientFunds, nil
		}

		//We have enough funds, do the database updates for transfer Pass 2

		for _, wt := range payload.WalletTransfers {
			res := bytes.Compare(wt.AssetID, wallet.Key())
			if res == 0 {
				//this is money coming back to self, just ignore it
				continue
			}
			amount := wt.Amount
			destinationAssetID := wt.AssetID
			app.addToBalanceKey(destinationAssetID, amount)
			app.subtractFromBalanceKey(wallet.Key(), amount)
		}
	}
	return CodeTypeOK, events
}

func (app *QredoChain) processWalletCreate(wallet *assets.Wallet, rawAsset []byte, txHash []byte, lightWeight bool) (code TransactionCode, events []abcitypes.Event) {
	if lightWeight == false {
		err1 := app.Set(txHash, rawAsset)
		if err1 != nil {
			return CodeDatabaseFail, nil
		}
		err2 := app.Set(wallet.Key(), txHash)
		if err2 != nil {
			return CodeDatabaseFail, nil
		}
		events = processTags(wallet.CurrentAsset.Asset.Tags)

		app.setBalanceKey(wallet.Key(), 0)

	}
	return CodeTypeOK, events
}

func (app *QredoChain) processGroup(group *assets.Group, lightWeight bool) (code TransactionCode, events []abcitypes.Event) {
	fmt.Printf("Process an Group\n")

	return CodeFailVerfication, events
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

func processTags(tags map[string][]byte) []types.Event {
	var attributes []kv.Pair
	for key, value := range tags {
		kvpair := kv.Pair{Key: []byte(key), Value: value}
		attributes = append(attributes, kvpair)
	}
	events := []types.Event{
		{
			Type:       "tag",
			Attributes: attributes,
		},
	}
	return events
}
