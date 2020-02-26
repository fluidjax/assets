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
		var code qredochain.TransactionCode
		txid, code, err = cliTool.NodeConn.PostTx(group)
		if code != 0 {
			print(code)
			print(err.Error())
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
