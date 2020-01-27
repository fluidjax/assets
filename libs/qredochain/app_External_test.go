package qredochain

import (
	"os"
	"testing"

	"github.com/qredo/assets/libs/assets"
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
	ShutDown()
	os.Exit(code)
}

func Test_External_Query(t *testing.T) {
	//Bring up a Node
	nc := StartTestConnectionNode(t)

	i, err := assets.NewIDDoc("testdoc")
	assert.Nil(t, err, "Error should be nil", err)

	txid, errorCode, err := nc.PostTx(i)
	assert.True(t, errorCode == CodeFailVerfication, "Error should not be nil", err)

	i.Sign(i)
	txid, errorCode, err = nc.PostTx(i)
	assert.True(t, errorCode == CodeTypeOK, "Error should be nil", err)

	assert.Nil(t, err, "Error should be nil", err)
	assert.NotNil(t, txid, "TXID shouldnt be nil", err)

}
