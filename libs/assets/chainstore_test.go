package assets

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/qredo/assets/libs/store"
	"github.com/stretchr/testify/assert"
)

func Test_IDDocPostTX(t *testing.T) {
	store := store.NewChainstore()

	i, err := NewIDDoc("1st Ref")
	i.DataStore = store
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	//serialize the whole transaction
	serializedTX1, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX1)

	txid, err := PostTx(base64EncodedTX, "127.0.0.1:26657")
	assert.Nil(t, err, "Error should be nil")
	txid, err = PostTx(base64EncodedTX, "127.0.0.1:26657")

	//Change 1 field and post again
	payload, err := i.Payload()
	payload.AuthenticationReference = "Different ref"
	err = i.Sign(i)
	serializedTX2, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	assert.False(t, bytes.Compare(serializedTX1, serializedTX2) == 0, "Serialize TX should be different")
	base64EncodedTX = base64.StdEncoding.EncodeToString(serializedTX2)
	txid, err = PostTx(base64EncodedTX, "127.0.0.1:26657")
	print(txid)
}

func Test_WalletPostTX(t *testing.T) {
	store := store.NewChainstore()

	i, err := NewIDDoc("1st Ref")
	i.DataStore = store
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	w, err := NewWallet(i, "BTC")
	wallet, err := w.Payload()
	wallet.SpentBalance = 100

	assert.Nil(t, err, "Error should be nil")

	//serialize the whole transaction
	serializedTX1, err := w.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX1)

	_, err = PostTx(base64EncodedTX, "127.0.0.1:26657")
	assert.Nil(t, err, "Error should be nil")

	//update wallet
	w2, _ := NewUpdateWallet(w, i)
	//serialize the whole transaction
	serializedTX2, err := w2.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX2 := base64.StdEncoding.EncodeToString(serializedTX2)

	_, err = PostTx(base64EncodedTX2, "127.0.0.1:26657")
	assert.Nil(t, err, "Error should be nil")

}
