package qc

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
)

type MPCJSON struct {
	Type      int64  `json:"type"`
	Address   string `json:"Address"`
	Signature string `json:"Signature"`
	AssetID   string `json:"AssetID"`
}

func (cliTool *CLITool) CreateMPCWithJSON(jsonParams string, broadcast bool) (err error) {
	cJSON := &MPCJSON{}
	err = json.Unmarshal([]byte(jsonParams), cJSON)
	if err != nil {
		return err
	}

	mpc, err := assets.NewMPC()
	if err != nil {
		return err
	}
	payload, err := mpc.Payload()
	if err != nil {
		return err
	}

	payload.Type = protobuffer.PBMPCType(cJSON.Type)

	payload.Address = []byte(cJSON.Address)
	payload.Signature = []byte(cJSON.Signature)

	assetID, err := hex.DecodeString(cJSON.AssetID)
	if err != nil {
		return errors.New("Invalid AssetID")
	}
	payload.AssetID = assetID
	mpc.AssetKeyFromPayloadHash()

	txid := ""
	if broadcast == true {
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(mpc)
		if code != 0 {
			print(err.Error())
			return errors.Wrap(err, "TX Fails verifications")
		}
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := mpc.SerializeSignedAsset()
	if err != nil {
		return err
	}

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", mpc.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", mpc.CurrentAsset)

	ppResult()

	return nil
}
