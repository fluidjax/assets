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
	"crypto/rand"
	"crypto/sha256"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/keystore"
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

func containsString(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func SetupIDDocs(store DataSource) (*IDDoc, *IDDoc, *IDDoc, *IDDoc) {
	idP, _ := NewIDDoc("Primary")
	idP.DataStore = store
	idP.Save()

	idT1, _ := NewIDDoc("1")
	idT1.DataStore = store
	idT1.Save()

	idT2, _ := NewIDDoc("2")
	idT2.DataStore = store
	idT2.Save()

	idT3, _ := NewIDDoc("3")
	idT3.DataStore = store
	idT3.Save()

	return idP, idT1, idT2, idT3
}

// Verify - generic verify function
func Verify(msg []byte, signature []byte, iddoc *IDDoc) (bool, error) {
	idDocPayload, err := iddoc.Payload()
	if err != nil {
		return false, err
	}
	blsPK := idDocPayload.GetBLSPublicKey()
	rc := crypto.BLSVerify(msg, blsPK, signature)
	if rc == 0 {
		return true, nil
	}
	return false, nil
}

// Sign - generic Sign Function
func Sign(msg []byte, iddoc *IDDoc) (signature []byte, err error) {
	if iddoc == nil {
		return nil, errors.New("Sign - supplied IDDoc is nil")
	}
	if iddoc.Seed == nil {
		return nil, errors.New("Unable to Sign IDDoc - No Seed")
	}
	_, blsSK, err := keystore.GenerateBLSKeys(iddoc.Seed)
	if err != nil {
		return nil, err
	}
	rc, signature := crypto.BLSSign(msg, blsSK)
	if rc != 0 {
		return nil, errors.New("Failed to Sign Asset")
	}
	return signature, nil
}

func TxHash(rawTX []byte) []byte {
	txHashA := sha256.Sum256(rawTX)
	txHash := txHashA[:]
	return txHash
}
