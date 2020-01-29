package qredochain

//Test the QredoChain app
//Start an App instance on standard port, and wait for it to complete initialization
//Perform test and terminate
//The App posts to itself and accesses data via the Badger Load/Save StoreInterface
//These tests are internal because its the App talking to itself.

import (
	"encoding/base64"
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/stretchr/testify/assert"
)

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
