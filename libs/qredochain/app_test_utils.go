package qredochain

//Utilties and globals used to Test the App

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/qredo/assets/libs/datastore"
	"github.com/qredo/assets/libs/logger"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/node"
)

var done chan bool
var ready chan bool
var app *KVStoreApplication
var tnode *node.Node

func ShutDown() {
	if tnode != nil {
		tnode.Stop()
		tnode.Wait()
	}
}

func InitiateChain() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	app = NewKVStoreApplication(db)

	flag.Parse()

	tnode, err := NewTendermint(app, "/tmp/example/config/config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	tnode.Start()
	defer func() {
		tnode.Stop()
		tnode.Wait()
	}()

	done = make(chan bool, 1)
	ready <- true //notify the server is up
	<-done        //wait
}

func StartTestChain() {
	go InitiateChain()
	ready = make(chan bool, 1)
	<-ready //wait for server to come up
}

func StartTestConnectionNode(t *testing.T) *NodeConnector {
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
	return nc
}
