package qredochain

import (
	"github.com/dgraph-io/badger"
)

//Set - low level set
func (app *QredoChain) RawSet(key, value []byte) error {
	err := app.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, value)
		err := txn.SetEntry(e)
		return err
	})
	return err
}

func (app *QredoChain) BatchGet(key []byte) ([]byte, error) {
	var res []byte
	item, err := app.CurrentBatch.Get(key)
	if item == nil {
		return nil, nil
	}
	err = item.Value(func(val []byte) error {
		res = append([]byte{}, val...) //this copies the item so we can use it outside the closure
		return nil
	})
	return res, err
}

func (app *QredoChain) Set(key []byte, data []byte) (string, error) {
	txn := app.CurrentBatch
	err := txn.Set(key, data)
	return "", err
}
