package assets

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

func Test_GroupTruthTable(t *testing.T) {
	store := NewMapstore()
	idInitiator, idT1, idT2, idT3 := SetupIDDocs(store)
	expression := "t1 + t2 + t3 > 1 "
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	t1, _ := NewGroup(idInitiator, protobuffer.PBGroupType_trusteeGroup)
	t1.AddTransfer(protobuffer.PBTransferType_transferPush, expression, participants)

	//Create another based on previous, ie. AnUpdateGroup
	res, err := t1.TruthTable(protobuffer.PBTransferType_transferPush)
	assert.Nil(t, err, "Truth table return should be nil")

	displayRes := fmt.Sprintln("[", strings.Join(res, "], ["), "]")
	assert.Equal(t, displayRes, "[ 0 + t2 + t3 > 1 ], [t1 + 0 + t3 > 1 ], [t1 + t2 + 0 > 1 ], [t1 + t2 + t3 > 1  ]\n", "Truth table invalid")
}

func Test_GroupRuleAdd(t *testing.T) {
	store := NewMapstore()
	idInitiator, idT1, idT2, idT3 := SetupIDDocs(store)
	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 "
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	t1, _ := NewGroup(idInitiator, protobuffer.PBGroupType_trusteeGroup)
	t1.store = idInitiator.store
	t1.AddTransfer(protobuffer.PBTransferType_transferPush, expression, participants)

	//Create another Group based on previous, ie. AnUpdateGroup
	t2, _ := NewUpdateGroup(t1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	t2.Asset.TransferType = protobuffer.PBTransferType_transferPush

	//Generate Signatures for each Participant - note they are signing the new Group with the TransferType set!
	sigT1, _ := t2.SignPayload(idT1)
	sigT2, _ := t2.SignPayload(idT2)
	sigT3, _ := t2.SignPayload(idT3)

	//Everything is sign by the s & Principal

	// //Pass correct
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: nil},
	}
	validTransfer1, _ := t2.IsValidTransfer(protobuffer.PBTransferType_transferPush, transferSignatures1)
	assert.True(t, validTransfer1, "Transfer should be valid")

	//Fail not enough threshold
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: nil},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: nil},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, _ = t2.IsValidTransfer(protobuffer.PBTransferType_transferPush, transferSignatures1)
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Pass too many correct
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, err := t2.IsValidTransfer(protobuffer.PBTransferType_transferPush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer should be valid")
}

func SetupTrusteeGroup(store *Mapstore) (*IDDoc, *IDDoc, *IDDoc, *Group) {
	tgInitiator, tgT1, tgT2, tgT3 := SetupIDDocs(store)

	w, _ := NewGroup(tgInitiator, protobuffer.PBGroupType_trusteeGroup)

	expression := "x1 + x2 + x3 > 1 "
	participants := &map[string][]byte{
		"x1": tgT1.Key(),
		"x2": tgT2.Key(),
		"x3": tgT3.Key(),
	}
	w.ConfigureGroup(expression, participants, "Trsutee Group Test Description")
	w.store = store
	w.Save()

	return tgT1, tgT2, tgT3, w
}

