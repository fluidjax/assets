package assets

import (
	"bytes"
	"crypto/rand"
)

//RandomBytes - generate n random bytes
func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

func contains(s [][]byte, e []byte) bool {
	for _, a := range s {
		res := bytes.Compare(a, e)
		if res == 0 {
			return true
		}
	}
	return false
}
