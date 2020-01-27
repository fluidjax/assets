package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Simple_Mapstore_LoadSave(t *testing.T) {
	store := NewMapstore()
	assert.NotNil(t, store, "Store is nil")
	testData := []byte("Some test data")
	testKey := []byte("TestKey")

	err := store.Save(testKey, testData)
	assert.Nil(t, err, "Save returned error")

	retrievedData, err := store.Load(testKey)
	assert.NotNil(t, retrievedData, "Retrieve data is nil")
	assert.True(t, string(retrievedData) == string(testData), "Failed to retrieve data")
}
