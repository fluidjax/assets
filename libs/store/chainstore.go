package store

import (
	"encoding/hex"
)

//Mapstore -
type Chainstore struct {
	Store map[string][]byte
}

//NewMapstore -
func NewChainstore() StoreInterface {
	m := &Mapstore{}
	m.Store = make(map[string][]byte)
	return m
}

func (m *Chainstore) Load(key []byte) ([]byte, error) {
	val := m.Store[hex.EncodeToString(key)]
	return val, nil
}

func (m *Chainstore) kSave(key []byte, data []byte) error {
	m.Store[hex.EncodeToString(key)] = data
	return nil

}
