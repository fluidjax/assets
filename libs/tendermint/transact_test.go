package tendermint

import (
	"encoding/base64"
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/datastore"
	"github.com/qredo/assets/libs/logger"
	"github.com/stretchr/testify/assert"
)

func Test_PostTX(t *testing.T) {
	i, err := assets.NewIDDoc("")
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	//serialize the whole transaction
	serializedTX, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX)

	txid, err := PostTx(base64EncodedTX, "127.0.0.1:26657")
	assert.Nil(t, err, "Error should be nil")

	print(txid)
}

func Test_NodeConnector(t *testing.T) {

	dsBackend, err := datastore.NewBoltBackend("datastore.dat")
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, dsBackend, "store should not be nil")

	store, err := datastore.NewStore(datastore.WithBackend(dsBackend), datastore.WithCodec(datastore.NewGOBCodec()))
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, store, "store should not be nil")

	logger, err := logger.NewLogger("text", "info")
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, logger, "logger should not be nil")

	nc, err := NewNodeConnector("127.0.0.1:26657", "NODEID", store, logger)
	assert.NotNil(t, nc, "tmConnector should not be nil")
	assert.Nil(t, err, "Error should be nil")

	//Build IDDoc
	i, err := assets.NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)

	nc.PostTx(i)

}
