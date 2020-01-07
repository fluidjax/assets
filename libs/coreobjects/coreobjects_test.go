package coreobjects

import (
	"testing"

	"github.com/gogo/protobuf/proto"
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
	res, err := i.Verify()
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
}

func Test_Serialize_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")

	data, err := i.PayloadSerialize()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, data, "Result should not be nil")

	i.Signature.SignatureAsset = nil
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

	//i2 = MakeIDDoc(retrieved)

	assert.Nil(t, err, "Error should be nil")
	print(retrieved)
	assdec := retrieved.GetDeclaration()
	iddoc := assdec.GetIddoc()

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
	assdec := retrieved.GetDeclaration()
	retrievedWallet := assdec.GetWallet()

	assert.Equal(t, testDescription, retrievedWallet.Description, "Load/Save failed")
}
