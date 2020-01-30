package qredochain

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/rpc/core"
)

//Test the QredoChain app
//Start an App instance on standard port, and wait for it to complete initialization
//Perform test and terminate
//These are external tests, that is other process such as Nodes querying the QredoChain for values.
//The external process can only access the chain via REST and have no access to badger.

func TestMain(m *testing.M) {
	StartTestChain()
	code := m.Run()
	time.Sleep(2 * time.Second)
	ShutDown()
	os.Exit(code)
}

func Test_IDOC(t *testing.T) {
	//Initialize a Node
	nc := StartTestConnectionNode(t)

	//Data Sources
	//Tendermint DB
	//Badger Consensus DB

	//REMOTE (nc)
	//Create an IDDoc, Post to Network connector
	//add to chain, and wait 2 seconds for the block
	i, txid, serializedIDDoc, err := buildTestIDDoc(t, nc)
	assert.Nil(t, err, "Error should be nil", err)

	//Tendermint DB
	//LOCAL TxSearch - Get the TXID from the internal db using tx_search
	query := fmt.Sprintf("tx.hash='%v'", txid)
	res, err := core.TxSearch(nil, query, true, 1, 30)
	compareAssets(t, res.Txs[0].Tx, serializedIDDoc, i.Key())

	//Remote TXSearch
	rquery := fmt.Sprintf("tx.hash='%s'", txid)
	resq, err := nc.TxSearch(rquery, false, 1, 1)
	compareAssets(t, resq.Txs[0].Tx, serializedIDDoc, i.Key())

	//Badger Consensus DB -
	//REMOTE (NC) - Query ABCI
	//Get from the Node using ABCIQuery
	txidBytes, _ := hex.DecodeString(txid)
	data, err := nc.tmClient.ABCIQuery("V", txidBytes)
	compareAssets(t, data.Response.GetValue(), serializedIDDoc, i.Key())
	data2, err := nc.tmClient.ABCIQuery("I", i.Key())
	compareAssets(t, data2.Response.GetValue(), serializedIDDoc, i.Key())

	//LOCAL from Badger
	ltx, err := app.Get(txidBytes)
	assert.Nil(t, err, "Error should be nil", err)
	compareAssets(t, ltx, serializedIDDoc, i.Key())

}

func buildTestIDDoc(t *testing.T, nc *NodeConnector) (*assets.IDDoc, string, []byte, error) {
	i, err := assets.NewIDDoc("testdoc")
	i.Sign(i)
	serializedIDDoc, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil", err)

	txid, errorCode, err := nc.PostTx(i)
	fmt.Println(txid)
	assert.True(t, errorCode == CodeTypeOK, "Error should be nil", err)
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
