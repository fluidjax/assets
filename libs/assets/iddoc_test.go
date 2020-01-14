package assets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IDDoc_Misc(t *testing.T) {
	i, err := NewIDDoc("")
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(nil)
	assert.NotNil(t, err, "Error should not be nil")
	_, err = i.Verify(nil)
	assert.NotNil(t, err, "Error should not be nil")
}

func Test_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	res, err := i.Verify(i)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
}

func Test_Serialize_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")

	data, err := i.SerializeAsset()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, data, "Result should not be nil")

	i.currentAsset.Asset = nil
	data, err = i.SerializeAsset()
	assert.NotNil(t, err, "Error should not be nil")
}

func Test_Save_Load(t *testing.T) {
	testName := "ABC!23"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	i.store = NewMapstore()
	key := i.Key()
	i.Save()
	retrieved, err := Load(i.store, key)
	assert.Nil(t, err, "Error should be nil")
	print(retrieved)
	iddoc := retrieved.Asset.GetIddoc()
	assert.Equal(t, testName, iddoc.AuthenticationReference, "Load/Save failed")
}
