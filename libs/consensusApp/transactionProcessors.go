package main

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/qredo/assets/libs/assets"
	"github.com/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

func (app *KVStoreApplication) processIDDoc(iddoc *assets.IDDoc, rawAsset []byte, txHash []byte) abcitypes.ResponseDeliverTx {
	if app.exists(txHash) {
		//Usually the tx cache takes care of this, but once its full, we need to stop duplicates of very old transactions
		dumpMessage(2, "Fail to add IDDoc - tx already in chain\n")
		return types.ResponseDeliverTx{Code: CodeAlreadyExists, Events: nil}
	}

	//IDDoc is immutable so if this AssetID already has a value we can't update it.
	if app.exists(iddoc.Key()) == true {
		dumpMessage(2, "Fail to add IDDoc - already exists\n")
		return types.ResponseDeliverTx{Code: CodeAlreadyExists, Events: nil}
	}

	//Check the IDDoc is valid
	if app.VerifyIDDoc(iddoc) == false {
		return types.ResponseDeliverTx{Code: CodeFailVerfication, Events: nil}
	}

	//Add pointer from AssetID to the txHash of the Object
	err := app.Set(iddoc.Key(), txHash)
	if err != nil {
		return types.ResponseDeliverTx{Code: CodeDatabaseFail, Events: nil}
	}
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: nil}
}

func (app *KVStoreApplication) processWallet(wallet *assets.Wallet, rawAsset []byte, txHash []byte) abcitypes.ResponseDeliverTx {
	if app.exists(txHash) {
		//Usually the tx cache takes care of this, but once its full, we need to stop duplicates of very old transactions
		dumpMessage(2, "Fail to add wallet - tx already in chain\n")
		return types.ResponseDeliverTx{Code: CodeAlreadyExists, Events: nil}
	}

	currentIndex, err := app.Get(wallet.Key())
	if err != nil {
		return types.ResponseDeliverTx{Code: CodeDatabaseFail, Events: nil}
	}
	var newAssetIndexString string

	if currentIndex == nil {
		msg := fmt.Sprintf("Add new wallet %v \n", hex.EncodeToString(wallet.Key()))
		dumpMessage(4, msg)

		//New Wallet
		if app.VerifyNewWallet(wallet) == false {
			dumpMessage(4, "Wallet failed verification")
			return types.ResponseDeliverTx{Code: CodeFailVerfication, Events: nil}
		}
		newAssetIndexString = IndexFormater(0)

	} else {
		//Check we are correctly incrementing the index
		currentIndexInteger, err := strconv.ParseInt(string(currentIndex), 10, 64)
		if err != nil {
			dumpMessage(4, "Failed to Parse Current Index")
			return types.ResponseDeliverTx{Code: CodeFailVerfication, Events: nil}
		}
		newIndex := wallet.CurrentAsset.Asset.Index
		if newIndex != currentIndexInteger+1 {
			dumpMessage(2, "Invalid Wallet Index\n")
			return types.ResponseDeliverTx{Code: CodeFailVerfication, Events: nil}
		}
		newAssetIndexString = IndexFormater(newIndex)
		//Wallet update
		if app.VerifyWalletUpdate(wallet) == false {
			dumpMessage(4, "Wallet failed verification")
			return types.ResponseDeliverTx{Code: CodeFailVerfication, Events: nil}
		}
	}

	//Write the Pointer Key
	// ABCDE.0 = txHash
	pointerKey := KeySuffix(wallet.Key(), newAssetIndexString)
	msg := fmt.Sprintf("Wallet set (assetid.index:tx) %v:%v \n", pointerKey, hex.EncodeToString(txHash))
	dumpMessage(4, msg)
	err = app.Set(pointerKey, txHash)
	if err != nil {
		return types.ResponseDeliverTx{Code: CodeDatabaseFail, Events: nil}
	}

	//Write the lastest index to the asset key
	// ABCDE = 0
	msg = fmt.Sprintf("Wallet set (assetid:latest_index) %v : %v \n", hex.EncodeToString(wallet.Key()), newAssetIndexString)
	dumpMessage(4, msg)
	err = app.Set(wallet.Key(), []byte(newAssetIndexString))
	if err != nil {
		return types.ResponseDeliverTx{Code: CodeDatabaseFail, Events: nil}
	}
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: nil}
}

func (app *KVStoreApplication) processGroup(group *assets.Group) error {
	fmt.Printf("Process an Group\n")
	return nil
}

func (app *KVStoreApplication) VerifyIDDoc(iddoc *assets.IDDoc) bool {
	return true
}

func (app *KVStoreApplication) VerifyWalletUpdate(iddoc *assets.Wallet) bool {
	return true
}

func (app *KVStoreApplication) VerifyNewWallet(iddoc *assets.Wallet) bool {
	return true
}

func (app *KVStoreApplication) VerifyGroupUpdate(iddoc *assets.Group) bool {
	return true
}

func (app *KVStoreApplication) VerifyNewGroup(iddoc *assets.Group) bool {
	return true
}
