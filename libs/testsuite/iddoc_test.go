package testsuite

import (
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/stretchr/testify/assert"
)

func Test_IDDoc(t *testing.T) {
	StartTestChain()

	i, err := assets.NewIDDoc("1st Ref")
	i.DataStore = app
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	i.CurrentAsset.Signature = nil

	txid, code, err := nc.PostTx(i)

	println(txid)
	println(code)
	println(err)
	// //serialize the whole transaction
	// serializedTX1, err := i.SerializeSignedAsset()
	// assert.Nil(t, err, "Error should be nil")
	// base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX1)
	// _, err = qredochain.PostTx(base64EncodedTX, "127.0.0.1:26657")
	// assert.Nil(t, err, "Error should be nil")

	dump(i.Key())

}
