package qredochain

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/logger"
	"github.com/qredo/assets/libs/protobuffer"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
)

var done chan bool
var ready chan bool
var app *QredoChain
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
	app = NewQredoChain(db)

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
	return nc
}

func buildTestIDDocLocal(t *testing.T) (*assets.IDDoc, string, []byte, error) {
	i, err := assets.NewIDDoc("testdoc2")
	i.AddTag("qredo_test_tag2", []byte("abc2"))
	i.Sign(i)
	serializedIDDoc, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil", err)

	result, err := core.BroadcastTxCommit(&rpctypes.Context{}, serializedIDDoc)
	assert.Nil(t, err, "Error should be nil", err)
	time.Sleep(2 * time.Second)
	return i, result.Hash.String(), serializedIDDoc, nil
}

func buildTestIDDoc(t *testing.T, nc *NodeConnector) (*assets.IDDoc, string, []byte, error) {
	i, err := assets.NewIDDoc("testdoc")
	i.AddTag("qredo_test_tag", []byte("abc"))
	i.Sign(i)
	serializedIDDoc, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil", err)

	txid, chainerr := nc.PostTx(i)
	fmt.Println(txid)
	assert.Nil(t, chainerr, "Error should be nil")
	time.Sleep(2 * time.Second)
	return i, txid, serializedIDDoc, nil
}

func compareAssets(t *testing.T, a1, a2, key []byte) {
	msg1 := &protobuffer.PBSignedAsset{}
	err1 := proto.Unmarshal(a1, msg1)
	assert.Nil(t, err1, "Error should be nil", err1)

	msg2 := &protobuffer.PBSignedAsset{}
	err2 := proto.Unmarshal(a2, msg2)
	assert.Nil(t, err2, "Error should be nil", err2)

	i1, err3 := assets.ReBuildIDDoc(msg1, key)
	assert.Nil(t, err3, "Error should be nil", err3)

	i2, err4 := assets.ReBuildIDDoc(msg2, key)
	assert.Nil(t, err4, "Error should be nil", err4)

	assert.True(t, i1.Hash() == i2.Hash(), "Keys dont match")
}
