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

func (cliTool *CLITool) PrepareKVUpdateWithJSON(jsonParams string) (err error) {
	//Load existing kv from AssetID
	//Load all the IDDocs

	//kv from JSON
	uwJSON := &KVUpdatePayload{}
	err = json.Unmarshal([]byte(jsonParams), uwJSON)
	if err != nil {
		return err
	}
	updatedKV, err := cliTool.kVFromKVUpdateJSON(uwJSON)
	if err != nil {
		return err
	}

	//Updated kv complete, return for signing
	msg, err := updatedKV.SerializeAsset()
	if err != nil {
		return errors.New("Failed to serialize payload")
	}

	addResultBinaryItem("serializedUpdate", msg)
	ppResult()
	return nil
}

func (cliTool *CLITool) AggregateKVSign(jsonParams string, broadcast bool) (err error) {
	//Decode the JSON
	agJSON := &KVUpdate{}
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

	//Rebuild the KV from the TX supplied
	updatedKV, err := cliTool.kVFromKVUpdateJSON(&agJSON.KVUpdatePayload)

	err = updatedKV.AggregatedSign(transferSignatures)
	if err != nil {
		return err
	}

	verify, err := updatedKV.FullVerify()

	if verify == false {
		return errors.New("Error failed to verify final update kv transaction")
	}

	txid := ""
	if broadcast == true {
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(updatedKV)
		if code != 0 {
			return errors.Wrap(err, "TX Fails verifications")
		}
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := updatedKV.SerializeSignedAsset()
	if err != nil {
		return err
	}

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", updatedKV.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", updatedKV.CurrentAsset)
	ppResult()
	return
}
