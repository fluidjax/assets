//Mapstore is a fast simply memory based storage of key values
//Use for testing  and in-memory caching

package assets

import (
	"encoding/hex"
)

//Mapstore -
type Mapstore struct {
	Store map[string][]byte
}

//NewMapstore -
func NewMapstore() DataSource {
	m := &Mapstore{}
	m.Store = make(map[string][]byte)
	return m
}

func (m *Mapstore) BatchGet(key []byte) ([]byte, error) {
	val := m.Store[hex.EncodeToString(key)]
	return val, nil
}

func (m *Mapstore) RawGet(key []byte) ([]byte, error) {
	val := m.Store[hex.EncodeToString(key)]
	return val, nil
}

func (m *Mapstore) Set(key []byte, data []byte) (string, error) {
	m.Store[hex.EncodeToString(key)] = data
	return "", nil
}
