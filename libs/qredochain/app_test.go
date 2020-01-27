package qredochain

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/tendermint"
	"github.com/stretchr/testify/assert"
)

var configFile string
var done chan bool
var ready chan bool
var app *KVStoreApplication

func init() {

	//	flag.StringVar(&configFile, "config", "/home/ubuntu/node/config/config.toml", "Path to config.toml")
	flag.StringVar(&configFile, "config", "/tmp/example/config/config.toml", "Path to config.toml")
}

func StartTestChain() {

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	app = NewKVStoreApplication(db)

	flag.Parse()

	node, err := NewTendermint(app, configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	node.Start()
	defer func() {
		node.Stop()
		node.Wait()
	}()

	done = make(chan bool, 1)
	ready <- true //notify the server is up
	<-done        //wait
}

func Test_ChainPutGet(t *testing.T) {
	go StartTestChain()
	ready = make(chan bool, 1)
	<-ready //wait for server to come up

	i, err := assets.NewIDDoc("1st Ref")
	i.DataStore = app

	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)

	//serialize the whole transaction
	serializedTX1, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil")
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX1)
	_, err = tendermint.PostTx(base64EncodedTX, "127.0.0.1:26657")
	assert.Nil(t, err, "Error should be nil")

	testData := []byte("Some test data")
	testKey := []byte("TestKey")

	//Make a test transaction
	app.currentBatch = app.db.NewTransaction(true)
	err = app.Save(testKey, testData)
	assert.Nil(t, err, "Save returned error")
	app.currentBatch.Commit()

	app.currentBatch = app.db.NewTransaction(false)
	retrievedData, err := app.Load(testKey)
	assert.NotNil(t, retrievedData, "Retrieve data is nil")
	assert.True(t, string(retrievedData) == string(testData), "Failed to retrieve data")
	app.currentBatch.Commit()

}
