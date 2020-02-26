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

func (cliTool *CLITool) groupFromGroupUpdateJSON(signedUpdate *GroupUpdatePayload) (*assets.Group, error) {
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

	//Get the Existing Group
	existingGroupKey, err := hex.DecodeString(signedUpdate.ExistingGroupAssetID)
	if err != nil {
		return nil, err
	}

	originalGroup, err := assets.LoadGroup(cliTool.NodeConn, existingGroupKey)
	if err != nil {
		return nil, err
	}

	originalGroup.DataStore = cliTool.NodeConn

	//Make New Group based on Existing
	updatedGroup, err := assets.NewUpdateGroup(originalGroup, newOwnerIDDoc)
	if err != nil {
		return nil, err
	}

	//add keys
	payload := updatedGroup.Payload()
	payload.Type = protobuffer.PBGroupType(signedUpdate.Group.Type)
	payload.Description = signedUpdate.Group.Description

	payload.GroupFields = buildKV(&signedUpdate.Group.GroupFields)
	payload.Participants = buildKV(&signedUpdate.Group.Participants)

	var truths []string
	for _, trans := range signedUpdate.Transfer {

		binParticipants := map[string][]byte{}
		for _, v := range trans.Participants {
			binVal, err := hex.DecodeString(v.ID)
			if err != nil {
				return nil, err
			}
			binParticipants[v.Name] = binVal
		}
		transferType := protobuffer.PBTransferType(trans.TransferType)
		updatedGroup.AddTransfer(transferType, trans.Expression, &binParticipants, trans.Description)
		truthTable, err := updatedGroup.TruthTable(transferType)
		if err != nil {
			return nil, err
		}

		for _, v := range truthTable {
			x := fmt.Sprintf("%d:%s", trans.TransferType, v)
			truths = append(truths, base64.StdEncoding.EncodeToString([]byte(x)))
		}
	}

	updatedGroup.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType(signedUpdate.TransferType)

	return updatedGroup, nil
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
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(updatedGroup)
		if code != 0 {
			print(err.Error())
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
