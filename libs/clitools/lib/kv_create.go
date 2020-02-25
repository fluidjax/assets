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

func (cliTool *CLITool) CreateKVJSON(jsonParams string, broadcast bool) (err error) {
	cKVJSON := &CreateKVJSON{}
	err = json.Unmarshal([]byte(jsonParams), cKVJSON)
	if err != nil {
		return err
	}
	seedHex := cKVJSON.Ownerseed
	iddoc, err := cliTool.GetIDDocForSeed(seedHex)
	if err != nil {
		return err
	}

	kv, err := assets.NewKVAsset(iddoc, protobuffer.PBKVAssetType(cKVJSON.KVAssetType))
	kv.SetKeyString(cKVJSON.AssetID)
	kv.DataStore = cliTool.NodeConn

	//add keys
	for _, pair := range cKVJSON.KV {
		key := pair.Key
		value := pair.Value
		kv.SetKV(key, []byte(value))
	}

	if err != nil {
		return err
	}

	var truths []string
	for _, trans := range cKVJSON.Transfer {

		binParticipants := map[string][]byte{}
		for _, v := range trans.Participants {
			binVal, err := hex.DecodeString(v.ID)
			if err != nil {
				return err
			}
			binParticipants[v.Name] = binVal
		}
		transferType := protobuffer.PBTransferType(trans.TransferType)
		kv.AddTransfer(transferType, trans.Expression, &binParticipants, trans.Description)
		truthTable, err := kv.TruthTable(transferType)
		if err != nil {
			return err
		}

		for _, v := range truthTable {
			x := fmt.Sprintf("%d:%s", trans.TransferType, v)
			truths = append(truths, base64.StdEncoding.EncodeToString([]byte(x)))
		}

	}

	kv.Sign(iddoc)

	txid := ""
	if broadcast == true {
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(kv)
		if code != 0 {
			print(code)
			print(err.Error())
			return errors.Wrap(err, "TX Fails verifications")
		}
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := kv.SerializeSignedAsset()
	if err != nil {
		return err
	}

	res["truthtable"] = truths

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", kv.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", kv.CurrentAsset)

	ppResult()
	return nil
}
