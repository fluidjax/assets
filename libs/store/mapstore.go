//Mapstore is a fast simply memory based storage of key values
//Use for testing  and in-memory caching

package store

import (
	"crypto/sha256"
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
	keyhash := sha256.Sum256(key)
	val := m.Store[hex.EncodeToString(keyhash[:])]
	return val, nil
}

func (m *Mapstore) Save(key []byte, data []byte) error {
	keyhash := sha256.Sum256(key)
	m.Store[hex.EncodeToString(keyhash[:])] = data
	return nil
}
