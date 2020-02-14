package qredochain

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/qredo/assets/libs/logger"
	"github.com/qredo/assets/libs/store"

	"github.com/qredo/assets/libs/assets"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/rpc/core"
)

//Wrap all tests so it starts a test Chain
// func TestMain(m *testing.M) {
// 	StartTestChain()
// 	code := m.Run()
// 	time.Sleep(2 * time.Second)
// 	ShutDown()
// 	os.Exit(code)
// }

func Test_LoadSave(t *testing.T) {
	i, err := assets.NewIDDoc("1st Ref")
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	//serialize the whole transaction
	serializedTX1, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX1)
	txid, err := PostTx(base64EncodedTX, "127.0.0.1:26657")

	assert.Nil(t, err, "Error should be nil", err)
	assert.NotNil(t, txid, "txid should not be nil", err)

}

func Test_ChainPutGet(t *testing.T) {
	//StartTestChain() //included in TestMain()

	i, err := assets.NewIDDoc("1st Ref")
	i.DataStore = app

	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	//serialize the whole transaction
	serializedTX1, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX1)
	_, err = PostTx(base64EncodedTX, "127.0.0.1:26657")
	assert.Nil(t, err, "Error should be nil")

	testData := []byte("Some test data")
	testKey := []byte("TestKey")

	//Make a test transaction
	app.currentBatch = app.db.NewTransaction(true)
	err = app.Save(testKey, testData)

	//testKey, _ = hex.DecodeString("B8F4D3CBFFFFD9D4D12A69AD8236F2A1295B4DCEBE018C7D6307FF7FAABD0CF9")

	assert.Nil(t, err, "Save returned error")
	app.currentBatch.Commit()

	app.currentBatch = app.db.NewTransaction(false)
	retrievedData, err := app.Load(testKey)
	assert.NotNil(t, retrievedData, "Retrieve data is nil")
	assert.True(t, string(retrievedData) == string(testData), "Failed to retrieve data")
	app.currentBatch.Commit()

}

func Test_External_Query(t *testing.T) {
	//Bring up a Node
	nc := StartTestConnectionNode(t)
	defer nc.Stop()

	i, err := assets.NewIDDoc("testdoc")
	i.AddTag("tagkey", []byte("tagvalue"))
	assert.Nil(t, err, "Error should be nil", err)

	txid, errorCode, err := nc.PostTx(i)
	assert.True(t, errorCode == CodeFailVerfication, "Error should not be nil", err)

	i.Sign(i)
	txid, errorCode, err = nc.PostTx(i)
	print("TXID ", txid, "\n")
	print("Key  ", hex.EncodeToString(i.Key()), "\n")

	assert.True(t, errorCode == CodeTypeOK, "Error should be nil", err)

	assert.Nil(t, err, "Error should be nil", err)
	assert.NotNil(t, txid, "TXID shouldnt be nil", err)
}

