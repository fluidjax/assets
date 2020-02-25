package qc

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
)

func (cliTool *CLITool) PrepareWalletUpdateWithJSON(jsonParams string) (err error) {
	//Load existing wallet from AssetID
	//Load all the IDDocs

	//wallet from JSON
	uwJSON := &WalletUpdatePayload{}
	err = json.Unmarshal([]byte(jsonParams), uwJSON)
	if err != nil {
		return err
	}
	updatedWallet, err := cliTool.walletFromWalletUpdateJSON(uwJSON)
	if err != nil {
		return err
	}

	//Updated wallet complete, return for signing
	msg, err := updatedWallet.SerializeAsset()
	if err != nil {
		return errors.New("Failed to serialize payload")
	}

	addResultBinaryItem("serializedUpdate", msg)
	ppResult()
	return nil
}

func (cliTool *CLITool) walletFromWalletUpdateJSON(signedUpdate *WalletUpdatePayload) (*assets.Wallet, error) {
	//Decode the JSON

	//Get the New Owner IDDoc
	idNewOwnerKey, err := hex.DecodeString(signedUpdate.Newowner)
	if err != nil {
		return nil, err
	}

	newOwnerIDDoc, err := assets.LoadIDDoc(cliTool.NodeConn, idNewOwnerKey)
	if err != nil {
		return nil, err
	}

	//Get the Existing Wallet
	existingWalletKey, err := hex.DecodeString(signedUpdate.ExistingWalletAssetID)
	if err != nil {
		return nil, err
	}

	originalWallet, err := assets.LoadWallet(cliTool.NodeConn, existingWalletKey)
	if err != nil {
		return nil, err
	}
	originalWallet.DataStore = cliTool.NodeConn
	//Make New Wallet based on Existing
	updatedWallet, err := assets.NewUpdateWallet(originalWallet, newOwnerIDDoc)
	if err != nil {
		return nil, err
	}

	var truths []string
	for _, trans := range signedUpdate.Transfer {

		binParticipants := map[string][]byte{}
		for _, v := range trans.Participants {
			binVal, err := hex.DecodeString(v.ID)
			if err != nil {
				return nil, err
			}
			binParticipants[v.Name] = binVal
		}
		transferType := protobuffer.PBTransferType(trans.TransferType)
		updatedWallet.AddTransfer(transferType, trans.Expression, &binParticipants, trans.Description)
		truthTable, err := updatedWallet.TruthTable(transferType)
		if err != nil {
			return nil, err
		}

		for _, v := range truthTable {
			x := fmt.Sprintf("%d:%s", trans.TransferType, v)
			truths = append(truths, base64.StdEncoding.EncodeToString([]byte(x)))
		}
	}

	//Add in the WalletTransfers - ie. payment destinations
	for _, wt := range signedUpdate.WalletTransfers {
		to, err := hex.DecodeString(wt.To)
		if err != nil {
			return nil, err
		}
		assetID, err := hex.DecodeString(wt.Assetid)
		if err != nil {
			return nil, err
		}

		updatedWallet.AddWalletTransfer(to, wt.Amount, assetID)
	}

	updatedWallet.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType(signedUpdate.TransferType)
	updatedWallet.DataStore = cliTool.NodeConn

	return updatedWallet, nil
}

func (cliTool *CLITool) AggregateWalletSign(jsonParams string, broadcast bool) (err error) {
	//Decode the JSON
	agJSON := &WalletUpdate{}
	err = json.Unmarshal([]byte(jsonParams), agJSON)
	if err != nil {
		return err
	}

	var transferSignatures []assets.SignatureID
	for _, sig := range agJSON.Sigs {
		key, err := hex.DecodeString(sig.ID)
		if err != nil {
			return err
		}
		iddoc, err := assets.LoadIDDoc(cliTool.NodeConn, key)
		signature, err := hex.DecodeString(sig.Signature)
		if err != nil {
			return err
		}
		sid := assets.SignatureID{IDDoc: iddoc, Abbreviation: sig.Abbreviation, Signature: signature}
		transferSignatures = append(transferSignatures, sid)
	}

	//Rebuild the Wallet from the TX supplied
	updatedWallet, err := cliTool.walletFromWalletUpdateJSON(&agJSON.WalletUpdatePayload)

	err = updatedWallet.AggregatedSign(transferSignatures)
	if err != nil {
		return err
	}

	verify, err := updatedWallet.FullVerify()

	if verify == false {
		return errors.New("Error failed to verify final update wallet transaction")
	}

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
	return
}
