package coreobjects

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

var (
	store map[string]proto.Message
)

func Test_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")
	i.SetTestKey()
	i.Sign()
	i.Description()
	res, err := i.Verify()
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	i.Dump()
}

func Test_Serialize_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")

	data, err := i.PayloadSerialize()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, data, "Result should not be nil")

	i.Signature.Asset = nil
	data, err = i.PayloadSerialize()
	assert.NotNil(t, err, "Error should not be nil")
}

func Test_Save_Load(t *testing.T) {
	testName := "ABC!23"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.SetTestKey()
	i.Sign()
	i.store = NewMapstore()
	key := i.key
	i.Save()
	retrieved, err := Load(i.store, key)
	assert.Nil(t, err, "Error should be nil")
	print(retrieved)
	iddoc := retrieved.Asset.GetIddoc()
	assert.Equal(t, testName, iddoc.AuthenticationReference, "Load/Save failed")
}

func Test_Wallet(t *testing.T) {
	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.SetTestKey()
	i.Sign()
	i.store = NewMapstore()
	i.Save()

	w, err := NewWallet(i)
	walletContents := w.AssetPayload()
	walletContents.Description = testDescription
	w.SetTestKey()
	w.Sign(i)
	assert.NotNil(t, w.Signature.Signature, "Signature is empty")
	res, err := w.Verify(i)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	w.Save()

	retrieved, err := Load(i.store, w.key)
	retrievedWallet := retrieved.Asset.GetWallet()

	assert.Equal(t, testDescription, retrievedWallet.Description, "Load/Save failed")
}

func Test_RuleAdd(t *testing.T) {
	idP, _ := NewIDDoc("Primary")
	idT1, _ := NewIDDoc("trustee1")
	idT2, _ := NewIDDoc("trustee2")
	idT3, _ := NewIDDoc("trustee3")

	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.key,
		"t1": idT1.key,
		"t2": idT2.key,
		"t3": idT3.key,
	}

	w1, _ := NewWallet(idP)
	w1.SetTestKey()
	w1.AddTransfer(protobuffer.TransferType_settlePush, expression, participants)
	//w1.Dump()

	//Create another based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//Generate Signatures for each Participant

	sigP, _ := w2.SignPayload(idP)
	sigT1, _ := w2.SignPayload(idT1)
	sigT2, _ := w2.SignPayload(idT2)
	sigT3, _ := w2.SignPayload(idT3)

	res1, _ := w2.VerifyPayload(sigT1, idT1)
	assert.True(t, res1, "Sig fails to verify")
	res2, _ := w2.VerifyPayload(sigT2, idT2)
	assert.True(t, res2, "Sig fails to verify")
	res3, _ := w2.VerifyPayload(sigT3, idT3)
	assert.True(t, res3, "Sig fails to verify")

	//We should have 3 valid sigs.

	//Pass correct
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Signature: sigP},
		SignatureID{IDDoc: idT1, Signature: sigT1},
		SignatureID{IDDoc: idT2, Signature: sigT2},
		SignatureID{IDDoc: idT3, Signature: nil},
	}
	validTransfer1, _ := w2.IsValidTransfer(protobuffer.TransferType_settlePush, transferSignatures1)
	assert.True(t, validTransfer1, "Transfer should be valid")

	//Pass too many correct
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Signature: sigP},
		SignatureID{IDDoc: idT1, Signature: sigT1},
		SignatureID{IDDoc: idT2, Signature: sigT2},
		SignatureID{IDDoc: idT3, Signature: sigT3},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.TransferType_settlePush, transferSignatures1)
	assert.True(t, validTransfer1, "Transfer should be valid")

	//Fail not enough threshold
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Signature: sigP},
		SignatureID{IDDoc: idT1, Signature: nil},
		SignatureID{IDDoc: idT2, Signature: nil},
		SignatureID{IDDoc: idT3, Signature: sigT3},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.TransferType_settlePush, transferSignatures1)
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Fail no principal
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Signature: nil},
		SignatureID{IDDoc: idT1, Signature: sigT1},
		SignatureID{IDDoc: idT2, Signature: sigT2},
		SignatureID{IDDoc: idT3, Signature: sigT3},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.TransferType_settlePush, transferSignatures1)
	assert.False(t, validTransfer1, "Transfer should be invalid")
}

func Test_TruthTable(t *testing.T) {

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

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.key,
		"t1": idT1.key,
		"t2": idT2.key,
		"t3": idT3.key,
	}

	w1, _ := NewWallet(idP)
	w1.SetTestKey()
	w1.AddTransfer(protobuffer.TransferType_settlePush, expression, participants)
	//w1.Dump()

	//Create another based on previous, ie. AnUpdateWallet

	tt, err := w1.TruthTable(protobuffer.TransferType_settlePush)
	assert.Nil(t, err, "Truth table return should be nil")

	print(tt)
}
