package qc

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
)

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

func (cliTool *CLITool) kVFromKVUpdateJSON(signedUpdate *KVUpdatePayload) (*assets.KVAsset, error) {
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

	//Get the Existing KV
	existingKVKey, err := hex.DecodeString(signedUpdate.ExistingKVAssetID)
	if err != nil {
		return nil, err
	}

	originalKV, err := assets.LoadKVAsset(cliTool.NodeConn, existingKVKey)
	if err != nil {
		return nil, err
	}

	//Make New KV based on Existing
	updatedKV, err := assets.NewUpdateKVAsset(originalKV, newOwnerIDDoc)
	if err != nil {
		return nil, err
	}

	//add keys
	for _, pair := range signedUpdate.KV {
		key := pair.Key
		value := pair.Value
		updatedKV.SetKV(key, []byte(value))
	}

	updatedKV.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType(signedUpdate.TransferType)
	updatedKV.DataStore = cliTool.NodeConn

	return updatedKV, nil
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
			print(err.Error())
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
