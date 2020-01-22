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

	"github.com/stretchr/testify/assert"
)

func Test_Wallet_Signing(t *testing.T) {
	testName := "ABC!23"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	i.Store = NewMapstore()
	i.Save()
	w, err := NewWallet(i, "BTC")

	sig1, err := w.SignAsset(i)
	payload, err := w.SerializeAsset()
	sig2, err := Sign(payload, i)
	res := bytes.Compare(sig1, sig2)
	assert.True(t, res == 0, "Compare should be 0")
	verify, err := w.VerifyAsset(sig1, i)
	assert.True(t, verify, "Verify should be true")
	verify2, err := Verify(payload, sig2, i)
	assert.True(t, verify2, "Verify should be true")
}

func Test_Wallet(t *testing.T) {
	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.Sign(i)
	i.Store = NewMapstore()
	i.Save()
	w, err := NewWallet(i, "BTC")
	walletContents, _ := w.Payload()
	walletContents.Description = testDescription
	w.Sign(i)
	assert.NotNil(t, w.CurrentAsset.Signature, "Signature is empty")
	res, err := w.Verify(i)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	w.Save()
	retrieved, err := Load(i.Store, w.Key())
	retrievedWallet := retrieved.Asset.GetWallet()
	assert.Equal(t, testDescription, retrievedWallet.Description, "Load/Save failed")
}

func Test_Wallet_Empty(t *testing.T) {
	w, err := NewWallet(nil, "BTC")
	assert.NotNil(t, err, "Error should not be nil")
	_, err = w.Payload()
	assert.NotNil(t, err, "Error should not be nil")
}
