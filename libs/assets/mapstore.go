/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package assets

import (
	"encoding/hex"
)

//Mapstore -
type Mapstore struct {
	Store map[string][]byte
}

//StoreInterface -
type StoreInterface interface {
	Load([]byte) ([]byte, error)
	Save([]byte, []byte) error
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
