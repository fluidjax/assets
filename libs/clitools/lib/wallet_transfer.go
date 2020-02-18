package qc

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/qredochain"
)

type UpdateWalletJSON struct {
	ExistingWalletAssetID string     `json:"ExistingWalletAssetID"`
	NewOwner              string     `json:"newownerseed"`
	TransferType          int64      `json:"transferType"`
	Currency              string     `json:"currency"`
	Transfer              []Transfer `json:"Transfer"`
}

func (cliTool *CLITool) UpdateWalletWithJSON(jsonParams string, broadcast bool) (err error) {
	//Load existing wallet from AssetID
	//Load all the IDDocs

	uwJSON := &UpdateWalletJSON{}
	err = json.Unmarshal([]byte(jsonParams), uwJSON)
	if err != nil {
		return err
	}

	idNewOwnerKey, err := hex.DecodeString(uwJSON.ExistingWalletAssetID)
	newOwnerIDDoc, err := assets.LoadIDDoc(cliTool.NodeConn, idNewOwnerKey)

	existingWalletKey, err := hex.DecodeString(uwJSON.ExistingWalletAssetID)
	originalWallet, err := assets.LoadWallet(cliTool.NodeConn, existingWalletKey)

	updatedWallet, err := assets.NewUpdateWallet(originalWallet, newOwnerIDDoc)

	txid := ""
	if broadcast == true {
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(updatedWallet)
		if code != 0 {
			print(err.Error())
			return errors.Wrap(err, "TX Fails verifications")
		}
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := updatedWallet.SerializeSignedAsset()
	if err != nil {
		return err
	}
	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", updatedWallet.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", updatedWallet.CurrentAsset)

	ppResult()

	return nil
}
