package lib

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
)

//CreateWalletWithJSON -
func (cliTool *CLITool) CreateWalletWithJSON(jsonParams string, broadcast bool) (err error) {
	cwJSON := &CreateWalletJSON{}
	err = json.Unmarshal([]byte(jsonParams), cwJSON)
	if err != nil {
		return err
	}
	seedHex := cwJSON.Ownerseed
	iddoc, err := cliTool.GetIDDocForSeed(seedHex)
	if err != nil {
		return err
	}

	wallet, err := assets.NewWallet(iddoc, protobuffer.PBCryptoCurrency(cwJSON.Currency))
	wallet.DataStore = cliTool.NodeConn
	if err != nil {
		return err
	}

	var truths []string
	for _, trans := range cwJSON.Transfer {

		binParticipants := map[string][]byte{}
		for _, v := range trans.Participants {
			binVal, err := hex.DecodeString(v.ID)
			if err != nil {
				return err
			}
			binParticipants[v.Name] = binVal
		}
		transferType := protobuffer.PBTransferType(trans.TransferType)
		wallet.AddTransfer(transferType, trans.Expression, &binParticipants, trans.Description)
		truthTable, err := wallet.TruthTable(transferType)
		if err != nil {
			return err
		}

		for _, v := range truthTable {
			x := fmt.Sprintf("%d:%s", trans.TransferType, v)
			truths = append(truths, base64.StdEncoding.EncodeToString([]byte(x)))
		}

	}

	wallet.Sign(iddoc)

	txid := ""
	if broadcast == true {
		var err error
		txid, err = cliTool.NodeConn.PostTx(wallet)
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := wallet.SerializeSignedAsset()
	if err != nil {
		return err
	}

	res["truthtable"] = truths

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", wallet.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", wallet.CurrentAsset)

	ppResult()
	return nil
}

func (cliTool *CLITool) PrepareWalletUpdateWithJSON(jsonParams string) (err error) {
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
		var err error
		txid, err = cliTool.NodeConn.PostTx(updatedWallet)
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
