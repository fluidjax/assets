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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

func Test_Group(t *testing.T) {

	//var store *StoreInterface
	store := NewMapstore()
	idInitiator, idT1, idT2, idT3 := SetupIDDocs(&store)

	_ = idInitiator

	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	i.Store = &store
	i.Save()

	w, err := NewGroup(i, protobuffer.PBGroupType_TrusteeGroup)
	expression := "t1 + t2 + t3 > 1 "
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w.ConfigureGroup(expression, participants, testDescription)

	b := proto.NewBuffer(nil)
	b.SetDeterministic(true)
	b.Marshal(w.CurrentAsset.Asset)
	msg1, _ := b.DecodeRawBytes(true)
	res1 := sha256.Sum256(msg1)
	fmt.Println(hex.EncodeToString(res1[:]))

	w.Sign(i)

	assert.NotNil(t, w.CurrentAsset.Signature, "Signature is empty")
	res, err := w.Verify(i)

	c := proto.NewBuffer(nil)
	c.SetDeterministic(true)
	c.Marshal(w.CurrentAsset.Asset)
	msg2, _ := c.DecodeRawBytes(true)
	res2 := sha256.Sum256(msg2)
	fmt.Println(hex.EncodeToString(res2[:]))

	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	w.Save()

	retrieved, err := Load(*i.Store, w.Key())
	retrievedGroup := retrieved.Asset.GetGroup()

	assert.Equal(t, testDescription, retrievedGroup.Description, "Load/Save failed")
}

// Tests to ensure objects are determinsitically serialized - something protobuffers by default does not guarantee
// This bug manifests when serializing maps, - their order when serialized is not guaranteed.
// Fix is to use a fork of proto buffers https://github.com/gogo/protobuf
func Test_Determinism(t *testing.T) {
	store := NewMapstore()
	_, idT1, idT2, idT3 := SetupIDDocs(&store)
	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")

	//Non determinsitic serialization is intermittent, 10 times should be enough to show the error
	//However occasionaly it may not occur on the run of every test, increase iterations if debuggging as issue.
	for j := 0; j < 10; j++ {
		w, _ := NewGroup(i, protobuffer.PBGroupType_TrusteeGroup)
		expression := "t1 + t2 + t3 > 1 "
		participants := &map[string][]byte{
			"t1": idT1.Key(),
			"t2": idT2.Key(),
			"t3": idT3.Key(),
		}
		w.ConfigureGroup(expression, participants, testDescription)
		q1, _ := w.SerializeAsset()
		q2, _ := w.SerializeAsset()

		resq := bytes.Compare(q1, q2)
		assert.True(t, resq == 0, "Two consecutives serializations yield different results")
	}
}
