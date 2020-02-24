package qredochain

import (
	"encoding/binary"
	"github.com/dgraph-io/badger"
)


const (
	hashPath           = "hash"
	heightPath         = "height"
)


type AppDB struct {
	*badger.DB
}


func NewAppDB(db *badger.DB) *AppDB {
	d := new(AppDB)
	d.DB = db
	return d
}

func (appDB *AppDB) GetLastHeight() uint64 {
	result := appDB.Get([]byte(heightPath))
	var height uint64

	if result != nil {
		height = binary.BigEndian.Uint64(result)
	}
	return height
}

func (appDB *AppDB) SetLastHeight(height uint64){
	h := make([]byte, 8)
	binary.BigEndian.PutUint64(h, height)
	appDB.Set([]byte(heightPath), h)
}


func (appDB *AppDB) Get(key []byte) ([]byte) {
	var res []byte
	err := appDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if item == nil {
			return nil
		}
		err = item.Value(func(val []byte) error {
			res = append([]byte{}, val...) //this copies the item so we can use it outside the closure
			return nil
		})
		return err
	})
	if err != nil {
		return nil 
	}
	return res
}


func (appDB *AppDB) Set(key ,value []byte) error {
	err := appDB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key,value)
		err := txn.SetEntry(e)
		return err
	  })
	 return err
}