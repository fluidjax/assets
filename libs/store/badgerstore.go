package store

import (
	"encoding/hex"
)

//Mapstore -
type Badgerstore struct {
	Store map[string][]byte
}

//NewMapstore -
func NewBadgerstore() StoreInterface {
	m := &Mapstore{}
	m.Store = make(map[string][]byte)
	return m
}

func (m *Badgerstore) Load(key []byte) ([]byte, error) {
	val := m.Store[hex.EncodeToString(key)]
	return val, nil
}

func (m *Badgerstore) Save(key []byte, data []byte) error {
	m.Store[hex.EncodeToString(key)] = data
	return nil

}
