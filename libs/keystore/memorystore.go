// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package keystore

import (
	"sync"
)

// MemoryStore is the in-memory implementation of key store
type MemoryStore struct {
	sync.RWMutex
	keys map[string][]byte
}

// NewMemoryStore creates a new MemoryStore
func NewMemoryStore() (Store, error) {
	return &MemoryStore{
		keys: map[string][]byte{},
	}, nil
}

// Set stores multiple keys at once
func (ms *MemoryStore) Set(name string, key []byte) error {
	ms.Lock()
	defer ms.Unlock()

	ms.keys[name] = make([]byte, len(key))
	copy(ms.keys[name], key)

	return nil
}

// Get retrieves multiple keys
func (ms *MemoryStore) Get(name string) ([]byte, error) {
	ms.RLock()
	defer ms.RUnlock()

	key, ok := ms.keys[name]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return key, nil
}

// SetSeed stores the secret seed
func (ms *MemoryStore) SetSeed(seed []byte) error {
	return ms.Set("seed", seed)
}

// GetSeed returns the stored seed
func (ms *MemoryStore) GetSeed() ([]byte, error) {
	return ms.Get("seed")
}

// GetBLSKeys returns the BLS keys
func (ms *MemoryStore) GetBLSKeys() (blsPublic []byte, blsSecret []byte, err error) {
	keySeed, err := ms.GetSeed()
	if err != nil {
		return nil, nil, err
	}

	return GenerateBLSKeys(keySeed)
}

// GetSIKEKeys returns the SIKE keys
func (ms *MemoryStore) GetSIKEKeys() (sikePublic []byte, sikeSecret []byte, err error) {
	keySeed, err := ms.GetSeed()
	if err != nil {
		return nil, nil, err
	}

	return GenerateSIKEKeys(keySeed)
}
