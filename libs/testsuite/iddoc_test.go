package testsuite

import (
	"encoding/hex"
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/stretchr/testify/assert"
)

func Test_IDDoc(t *testing.T) {
	StartTestChain()

	i := BuildTestIDDoc(t)
	txid, chainErr := nc.PostTx(i)
	assert.Nil(t, chainErr, "Error should be nil")
	assert.NotEqual(t, txid, "", "TxID should not be empty")

	i = BuildTestIDDoc(t)
	i.CurrentAsset.Signature = nil

	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")

	// txid, chainErr := nc.PostTx(i)

	// c := assets.TransactionCode(code)
	// println(txid)
	// println(c.String())
	// println(err)
	// // //serialize the whole transaction
	// // serializedTX1, err := i.SerializeSignedAsset()
	// // assert.Nil(t, err, "Error should be nil")
	// // base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX1)
	// // _, err = qredochain.PostTx(base64EncodedTX, "127.0.0.1:26657")
	// // assert.Nil(t, err, "Error should be nil")

	// dump(i.Key())

}

func BuildTestIDDoc(t *testing.T) *assets.IDDoc {
	randBytes, err := cryptowallet.RandomBytes(32)
	if err != nil {
		panic("Fail to create random string")
	}

	i, err := assets.NewIDDoc(hex.EncodeToString(randBytes))
	i.DataStore = app
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)
	return i
}
