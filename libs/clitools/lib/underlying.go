package qc

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
)

type UnderlyingJSON struct {
	Type               int64  `json:"type"`
	CryptoCurrencyCode int64  `json:"CryptoCurrencyCode"`
	Proof              string `json:"Proof"`
	Amount             int64  `json:"Amount"`
	Address            string `json:"Address"`
	TxID               string `json:"TxID"`
}

func (cliTool *CLITool) CreateUnderlyingWithJSON(jsonParams string, broadcast bool) (err error) {
	cJSON := &UnderlyingJSON{}
	err = json.Unmarshal([]byte(jsonParams), cJSON)
	if err != nil {
		return err
	}
	under, err := assets.NewUnderlying()
	if err != nil {
		return err
	}
	payload, err := under.Payload()
	if err != nil {
		return err
	}

	payload.Type = protobuffer.PBUnderlyingType(cJSON.Type)
	payload.CryptoCurrencyCode = protobuffer.PBCryptoCurrency(cJSON.CryptoCurrencyCode)

	proofBin, err := hex.DecodeString(cJSON.Proof)
	payload.Proof = proofBin

	amount := int64(cJSON.Amount)
	payload.Amount = amount
	payload.Address = cJSON.Address
	payload.TxID = cJSON.TxID
	under.AssetKeyFromPayloadHash()

	txid := ""
	if broadcast == true {
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(under)
		if code != 0 {
			print(err.Error())
			return errors.Wrap(err, "TX Fails verifications")
		}
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := under.SerializeSignedAsset()
	if err != nil {
		return err
	}

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", under.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", under.CurrentAsset)

	ppResult()

	return nil
}
