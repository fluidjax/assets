package store

import (
	"encoding/hex"
)

//Mapstore -
type Mapstore struct {
	Store map[string][]byte
}

//NewMapstore -
func NewMapstore() StoreInterface {
	m := &Mapstore{}
	m.Store = make(map[string][]byte)
	return m
}

func (m *Mapstore) Load(key []byte) ([]byte, error) {
	val := m.Store[hex.EncodeToString(key)]
	return val, nil
}

func (m *Mapstore) Save(key []byte, data []byte) error {
	m.Store[hex.EncodeToString(key)] = data
	return nil

}
