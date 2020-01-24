package qredochain

import (
	"encoding/base64"
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/tendermint"
	"github.com/stretchr/testify/assert"
)

func Test_LoadSave(t *testing.T) {
	i, err := assets.NewIDDoc("1st Ref")
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	//serialize the whole transaction
	serializedTX1, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX1)
	txid, err := tendermint.PostTx(base64EncodedTX, "127.0.0.1:26657")
	print(txid)

	

}
