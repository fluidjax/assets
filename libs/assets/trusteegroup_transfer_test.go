package assets

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

func Test_TrusteeGroupTruthTable(t *testing.T) {
	idInitiator, idT1, idT2, idT3 := SetupIDDocs()
	expression := "t1 + t2 + t3 > 1 "
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}
	t1, _ := NewTrusteeGroup(idInitiator)
	t1.AddTransfer(protobuffer.PBTransferType_transferPush, expression, participants)

	//Create another based on previous, ie. AnUpdateTrusteeGroup
	res, err := t1.TruthTable(protobuffer.PBTransferType_transferPush)
	assert.Nil(t, err, "Truth table return should be nil")

	displayRes := fmt.Sprintln("[", strings.Join(res, "], ["), "]")
	assert.Equal(t, displayRes, "[ 0 + t2 + t3 > 1 ], [t1 + 0 + t3 > 1 ], [t1 + t2 + 0 > 1 ], [t1 + t2 + t3 > 1  ]\n", "Truth table invalid")

}

func Test_TrusteeGroupRuleAdd(t *testing.T) {
	idInitiator, idT1, idT2, idT3 := SetupIDDocs()
	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 "
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	t1, _ := NewTrusteeGroup(idInitiator)
	t1.store = idInitiator.store
	t1.AddTransfer(protobuffer.PBTransferType_transferPush, expression, participants)

	//Create another TrusteeGroup based on previous, ie. AnUpdateTrusteeGroup
	t2, _ := NewUpdateTrusteeGroup(t1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	t2.Asset.TransferType = protobuffer.PBTransferType_transferPush

	//Generate Signatures for each Participant - note they are signing the new TrusteeGroup with the TransferType set!
	sigT1, _ := t2.SignPayload(idT1)
	sigT2, _ := t2.SignPayload(idT2)
	sigT3, _ := t2.SignPayload(idT3)

	//Everything is sign by the Trustees & Principal

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

func Test_TrusteeGroupAggregationAndVerify(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs()

	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1"
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	t1, _ := NewTrusteeGroup(idP)
	t1.store = idP.store
	t1.AddTransfer(protobuffer.PBTransferType_transferPush, expression, participants)

	//Create another TrusteeGroup based on previous, ie. AnUpdateTrusteeGroup
	t2, _ := NewUpdateTrusteeGroup(t1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	t2.Asset.TransferType = protobuffer.PBTransferType_transferPush

	//Generate Signatures for each Participant - note they are signing the new TrusteeGroup with the TransferType set!
	sigT1, _ := t2.SignPayload(idT1)
	sigT2, _ := t2.SignPayload(idT2)
	sigT3, _ := t2.SignPayload(idT3)

	//Everything is sign by the Trustees & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}

	validTransfer1, err := t2.IsValidTransfer(protobuffer.PBTransferType_transferPush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer is invalid boolean doesn't return true")

	//Build Aggregated Signature and insert into the releavnt TrusteeGroup Fields
	err = t2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check trusteeGroup2 validatity based on previous Version
	verify, err := t2.FullVerify(t2.previousAsset)
	assert.True(t, verify, "Verify should be True")
	assert.Nil(t, err, "Error should be nil")

}

func Test_TrusteeGroupAggregationAndVerifyFailingTransfer(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs()

	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1"
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	t1, _ := NewTrusteeGroup(idP)
	t1.store = idP.store
	t1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another TrusteeGroup based on previous, ie. AnUpdateTrusteeGroup
	t2, _ := NewUpdateTrusteeGroup(t1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	t2.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new TrusteeGroup with the TransferType set!
	sigP, _ := t2.SignPayload(idP)
	sigT1, _ := t2.SignPayload(idT1)

	//Everything is sign by the Trustees & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
	}

	validTransfer1, err := t2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Build Aggregated Signature and insert into the releavnt TrusteeGroup Fields
	err = t2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check trusteeGroup2 validatity based on previous Version
	verify, err := t2.FullVerify(t2.previousAsset)
	assert.False(t, verify, "Verify should be False")
	assert.NotNil(t, err, "Error should describe the failure")

}
