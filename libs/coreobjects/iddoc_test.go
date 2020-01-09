package coreobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")
	i.Sign()
	res, err := i.Verify()
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	i.Dump()
}

func Test_Serialize_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")

	data, err := i.serializePayload()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, data, "Result should not be nil")

	i.PBSignedAsset.Asset = nil
	data, err = i.serializePayload()
	assert.NotNil(t, err, "Error should not be nil")
}

func Test_Save_Load(t *testing.T) {
	testName := "ABC!23"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign()
	i.store = NewMapstore()
	key := i.Key()
	i.Save()
	retrieved, err := Load(i.store, key)
	assert.Nil(t, err, "Error should be nil")
	print(retrieved)
	iddoc := retrieved.Asset.GetIddoc()
	assert.Equal(t, testName, iddoc.AuthenticationReference, "Load/Save failed")
}
