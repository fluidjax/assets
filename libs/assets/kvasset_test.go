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
	"bytes"
	"testing"

	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/store"
	"github.com/stretchr/testify/assert"
)

func Test_KVAsset(t *testing.T) {

	testVal1 := []byte("Field 1 contents")
	testVal2 := []byte("Field 2 contents")
	testUpdatedVal1 := []byte("Field 1 updated")

	//var store *StoreInterface
	store := store.NewMapstore()

	//Generate an Idoc
	testName := "ABC!23"

	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	i.DataStore = store
	i.Save()

	//Generate the Asset & Sign
	a, err := NewKVAsset(i, protobuffer.PBKVAssetType_UnspecifiedKVAsset)

	//add some test fields
	a.SetKV("Field1", testVal1)
	a.SetKV("Field2", testVal2)
	a.SetKV("Field3", testVal1)

	immutable := []string{"Field2", "Field3"}
	a.SetImmutable(immutable)

	//Sign asset
	a.Sign(i)

	//Save to store
	a.Save()
	key := a.Key()

	reconstitutedData, err := Load(store, key)
	ra, err := ReBuildKVAsset(reconstitutedData, key)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, reconstitutedData, "reconstitutedData should not be nil")
	assert.NotNil(t, ra, "rAsset should not be nil")

	q1, _ := a.SerializeAsset()
	q2, _ := ra.SerializeAsset()

	resq := bytes.Compare(q1, q2)
	assert.True(t, resq == 0, "Reconstructed Asset is not identical to the original")

	val1, err := a.GetKV("Field1")
	resq = bytes.Compare(val1, testVal1)
	assert.True(t, resq == 0, "Reconstructed Asset is not identical to the original")

	val2, err := a.GetKV("Field2")
	resq = bytes.Compare(val2, testVal1)
	assert.True(t, resq == 1, "Reconstructed Asset has invalid field")
	resq = bytes.Compare(val2, testVal2)
	assert.True(t, resq == 0, "Reconstructed Asset is not identical to the original")

	val3, err := a.GetKV("Not Exists")
	assert.Nil(t, val3, "Non existence key should be nil")
	assert.Nil(t, err, "Non existence key should not return an error")

	//Test Update
	ua, err := NewUpdateKVAsset(a, i)
	ua.SetKV("Field1", testUpdatedVal1)
	ua.Sign(i)
	ua.Save()

	//Retreve updated
	reconstitutedData, err = Load(store, key)
	rua, err := ReBuildKVAsset(reconstitutedData, key)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, reconstitutedData, "reconstitutedData should not be nil")
	assert.NotNil(t, rua, "Asset should not be nil")

	updatedField1Value, err := rua.GetKV("Field1")
	assert.Nil(t, err, "Error should be nil")
	resq = bytes.Compare(updatedField1Value, testUpdatedVal1)
	assert.True(t, resq == 0, "Reconstructed Updated Asset doesnt have updated field")

	//Check Mandatory
	//Make another update
	ua2, err := NewUpdateKVAsset(ua, i)
	ua2.CurrentAsset.Asset.Payload = ua2.PreviousAsset.Asset.GetPayload()

	err = ua2.SetKV("Field2", testUpdatedVal1)
	assert.NotNil(t, err, "Error should not be nil")

	notUpdateValue, err := ua2.GetKV("Field2")
	resq = bytes.Compare(testVal2, notUpdateValue)
	assert.True(t, resq == 0, "Reconstructed Asset should not have changed")

	err = ua2.SetKV("Field1", testUpdatedVal1)
	assert.Nil(t, err, "Error should be nil")

	//Sav, restore, re-check
	ua2.Save()
	reconstitutedData3, err := Load(store, key)
	assert.Nil(t, err, "Error should be nil")	
	rua3, err := ReBuildKVAsset(reconstitutedData3, key)
	assert.Nil(t, err, "Error should be nil")

	f2, err := rua3.GetKV("Field2")
	assert.Nil(t, err, "Error should be nil")

	resq = bytes.Compare(f2, notUpdateValue)
	assert.True(t, resq == 0, "Reconstructed Asset should not have changed")

}
