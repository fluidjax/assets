package qc

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
)

func (cliTool *CLITool) CreateGroupJSON(jsonParams string, broadcast bool) (err error) {
	cGroupJSON := &CreateGroupJSON{}
	err = json.Unmarshal([]byte(jsonParams), cGroupJSON)
	if err != nil {
		return err
	}
	seedHex := cGroupJSON.Ownerseed
	iddoc, err := cliTool.GetIDDocForSeed(seedHex)
	if err != nil {
		return err
	}

	group, err := assets.NewGroup(iddoc, protobuffer.PBGroupType(cGroupJSON.Group.Type))
	group.DataStore = cliTool.NodeConn

	if err != nil {
		return err
	}

	//Add in the Group Payload info
	payload := group.Payload()
	payload.Type = protobuffer.PBGroupType(cGroupJSON.Group.Type)
	payload.Description = cGroupJSON.Group.Description

	payload.GroupFields = buildKV(&cGroupJSON.Group.GroupFields)
	payload.Participants = buildKV(&cGroupJSON.Group.Participants)

	var truths []string
	for _, trans := range cGroupJSON.Transfer {

		binParticipants := map[string][]byte{}
		for _, v := range trans.Participants {
			binVal, err := hex.DecodeString(v.ID)
			if err != nil {
				return err
			}
			binParticipants[v.Name] = binVal
		}
		transferType := protobuffer.PBTransferType(trans.TransferType)
		group.AddTransfer(transferType, trans.Expression, &binParticipants, trans.Description)
		truthTable, err := group.TruthTable(transferType)
		if err != nil {
			return err
		}

		for _, v := range truthTable {
			x := fmt.Sprintf("%d:%s", trans.TransferType, v)
			truths = append(truths, base64.StdEncoding.EncodeToString([]byte(x)))
		}
	}

	group.Sign(iddoc)

	txid := ""
	if broadcast == true {
		var code assets.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(group)
		if code != 0 {
			print(code)
			return errors.Wrap(err, "TX Fails verifications")
		}
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := group.SerializeSignedAsset()
	if err != nil {
		return err
	}

	res["truthtable"] = truths

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", group.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", group.CurrentAsset)

	ppResult()
	return nil
}

func (cliTool *CLITool) PrepareGroupUpdateWithJSON(jsonParams string) (err error) {
	uwJSON := &GroupUpdatePayload{}
	err = json.Unmarshal([]byte(jsonParams), uwJSON)
	if err != nil {
		return err
	}
	updatedGroup, err := cliTool.groupFromGroupUpdateJSON(uwJSON)
	if err != nil {
		return err
	}

	//Updated Group complete, return for signing
	msg, err := updatedGroup.SerializeAsset()
	if err != nil {
		return errors.New("Failed to serialize payload")
	}

	addResultBinaryItem("serializedUpdate", msg)
	ppResult()
	return nil
}

func (cliTool *CLITool) AggregateGroupSign(jsonParams string, broadcast bool) (err error) {
	//Decode the JSON
	agJSON := &GroupUpdate{}
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

	//Rebuild the Group from the TX supplied
	updatedGroup, err := cliTool.groupFromGroupUpdateJSON(&agJSON.GroupUpdatePayload)

	err = updatedGroup.AggregatedSign(transferSignatures)
	if err != nil {
		return err
	}

	verify, err := updatedGroup.FullVerify()

	if verify == false {
		return errors.New("Error failed to verify final update Group transaction")
	}

	txid := ""
	if broadcast == true {
		var code assets.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(updatedGroup)
		if code != 0 {
			return errors.Wrap(err, "TX Fails verifications")
		}
		if err != nil {
			return err
		}
	}

	serializedSignedAsset, err := updatedGroup.SerializeSignedAsset()
	if err != nil {
		return err
	}

	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", updatedGroup.Key())
	addResultBinaryItem("serializedSignedAsset", serializedSignedAsset)
	addResultSignedAsset("object", updatedGroup.CurrentAsset)
	ppResult()
	return
}
