package qredochain

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
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

func Test_IDOC(t *testing.T) {
	//Bring up a Node
	nc := StartTestConnectionNode(t)
	i, err := assets.NewIDDoc("testdoc")
	i.Sign(i)
	txid, errorCode, err := nc.PostTx(i)
	fmt.Println(txid)
	assert.True(t, errorCode == CodeTypeOK, "Error should be nil", err)

	err = app.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {

				fmt.Printf("key=%s, value=%s\n", hex.EncodeToString(k), hex.EncodeToString(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	assert.Nil(t, err, "Error should be nil", err)

}
