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

package testsuite

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

func SignWallet(wallet *assets.Wallet, idP, idT1, idT2 *assets.IDDoc) {
	//Sign
	sigP, _ := wallet.SignAsset(idP)
	sigT1, _ := wallet.SignAsset(idT1)
	sigT2, _ := wallet.SignAsset(idT2)

	signatures := []assets.SignatureID{
		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
	}

	wallet.AggregatedSign(signatures)
}

func Test_Wallet_Create(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs(t)

	//Standard Wallet build
	wallet, err := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
	assert.Nil(t, err, "Truth table return should be nil")
	assert.NotNil(t, wallet, "Wallet should not be nil")

	//Add transfers
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}
	wallet.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")
	payload, _ := wallet.Payload()
	SignWallet(wallet, idP, idT1, idT2)

	//This is now a valid wallet
	//Systematically break it to check the errors
	//txid, err := wallet.Save()

	//Non Zero balance - Negative
	payload.SpentBalance = -100
	SignWallet(wallet, idP, idT1, idT2)
	txid, err := wallet.Save()
	assetError, _ := err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code - "+assetError.Code().String())
	assert.True(t, txid == "", "TXID should be empty")
	payload.SpentBalance = 0

	//Non Zero balance - Positive
	payload.SpentBalance = 100
	SignWallet(wallet, idP, idT1, idT2)
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code - "+assetError.Code().String())
	assert.True(t, txid == "", "TXID should be empty")
	payload.SpentBalance = 0

	///No Currency
	payload.Currency = 0
	SignWallet(wallet, idP, idT1, idT2)
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code - "+assetError.Code().String())
	assert.True(t, txid == "", "TXID should be empty")
	payload.Currency = 1

	///No Signature
	tempSig := wallet.CurrentAsset.Signature
	wallet.CurrentAsset.Signature = nil
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusErrorFailtoVerifySignature, "Incorrect Error code - "+assetError.Code().String())
	assert.True(t, txid == "", "TXID should be empty")
	wallet.CurrentAsset.Signature = tempSig

	///No Transfers
	tempTrans := wallet.CurrentAsset.Asset.Transferlist
	wallet.CurrentAsset.Asset.Transferlist = nil
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusWalletNoTransferRules, "Incorrect Error code - "+assetError.Code().String())
	assert.True(t, txid == "", "TXID should be empty")
	wallet.CurrentAsset.Asset.Transferlist = tempTrans

	//VALID!
	SignWallet(wallet, idP, idT1, idT2)
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.Nil(t, assetError, "Error should  be nil")
	assert.False(t, txid == "", "TXID should be empty")

	return

}

func Test_TruthTable(t *testing.T) {
	//Local test of the truth table
	idP, idT1, idT2, idT3 := SetupIDDocs(t)

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}
	w1, _ := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
	w1.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")
	//Create another based on previous, ie. AnUpdateWallet
	res, err := w1.TruthTable(protobuffer.PBTransferType_SettlePush)
	assert.Nil(t, err, "Truth table return should be nil")
	displayRes := fmt.Sprintln("[", strings.Join(res, "], ["), "]")
	assert.Equal(t, displayRes, "[ 0 + t2 + t3 > 1 & p], [t1 + 0 + t3 > 1 & p], [t1 + t2 + 0 > 1 & p], [t1 + t2 + t3 > 1 & p ]\n", "Truth table invalid")
}

func Test_Update(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs(t)

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
	w1.DataStore = idP.DataStore
	w1.AddTransfer(protobuffer.PBTransferType_TransferPush, expression, participants, "description")

	sigP, _ := w1.SignAsset(idP)
	err := w1.AddSigner(idP, "p", sigP)
	assert.Nil(t, err, "Error should be nil")

	_, err = w1.Save()
	assert.Nil(t, err, "Error should be nil")

	//Create another Wallet based on previous, ie. AnUpdateWallet
	w2, _ := assets.NewUpdateWallet(w1, idT1)
	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_TransferPush
	w2.AddTransfer(protobuffer.PBTransferType_TransferPush, expression, participants, "description")

	// //Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
	sigP, _ = w2.SignAsset(idP)
	sigT1, _ := w2.SignAsset(idT1)
	sigT2, _ := w2.SignAsset(idT2)
	sigT3, _ := w2.SignAsset(idT3)

	//Check not enough signatures
	transferSignatures1 := []assets.SignatureID{
		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
	}
	validTransfer1, _ := w2.IsValidTransfer(protobuffer.PBTransferType_TransferPush, transferSignatures1)
	assert.False(t, validTransfer1, "Transfer should be invalid - not enough Signatures")

	//Add signers one at a time
	err = w2.AddSigner(idP, "p", sigP)
	assert.Nil(t, err, "Error should be nil")
	err = w2.AddSigner(idT1, "t1", sigT1)
	assert.Nil(t, err, "Error should be nil")

	txid, err := w2.Save()
	assert.NotNil(t, err, "Error should not be nil")
	assert.True(t, txid == "", "TXID should be blank")
	assetError, _ := err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusTransferRulesFailed, "Incorrect Error code")

	//Check correct more than enough number of sigs
	transferSignatures1 = []assets.SignatureID{
		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		assets.SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_TransferPush, transferSignatures1)
	assert.True(t, validTransfer1, "Transfer should be valid")

	//Check correct number of sigs
	transferSignatures1 = []assets.SignatureID{
		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
	}
	validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_TransferPush, transferSignatures1)
	assert.True(t, validTransfer1, "Transfer should be valid")

	//Sign all in one go
	err = w2.AggregatedSign(transferSignatures1)
	assert.Nil(t, err, "Error should be nil")

	//Send Valid Update
	txid, err = w2.Save()
	assert.Nil(t, err, "Error should  be nil")
	assert.False(t, txid == "", "TXID should have a value")

}

func TestMain(m *testing.M) {
	StartTestChain()
	code := m.Run()
	ShutDown()
	os.Exit(code)
}
