package coreobjects

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

var (
	store map[string]proto.Message
)

func Test_TruthTable(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs()
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.AssetKeyFromPayloadHash()
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another based on previous, ie. AnUpdateWallet
	res, err := w1.TruthTable(protobuffer.PBTransferType_settlePush)
	assert.Nil(t, err, "Truth table return should be nil")

	displayRes := fmt.Sprintln("[", strings.Join(res, "], ["), "]")
	assert.Equal(t, displayRes, "[ 0 + t2 + t3 > 1 & p], [t1 + 0 + t3 > 1 & p], [t1 + t2 + 0 > 1 & p], [t1 + t2 + t3 > 1 & p ]\n", "Truth table invalid")

}

func SetupIDDocs() (*IDDoc, *IDDoc, *IDDoc, *IDDoc) {
	store := NewMapstore()

	idP, _ := NewIDDoc("Primary")
	idP.store = store
	idP.Save()

	idT1, _ := NewIDDoc("trustee1")
	idT1.store = store
	idT1.Save()

	idT2, _ := NewIDDoc("trustee2")
	idT2.store = store
	idT2.Save()

	idT3, _ := NewIDDoc("trustee3")
	idT3.store = store
	idT3.Save()

	return idP, idT1, idT2, idT3
}

func Test_RuleAdd(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs()
	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.store = idP.store
	w1.AssetKeyFromPayloadHash()
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another Wallet based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	w2.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
	sigP, _ := w2.SignPayload(idP)
	sigT1, _ := w2.SignPayload(idT1)
	sigT2, _ := w2.SignPayload(idT2)
	sigT3, _ := w2.SignPayload(idT3)

	//Everything is sign by the Trustees & Principal

	// //Pass correct
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: nil},
	}
	validTransfer1, _ := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.True(t, validTransfer1, "Transfer should be valid")

	//Fail not enough threshold
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: nil},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: nil},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Fail no principal
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: nil},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Pass too many correct
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer should be valid")
}

func Test_AggregationAndVerify(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs()

	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.store = idP.store
	w1.AssetKeyFromPayloadHash()
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another Wallet based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	w2.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
	sigP, _ := w2.SignPayload(idP)
	sigT1, _ := w2.SignPayload(idT1)
	sigT2, _ := w2.SignPayload(idT2)
	sigT3, _ := w2.SignPayload(idT3)

	//Everything is sign by the Trustees & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}

	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer is invalid boolean doesn't return true")

	//Build Aggregated Signature and insert into the releavnt Wallet Fields
	err = w2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check wallet2 validatity based on previous Version
	verify, err := w2.FullVerify(w2.previousAsset)
	assert.True(t, verify, "Verify should be True")
	assert.Nil(t, err, "Error should be nil")

}

func Test_AggregationAndVerifyFailingTransfer(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs()

	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.store = idP.store
	w1.AssetKeyFromPayloadHash()
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another Wallet based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	w2.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
	sigP, _ := w2.SignPayload(idP)
	sigT1, _ := w2.SignPayload(idT1)

	//Everything is sign by the Trustees & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
	}

	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Build Aggregated Signature and insert into the releavnt Wallet Fields
	err = w2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check wallet2 validatity based on previous Version
	verify, err := w2.FullVerify(w2.previousAsset)
	assert.False(t, verify, "Verify should be False")
	assert.NotNil(t, err, "Error should describe the failure")

}
