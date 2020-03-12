package qredochain

import (
	"encoding/binary"

	"github.com/dgraph-io/badger"
)

const (
	hashPath      = "hash"
	heightPath    = "height"
	blockHashPath = "blockhash"
)

//GetLastHeight get the last saved block height
func (app *QredoChain) GetLastHeight() (uint64, error) {
	result, err := app.RawGet([]byte(heightPath))
	if err != nil {
		return 0, err
	}

	var height uint64

	if result != nil {
		height = binary.BigEndian.Uint64(result)
	}
	return height, nil
}

//GetLastHeight get the last saved block height
func (app *QredoChain) GetLastBlockHash() ([]byte, error) {
	var hash [32]byte
	rawHash, err := app.RawGet([]byte(blockHashPath))
	if err != nil {
		return nil, err
	}
	copy(hash[:], rawHash)
	return hash[:], nil
}

//SetLastHeight get the last block height
func (app *QredoChain) SetLastBlockHash(hash []byte) error {
	return app.RawSet([]byte(blockHashPath), hash)
}

//SetLastHeight get the last block height
func (app *QredoChain) SetLastHeight(height uint64) error {
	h := make([]byte, 8)
	binary.BigEndian.PutUint64(h, height)
	return app.RawSet([]byte(heightPath), h)
}

//Get - low level get
func (app *QredoChain) RawGet(key []byte) ([]byte, error) {
	var res []byte
	err := app.DB.View(func(txn *badger.Txn) error {
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
		return nil, err
	}
	return res, err
}

func (app *QredoChain) GetAssetbyID(assetID []byte) ([]byte, error) { //Get Asset using Asset ID
	txid, err := app.RawGet(assetID)
	if err != nil {
		return nil, err
	}
	result, err := app.RawGet(txid)
	return result, nil
}

func (app *QredoChain) CacheSet(key []byte, data []byte) {
	return
}

func (app *QredoChain) CacheGet(key []byte) []byte {
	return nil
}
