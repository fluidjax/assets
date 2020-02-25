package qredochain

import (
	"encoding/binary"
	"github.com/dgraph-io/badger"
)


const (
	hashPath           = "hash"
	heightPath         = "height"
	blockHashPath         = "blockhash"
)


type AppDB struct {
	*badger.DB
}

//NewAppDB create a new Application Database object
func NewAppDB(db *badger.DB) *AppDB {
	d := new(AppDB)
	d.DB = db
	return d
}


//GetLastHeight get the last saved block height
func (appDB *AppDB) GetLastHeight() uint64 {
	result := appDB.RawGet([]byte(heightPath))
	var height uint64

	if result != nil {
		height = binary.BigEndian.Uint64(result)
	}
	return height
}

//GetLastHeight get the last saved block height
func (appDB *AppDB) GetLastBlockHash() []byte {
	var hash [32]byte
	rawHash := appDB.RawGet([]byte(blockHashPath))
	copy(hash[:], rawHash)
	return hash[:]
}


//SetLastHeight get the last block height
func (appDB *AppDB) SetLastBlockHash(hash []byte){
	appDB.RawSet([]byte(blockHashPath), hash)
}


//SetLastHeight get the last block height
func (appDB *AppDB) SetLastHeight(height uint64){
	h := make([]byte, 8)
	binary.BigEndian.PutUint64(h, height)
	appDB.RawSet([]byte(heightPath), h)
}

//Get - low level get 
func (appDB *AppDB) RawGet(key []byte) ([]byte) {
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

//Set - low level set  
func (appDB *AppDB) RawSet(key ,value []byte) error {
	err := appDB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key,value)
		err := txn.SetEntry(e)
		return err
	  })
	 return err
}