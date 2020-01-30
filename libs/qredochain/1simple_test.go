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
	//Bring up a Node
	nc := StartTestConnectionNode(t)
	i, err := assets.NewIDDoc("testdoc")
	i.Sign(i)
	serializedIDDoc, err := i.SerializeSignedAsset()
	assert.Nil(t, err, "Error should be nil", err)

	txid, errorCode, err := nc.PostTx(i)
	fmt.Println(txid)
	assert.True(t, errorCode == CodeTypeOK, "Error should be nil", err)
	time.Sleep(2 * time.Second)

	//Get from the Node using ABCIQuery
	txidBytes, _ := hex.DecodeString(txid)
	data, err := nc.tmClient.ABCIQuery("V", txidBytes)

	//Check its goodA
	compareAssets(t, data.Response.GetValue(), serializedIDDoc, i.Key())

	// msg := &protobuffer.PBSignedAsset{}
	// err = proto.Unmarshal(data.Response.GetValue(), msg)
	// assert.Nil(t, err, "Error should be nil", err)
	// i2, err := assets.ReBuildIDDoc(msg, i.Key())
	// assert.True(t, i.Hash() == i2.Hash(), "Keys dont match")

	//Get from the Node using indirect Asset ID
	data2, err := nc.tmClient.ABCIQuery("I", i.Key())

	compareAssets(t, data2.Response.GetValue(), serializedIDDoc, i.Key())

	// err = proto.Unmarshal(data2.Response.GetValue(), msg)
	// assert.Nil(t, err, "Error should be nil", err)
	// assert.True(t, i.Hash() == i2.Hash(), "Keys dont match")

	//Get the TXID from the internal db using tx_search
	query := fmt.Sprintf("tx.hash='%v'", txid)
	res, err := core.TxSearch(nil, query, true, 1, 30)
	print(res)

	// nc := StartTestConnectionNode(t)
	// i, err := assets.NewIDDoc("testdoc")
	// i.Sign(i)
	// txid, errorCode, err := nc.PostTx(i)
	// fmt.Println(txid)
	// assert.True(t, errorCode == CodeTypeOK, "Error should be nil", err)

	// data, err := nc.GetTx(txid)
	// assert.NotNil(t, data, "Data should not be nil", err)

	// err = app.db.View(func(txn *badger.Txn) error {
	// 	opts := badger.DefaultIteratorOptions
	// 	opts.PrefetchSize = 10
	// 	it := txn.NewIterator(opts)
	// 	defer it.Close()
	// 	for it.Rewind(); it.Valid(); it.Next() {
	// 		print("Inside loop")
	// 		item := it.Item()
	// 		k := item.Key()
	// 		err := item.Value(func(v []byte) error {

	// 			fmt.Printf("key=%s, value=%s\n", hex.EncodeToString(k), hex.EncodeToString(v))
	// 			return nil
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// })
	// assert.Nil(t, err, "Error should be nil", err)

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
