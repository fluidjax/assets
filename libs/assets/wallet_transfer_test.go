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
	"fmt"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

var (
	store map[string]proto.Message
)

func Test_TruthTable(t *testing.T) {
	store := NewMapstore()
	idP, idT1, idT2, idT3 := SetupIDDocs(store)
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}
	w1, _ := NewWallet(idP)
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)
	//Create another based on previous, ie. AnUpdateWallet
	res, err := w1.TruthTable(protobuffer.PBTransferType_settlePush)
	assert.Nil(t, err, "Truth table return should be nil")
	displayRes := fmt.Sprintln("[", strings.Join(res, "], ["), "]")
	assert.Equal(t, displayRes, "[ 0 + t2 + t3 > 1 & p], [t1 + 0 + t3 > 1 & p], [t1 + t2 + 0 > 1 & p], [t1 + t2 + t3 > 1 & p ]\n", "Truth table invalid")
}

func SetupIDDocs(store *Mapstore) (*IDDoc, *IDDoc, *IDDoc, *IDDoc) {
	idP, _ := NewIDDoc("Primary")
	idP.Store = store
	idP.Save()

	idT1, _ := NewIDDoc("1")
	idT1.Store = store
	idT1.Save()

	idT2, _ := NewIDDoc("2")
	idT2.Store = store
	idT2.Save()

	idT3, _ := NewIDDoc("3")
	idT3.Store = store
	idT3.Save()

	return idP, idT1, idT2, idT3
}

func Test_RuleAdd(t *testing.T) {
	store := NewMapstore()
	idP, idT1, idT2, idT3 := SetupIDDocs(store)
	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.Store = idP.Store
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another Wallet based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
	sigP, _ := w2.SignAsset(idP)
	sigT1, _ := w2.SignAsset(idT1)
	sigT2, _ := w2.SignAsset(idT2)
	sigT3, _ := w2.SignAsset(idT3)

	// //Pass correct
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: nil},
	}
	validTransfer1, _ := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.True(t, validTransfer1, "Transfer should be valid")

	//Fail not enough threshold
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: nil},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: nil},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Fail no principal
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: nil},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Pass too many correct
	transferSignatures1 = []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer should be valid")
}

func Test_AggregationAndVerify(t *testing.T) {
	store := NewMapstore()
	idP, idT1, idT2, idT3 := SetupIDDocs(store)
	idNewOwner, _ := NewIDDoc("NewOwner")
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}
	w1, _ := NewWallet(idP)
	w1.Store = idP.Store
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another Wallet based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
	sigP, _ := w2.SignAsset(idP)
	sigT1, _ := w2.SignAsset(idT1)
	sigT2, _ := w2.SignAsset(idT2)
	sigT3, _ := w2.SignAsset(idT3)

	//Everything is sign by the s & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}

	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, validTransfer1, "Transfer is invalid boolean doesn't return true")

	//Build Aggregated Signature and insert into the releavnt Wallet Fields
	err = w2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check wallet2 validatity based on previous Version
	verify, err := w2.FullVerify()
	assert.True(t, verify, "Verify should be True")
	assert.Nil(t, err, "Error should be nil")

}

func Test_AggregationAndVerifyFailingTransfer(t *testing.T) {
	store := NewMapstore()
	idP, idT1, idT2, idT3 := SetupIDDocs(store)
	idNewOwner, _ := NewIDDoc("NewOwner")
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.Store = idP.Store
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)

	//Create another Wallet based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
	sigP, _ := w2.SignAsset(idP)
	sigT1, _ := w2.SignAsset(idT1)

	//Everything is sign by the s & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
	}
	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Build Aggregated Signature and insert into the releavnt Wallet Fields
	err = w2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check wallet2 validatity based on previous Version
	verify, err := w2.FullVerify()
	assert.False(t, verify, "Verify should be False")
	assert.NotNil(t, err, "Error should describe the failure")
}

func Test_WalletTransfer(t *testing.T) {
	store := NewMapstore()
	idP, idT1, idT2, idT3 := SetupIDDocs(store)
	idNewOwner, _ := NewIDDoc("NewOwner")
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.Store = idP.Store
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)
	wallet, err := w1.Payload()

	//Wallet has already spent 100
	wallet.SpentBalance = 100

	//Create another Wallet based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//send 30 BTC to idT3
	w2.AddWalletTransfer(idT3.Key(), 30)
	w2.AddWalletTransfer(idT3.Key(), 22)

	//Change Payload to a SettlePush Type Transfer
	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
	sigP, _ := w2.SignAsset(idP)
	sigT1, _ := w2.SignAsset(idT1)

	//Everything is sign by the s & Principal

	//Add sufficient Signatures
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
	}
	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.Nil(t, err, "Error should be nil")
	assert.False(t, validTransfer1, "Transfer should be invalid")

	//Build Aggregated Signature and insert into the releavnt Wallet Fields
	err = w2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Check wallet2 validatity based on previous Version
	verify, err := w2.FullVerify()
	assert.False(t, verify, "Verify should be False")
	assert.NotNil(t, err, "Error should describe the failure")

	payload, _ := w2.Payload()
	assert.True(t, payload.SpentBalance == 152, "Invalid total spent")

}