func Test_GroupAggregationAndVerify(t *testing.T) {
	store := NewMapstore()
	idP, idT1, idT2, idT3 := SetupIDDocs(store)

	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1"
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	t1, _ := NewGroup(idP, protobuffer.PBGroupType_trusteeGroup)
	t1.store = idP.store
	t1.AddTransfer(protobuffer.PBTransferType_transferPush, expression, participants)

	//Create another Group based on previous, ie. AnUpdateGroup
	t2, _ := NewUpdateGroup(t1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	t2.Asset.TransferType = protobuffer.PBTransferType_transferPush

	//Generate Signatures for each Participant - note they are signing the new Group with the TransferType set!
	sigT1, _ := t2.SignPayload(idT1)
	sigT2, _ := t2.SignPayload(idT2)
	sigT3, _ := t2.SignPayload(idT3)

	//Everything is sign by the s & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}

	validTransfer1, err := t2.IsValidTransfer(protobuffer.PBTransferType_transferPush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer is invalid boolean doesn't return true")

	//Build Aggregated Signature and insert into the releavnt Group Fields
	err = t2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check Group2 validatity based on previous Version
	verify, err := t2.FullVerify(t2.previousAsset)
	assert.True(t, verify, "Verify should be True")
	assert.Nil(t, err, "Error should be nil")

}

func Test_Recusion_GroupAggregationAndVerify(t *testing.T) {
	store := NewMapstore()
	idP, _, idT2, idT3 := SetupIDDocs(store)
	idX1, idX2, _, Group := SetupTrusteeGroup(store)
	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "tg1 + t2 + t3 > 1"
	participants := &map[string][]byte{
		"tg1": Group.Key(),
		"t2":  idT2.Key(),
		"t3":  idT3.Key(),
	}

	t1, _ := NewGroup(idP, protobuffer.PBGroupType_trusteeGroup)
	t1.store = idP.store
	t1.AddTransfer(protobuffer.PBTransferType_transferPush, expression, participants)

	//Create another Group based on previous, ie. AnUpdateGroup
	t2, _ := NewUpdateGroup(t1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	t2.Asset.TransferType = protobuffer.PBTransferType_transferPush

	//Generate Signatures for each Participant - note they are signing the new Group with the TransferType set!
	//sigT1, _ := t2.SignPayload(idT1)
	sigT2, _ := t2.SignPayload(idT2)
	sigT3, _ := t2.SignPayload(idT3)

	sigX1, _ := t2.SignPayload(idX1)
	sigX2, _ := t2.SignPayload(idX2)
	//sigX3, _ := t2.SignPayload(idX3)

	//Everything is sign by the s & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idX1, Abbreviation: "tg1.x1", Signature: sigX1},
		SignatureID{IDDoc: idX2, Abbreviation: "tg1.x2", Signature: sigX2},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}

	validTransfer1, err := t2.IsValidTransfer(protobuffer.PBTransferType_transferPush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer is invalid boolean doesn't return true")

	//Build Aggregated Signature and insert into the releavnt Group Fields
	err = t2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check Group2 validatity based on previous Version
	verify, err := t2.FullVerify(t2.previousAsset)
	assert.True(t, verify, "Verify should be True")
	assert.Nil(t, err, "Error should be nil")
}

func Test_GroupAggregationAndVerifyFailingTransfer(t *testing.T) {
	store := NewMapstore()
	idP, idT1, idT2, idT3 := SetupIDDocs(store)

	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1"
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	t1, _ := NewGroup(idP, protobuffer.PBGroupType_trusteeGroup)
	t1.store = idP.store
	t1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another Group based on previous, ie. AnUpdateGroup
	t2, _ := NewUpdateGroup(t1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	t2.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new Group with the TransferType set!
	sigP, _ := t2.SignPayload(idP)
	sigT1, _ := t2.SignPayload(idT1)

	//Everything is sign by the s & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
	}

	validTransfer1, err := t2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Build Aggregated Signature and insert into the releavnt Group Fields
	err = t2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check Group2 validatity based on previous Version
	verify, err := t2.FullVerify(t2.previousAsset)
	assert.False(t, verify, "Verify should be False")
	assert.NotNil(t, err, "Error should describe the failure")
}

func Test_GroupTransferParser(t *testing.T) {
	store := NewMapstore()
	idP, idT1, idT2, idT3 := SetupIDDocs(store)
	_, groupMember1, groupMember2, groupMember3 := SetupIDDocs(store)
	_, groupMember4, groupMember5, _ := SetupIDDocs(store)

	expression := "t1 + t2 + t3 > 1"
	idNewOwner, _ := NewIDDoc("NewOwner")

	groupMembers := &map[string][]byte{
		"g1": groupMember1.Key(),
		"g2": groupMember2.Key(),
		"g3": groupMember3.Key(),
	}

	//Remove 1, add 4 & 5
	updatedGroupMembers := &map[string][]byte{
		"g1": groupMember2.Key(),
		"g2": groupMember3.Key(),
		"g3": groupMember4.Key(),
		"g4": groupMember5.Key(),
	}

	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	t1, _ := NewGroup(idP, protobuffer.PBGroupType_default)

	groupPayload := t1.Payload()
	groupPayload.Participants = *groupMembers

	t1.store = idP.store
	t1.AddTransfer(protobuffer.PBTransferType_transferPush, expression, participants)

	//Create another Group based on previous, ie. AnUpdateGroup
	t2, _ := NewUpdateGroup(t1, idNewOwner)
	t2Payload := t2.Payload()
	t2Payload.Participants = *updatedGroupMembers
	t2.Asset.TransferType = protobuffer.PBTransferType_transferPush

	unchanged, added, deleted, err := t2.ParseChanges()

	assert.True(t, len(deleted) == 1, "Deleted is wrong length")
	assert.True(t, contains(deleted, groupMember1.Key()), "Invalid Deleted Parsing")

	assert.True(t, len(added) == 2, "Added is wrong length")
	assert.False(t, contains(added, groupMember1.Key()), "Invalid Added Parsing x1")
	assert.False(t, contains(added, groupMember2.Key()), "Invalid Added Parsing x2")
	assert.False(t, contains(added, groupMember3.Key()), "Invalid Added Parsing x3")

	assert.True(t, contains(added, groupMember4.Key()), "Invalid Added Parsing 1")
	assert.True(t, contains(added, groupMember5.Key()), "Invalid Added Parsing 2")

	assert.True(t, len(unchanged) == 2, "Unchanged is wrong length")
	assert.True(t, contains(unchanged, groupMember2.Key()), "Invalid Unchanged Parsing 1")
	assert.True(t, contains(unchanged, groupMember3.Key()), "Invalid Unchanged Parsing 2")

	//Generate Signatures for each Participant - note they are signing the new Group with the TransferType set!
	sigT1, _ := t2.SignPayload(idT1)
	sigT2, _ := t2.SignPayload(idT2)
	sigT3, _ := t2.SignPayload(idT3)

	//Everything is sign by the s & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}

	validTransfer1, err := t2.IsValidTransfer(protobuffer.PBTransferType_transferPush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer is invalid boolean doesn't return true")

	//Build Aggregated Signature and insert into the releavnt Group Fields
	err = t2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check Group2 validatity based on previous Version
	verify, err := t2.FullVerify(t2.previousAsset)
	assert.True(t, verify, "Verify should be True")
	assert.Nil(t, err, "Error should be nil")
}
