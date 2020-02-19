package qc

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"

	//"github.com/golang/protobuf/proto"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/qredochain"
	//	"github.com/qredo/assets/libs/protobuffer"
)

type SignJSON struct {
	Seed string `json:"seed"`
	Msg  string `json:"msg"`
}

func (cliTool *CLITool) Sign(jsonParams string) (err error) {
	//Decode the JSON
	signJSON := &SignJSON{}
	err = json.Unmarshal([]byte(jsonParams), signJSON)
	if err != nil {
		return err
	}

	seed, err := hex.DecodeString(signJSON.Seed)
	if err != nil {
		return err
	}

	key := assets.KeyFromSeed(seed)
	iddoc, err := assets.LoadIDDoc(cliTool.NodeConn, key)
	iddoc.Seed = seed

	msgToSign, err := hex.DecodeString(signJSON.Msg)
	if err != nil {
		return err
	}

	signature, err := assets.Sign(msgToSign, iddoc)

	addResultBinaryItem("signature", signature)
	ppResult()
	return

}

// transfersignaturesjson.go

// aggregatesignjson.go

type AggregateSignJSON struct {
	Sigs         []Sig        `json:"Sigs"`
	WalletUpdate WalletUpdate `json:"walletUpdate"`
}

// sig.go

type Sig struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	Signature    string `json:"signature"`
}

func (cliTool *CLITool) AggregateSign(jsonParams string, broadcast bool) (err error) {
	//Decode the JSON
	agJSON := &AggregateSignJSON{}
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
	updatedWallet, err := cliTool.WalletFromWalletUpdateJSON(&agJSON.WalletUpdate)

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
