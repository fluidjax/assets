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
	txid, errorCode, err := nc.PostTx(i)
	fmt.Println(txid)
	assert.True(t, errorCode == CodeTypeOK, "Error should be nil", err)
	time.Sleep(2 * time.Second)

	//Get from the Node using ABCIQuery
	txidBytes, _ := hex.DecodeString(txid)
	data, err := nc.tmClient.ABCIQuery("V", txidBytes)

	//Check its good
	msg := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(data.Response.GetValue(), msg)
	assert.Nil(t, err, "Error should be nil", err)
	i2, err := assets.ReBuildIDDoc(msg, i.Key())
	assert.True(t, i.Hash() == i2.Hash(), "Keys dont match")

	//Get from the Node using indirect Asset ID
	data2, err := nc.tmClient.ABCIQuery("I", i.Key())
	err = proto.Unmarshal(data2.Response.GetValue(), msg)
	assert.Nil(t, err, "Error should be nil", err)
	assert.True(t, i.Hash() == i2.Hash(), "Keys dont match")

}