// Test checks that access to both the Tendermint underlying KV database and
// The Badger QredoChain Consensus database is accessible both
// locally (non RPC, nc connectors)
// and Remotely using http, connectors.
func Test_IDOC(t *testing.T) {
	//Initialize a Node
	nc := StartTestConnectionNode(t)

	//Post TX - Remote (note the 2 second wait for chain to create block)
	i, txid, serializedIDDoc, err := buildTestIDDoc(t, nc)

	assert.Nil(t, err, "Error should be nil", err)

	//Tendermint DB - Local Search
	query := fmt.Sprintf("tx.hash='%v'", txid)
	res, err := core.TxSearch(nil, query, true, 1, 30)
	compareAssets(t, res.Txs[0].Tx, serializedIDDoc, i.Key())

	//Tendermint DB - Remote Search
	rquery := fmt.Sprintf("tx.hash='%s'", txid)
	resq, err := nc.TxSearch(rquery, false, 1, 1)
	compareAssets(t, resq.Txs[0].Tx, serializedIDDoc, i.Key())

	//Badger Consensus DB - Remote search
	txidBytes, _ := hex.DecodeString(txid)
	data, err := nc.TmClient.ABCIQuery("V", txidBytes)
	compareAssets(t, data.Response.GetValue(), serializedIDDoc, i.Key())
	data2, err := nc.TmClient.ABCIQuery("I", i.Key())
	compareAssets(t, data2.Response.GetValue(), serializedIDDoc, i.Key())

	//Badger Consensus DB - Local Search
	ltx, err := app.Get(txidBytes)
	assert.Nil(t, err, "Error should be nil", err)
	compareAssets(t, ltx, serializedIDDoc, i.Key())

	//Tags Remote
	tresq, terr := nc.TxSearch("tag.qredo_test_tag='abc'", false, 1, 1)
	assert.Nil(t, terr, "Error should be nil", terr)
	assert.True(t, len(tresq.Txs) == 1, "Should have 1 matching tag")
	tresq2, terr2 := nc.TxSearch("tag.qredo_test_tag='notfound'", false, 1, 1)
	assert.Nil(t, terr2, "Error should be nil", terr2)
	assert.True(t, len(tresq2.Txs) == 0, "Should have 1 matching tag")

	//Tags Local
	tlresq, tlerr := core.TxSearch(nil, "tag.qredo_test_tag='abc'", false, 1, 1)
	assert.Nil(t, tlerr, "Error should be nil", tlerr)
	assert.True(t, len(tlresq.Txs) == 1, "Should have 1 matching tag")
	tlresq2, tlerr2 := core.TxSearch(nil, "tag.qredo_test_tag='notfound'", false, 1, 1)
	assert.Nil(t, tlerr2, "Error should be nil", tlerr2)
	assert.True(t, len(tlresq2.Txs) == 0, "Should have 1 matching tag")

	//Local Post
	i3, txid3, serializedIDDoc3, _ := buildTestIDDocLocal(t)
	txidBytes3, _ := hex.DecodeString(txid3)
	retID3, err6 := app.Get(txidBytes3)
	assert.Nil(t, err6, "Error should be nil", err6)
	compareAssets(t, retID3, serializedIDDoc3, i3.Key())
}

func Test_Subscribe(t *testing.T) {

}

func Test_IDDocPostTX(t *testing.T) {
	store := store.NewChainstore()

	i, err := assets.NewIDDoc("1st Ref")
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

	i, err := assets.NewIDDoc("1st Ref")
	i.DataStore = store
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	w, err := assets.NewWallet(i, "BTC")
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
	w2, _ := assets.NewUpdateWallet(w, i)
	//serialize the whole transaction
	serializedTX2, err := w2.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX2 := base64.StdEncoding.EncodeToString(serializedTX2)

	_, err = PostTx(base64EncodedTX2, "127.0.0.1:26657")
	assert.Nil(t, err, "Error should be nil")

}

func Test_NodeConnector(t *testing.T) {
	//qredochain.StartTestChain()

	// dsBackend, err := datastore.NewBoltBackend("datastore.dat")
	// assert.Nil(t, err, "Error should be nil")
	// assert.NotNil(t, dsBackend, "store should not be nil")

	// store, err := datastore.NewStore(datastore.WithBackend(dsBackend), datastore.WithCodec(datastore.NewGOBCodec()))
	// assert.Nil(t, err, "Error should be nil")
	// assert.NotNil(t, store, "store should not be nil")

	logger, err := logger.NewLogger("text", "info")
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, logger, "logger should not be nil")

	nc, err := NewNodeConnector("127.0.0.1:26657", "NODEID", nil, logger)
	assert.NotNil(t, nc, "tmConnector should not be nil")
	assert.Nil(t, err, "Error should be nil")

	//Build IDDoc
	i, err := assets.NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	nc.PostTx(i)
}

func Test_BadgerIDDocPostTX(t *testing.T) {
	store := store.NewChainstore()

	i, err := assets.NewIDDoc("1st Ref")
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
