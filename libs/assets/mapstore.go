package assets

import (
	"encoding/hex"
)

type Mapstore struct {
	store map[string][]byte
}

type MapstoreInterface interface {
	load([]byte) []byte
	save([]byte, []byte)
}

func NewMapstore() *Mapstore {
	m := &Mapstore{}
	m.store = make(map[string][]byte)
	return m
}

func (m *Mapstore) load(key []byte) ([]byte, error) {
	val := m.store[hex.EncodeToString(key)]
	return val, nil
}

func (m *Mapstore) save(key []byte, data []byte) error {
	m.store[hex.EncodeToString(key)] = data
	return nil

}
