package qc

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
)

type CreateWalletJSON struct {
	Ownerseed    string        `json:"ownerseed"`
	Expression   string        `json:"expression"`
	TransferType int64         `json:"transferType"`
	Participants []Participant `json:"participants"`
	Currency     string        `json:"currency"`
}

// participant.go

/*
{
	"ownerseed":"3772e3fa880e1912498d2fc48a367a2058c69ea4bf6ec3cf41fbbb6d8089f8868f3c46e31d8e9ab251ea5e4c6f5ded53",
	"expression":"t1 + t2 + t3 > 1 & p",
	"transferType":1,
	"participants":[
			{"name":"p","ID":"3b1c3e7563f6f174ba4cc01e77bd69f3999e6e81e74b5d044c69336e2751045a"},
			{"name":"t1","ID":"40cac6105b4a025a0815c96e630d75414982ff2b4aa5b500011fc59f50ad3c4d"},
			{"name":"t2","ID":"b29d6d6fb277eef333e1dfc79e4bed516cf18bf5ce3eae808a4d941c081f7afa"},
			{"name":"t3","ID":"cc4d921dc8f8ebe163ee476b4ce9ed06412be60d9d94c9e0316fa2321c2eaa20"}
			],
	"currency":"BTC"
	}
*/

type Participant struct {
	Name string `json:"name"`
	ID   string `json:"ID"`
}

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

	wallet, err := assets.NewWallet(iddoc, cwJSON.Currency)
	wallet.DataStore = cliTool.NodeConn
	if err != nil {
		return err
	}

	transferType := protobuffer.PBTransferType(cwJSON.TransferType)

	binParticipants := map[string][]byte{}
	for _, v := range cwJSON.Participants {
		binVal, err := hex.DecodeString(v.ID)
		if err != nil {
			return err
		}
		binParticipants[v.Name] = binVal
	}

	wallet.AddTransfer(transferType, cwJSON.Expression, &binParticipants)
	truthTable, err := wallet.TruthTable(transferType)
	if err != nil {
		return err
	}

	wallet.Sign(iddoc)

	txid := ""
	if broadcast == true {
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(wallet)
		if code != 0 {
			print(err.Error())
			return errors.Wrap(err, "TX Fails verifications")
		}
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := wallet.SerializeSignedAsset()
	if err != nil {
		return err
	}

	var truths []string
	for _, v := range truthTable {
		truths = append(truths, base64.StdEncoding.EncodeToString([]byte(v)))
	}
	res["truthtable"] = truths

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", wallet.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", wallet.CurrentAsset)

	ppResult()
	return nil
}

func (cliTool *CLITool) CreateWallet(seedHex string, broadcast bool) (err error) {

	iddoc, err := cliTool.GetIDDocForSeed(seedHex)
	if err != nil {
		return err
	}

	wallet, err := assets.NewWallet(iddoc, "BTC")
	if err != nil {
		return err
	}
	wallet.Sign(iddoc)

	txid := ""
	if broadcast == true {
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(wallet)
		if code != 0 {
			print(err.Error())
			return errors.Wrap(err, "TX Fails verifications")
		}

		if err != nil {
			return err
		}
	}

	serializedAsset, err := wallet.SerializeAsset()
	if err != nil {
		return err
	}

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", wallet.Key())
	addResultBinaryItem("serialized", serializedAsset)
	addResultSignedAsset("object", wallet.CurrentAsset)
	ppResult()

	return
}
