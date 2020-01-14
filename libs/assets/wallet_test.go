package assets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Wallet_Signing(t *testing.T) {
	testName := "ABC!23"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	i.store = NewMapstore()
	i.Save()
	w, err := NewWallet(i)


	sig1, err := w.SignAsset(i)
	payload, err := w.SerializeAsset()
	sig2, err := Sign(payload, i)
	res := bytes.Compare(sig1, sig2)
	assert.True(t, res == 0, "Compare should be 0")

	verify, err := w.VerifyAsset(sig1, i)
	assert.True(t, verify, "Verify should be true")

	verify2, err := Verify(payload, sig2, i)
	assert.True(t, verify2, "Verify should be true")

}

func Test_Wallet(t *testing.T) {
	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	i.store = NewMapstore()
	i.Save()
	w, err := NewWallet(i)
	walletContents, _ := w.Payload()
	walletContents.Description = testDescription
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

func Test_Wallet_Empty(t *testing.T) {
	w, err := NewWallet(nil)
	assert.NotNil(t, err, "Error should not be nil")
	_, err = w.Payload()
	assert.NotNil(t, err, "Error should not be nil")
}
