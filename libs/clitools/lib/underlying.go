package lib

import (
	"encoding/hex"
	"encoding/json"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
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
	payload.Address = []byte(cJSON.Address)
	payload.TxID = []byte(cJSON.TxID)
	under.AssetKeyFromPayloadHash()

	txid := ""
	if broadcast == true {
		var err error
		txid, err = cliTool.NodeConn.PostTx(under)
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
