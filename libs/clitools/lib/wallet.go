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

// createwalletjson.go
type CreateWalletJSON struct {
	TransferType int64      `json:"transferType"`
	Ownerseed    string     `json:"ownerseed"`
	Currency     string     `json:"currency"`
	Transfer     []Transfer `json:"Transfer"`
}

// transfer.go

type Transfer struct {
	TransferType int64         `json:"TransferType"`
	Expression   string        `json:"Expression"`
	Description  string        `json:"description"`
	Participants []Participant `json:"participants"`
}

// participant.go

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
