package qc

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
)

type WalletUpdate struct {
	ExistingWalletAssetID string     `json:"ExistingWalletAssetID"`
	NewOwner              string     `json:"newownerseed"`
	TransferType          int64      `json:"transferType"`
	Currency              string     `json:"currency"`
	Transfer              []Transfer `json:"Transfer"`
}

func (cliTool *CLITool) UpdateWalletWithJSON(jsonParams string) (err error) {
	//Load existing wallet from AssetID
	//Load all the IDDocs

	//wallet from JSON
	uwJSON := &WalletUpdate{}
	err = json.Unmarshal([]byte(jsonParams), uwJSON)
	if err != nil {
		return err
	}
	updatedWallet, err := cliTool.WalletFromWalletUpdateJSON(uwJSON)

	//Updated wallet complete, return for signing
	msg, err := updatedWallet.SerializeAsset()
	if err != nil {
		return errors.New("Failed to serialize payload")
	}

	addResultBinaryItem("serializedUpdate", msg)
	ppResult()
	return nil
}

func (cliTool *CLITool) WalletFromWalletUpdateJSON(walletUpdate *WalletUpdate) (*assets.Wallet, error) {
	//Decode the JSON

	//Get the New Owner IDDoc
	idNewOwnerKey, err := hex.DecodeString(walletUpdate.ExistingWalletAssetID)
	if err != nil {
		return nil, err
	}

	newOwnerIDDoc, err := assets.LoadIDDoc(cliTool.NodeConn, idNewOwnerKey)
	if err != nil {
		return nil, err
	}

	//Get the Existing Wallet
	existingWalletKey, err := hex.DecodeString(walletUpdate.ExistingWalletAssetID)
	if err != nil {
		return nil, err
	}

	originalWallet, err := assets.LoadWallet(cliTool.NodeConn, existingWalletKey)
	if err != nil {
		return nil, err
	}

	//Make New Wallet based on Existing
	updatedWallet, err := assets.NewUpdateWallet(originalWallet, newOwnerIDDoc)
	if err != nil {
		return nil, err
	}

	updatedWallet.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType(walletUpdate.TransferType)
	updatedWallet.DataStore = cliTool.NodeConn

	return updatedWallet, nil
}
