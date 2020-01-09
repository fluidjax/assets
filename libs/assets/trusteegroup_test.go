package assets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TrusteeGroup(t *testing.T) {
	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign()
	i.store = NewMapstore()
	i.Save()

	w, err := NewTrusteeGroup(i)
	trusteegroupContents := w.Payload()
	trusteegroupContents.Description = testDescription
	w.Sign(i)
	assert.NotNil(t, w.PBSignedAsset.Signature, "Signature is empty")
	res, err := w.Verify(i)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	w.Save()

	retrieved, err := Load(i.store, w.Key())
	retrievedTrusteeGroup := retrieved.Asset.GetTrusteeGroup()

	assert.Equal(t, testDescription, retrievedTrusteeGroup.Description, "Load/Save failed")
}
