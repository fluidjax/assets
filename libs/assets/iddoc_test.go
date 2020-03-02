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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IDDoc_Misc(t *testing.T) {
	i, err := NewIDDoc("")
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(nil)
	assert.NotNil(t, err, "Error should not be nil")
	assetError := i.Verify(nil)
	assert.NotNil(t, assetError, "Error should not be nil")
}

func Test_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	assetError := i.Verify(i)
	assert.Nil(t, assetError, "Error should be nil")
}

func Test_Serialize_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")

	data, err := i.SerializeAsset()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, data, "Result should not be nil")

	i.CurrentAsset.Asset = nil
	data, err = i.SerializeAsset()
	assert.NotNil(t, err, "Error should not be nil")
}

func Test_Save_Load(t *testing.T) {
	testName := "ABC!23"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	store := NewMapstore()
	i.DataStore = store
	key := i.Key()
	i.Save()
	retrieved, err := Load(i.DataStore, key)
	assert.Nil(t, err, "Error should be nil")
	print(retrieved)
	iddoc := retrieved.Asset.GetIddoc()
	assert.Equal(t, testName, iddoc.AuthenticationReference, "Load/Save failed")
}
