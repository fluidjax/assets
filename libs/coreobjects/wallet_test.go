package coreobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Wallet(t *testing.T) {
	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.AssetKeyFromPayloadHash()
	i.Sign()
	i.store = NewMapstore()
	i.Save()

	w, err := NewWallet(i)
	walletContents := w.Payload()
	walletContents.Description = testDescription
	w.AssetKeyFromPayloadHash()
	w.Sign(i)
	assert.NotNil(t, w.PBSignedAsset.Signature, "Signature is empty")
	res, err := w.Verify(i)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	w.Save()

	retrieved, err := Load(i.store, w.Key())
	retrievedWallet := retrieved.Asset.GetWallet()

	assert.Equal(t, testDescription, retrievedWallet.Description, "Load/Save failed")
}
