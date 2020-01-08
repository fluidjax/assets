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

	//idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 & p"
	participants := map[string][]byte{
		"p":  idP.key,
		"t1": idT1.key,
		"t2": idT2.key,
		"t3": idT3.key,
	}

	wDec, _ := NewWallet(idP)
	wDec.AddTransfer(protobuffer.TransferType_settlePush, expression, participants)
	wDec.Dump()

	//Create an Update
	// wUpdate, _ := NewUpdateWallet(idP, wDec)

	// //Generate Signatures for each Participant

	// //Verify Signature for Each Particpiant

	// transferSignatures := []SignatureID{
	// 	SignatureID{IDDocID: idP.key, Signature: nil},
	// 	SignatureID{IDDocID: idT1.key, Signature: nil},
	// 	SignatureID{IDDocID: idT2.key, Signature: nil},
	// 	SignatureID{IDDocID: idT3.key, Signature: nil},
	// 	SignatureID{IDDocID: idT4.key, Signature: nil},
	// }

}
