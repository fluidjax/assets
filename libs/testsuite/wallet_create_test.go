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

// var (
// 	mapstore map[string]proto.Message
// )

func Test_Wallet_Create(t *testing.T) {
	idP, idT1, idT2, idT3 := SetupIDDocs(t)

	//Standard Wallet build
	wallet, err := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
	assert.Nil(t, err, "Truth table return should be nil")
	assert.NotNil(t, wallet, "Wallet should not be nil")

	//Error no transfers
	txid, err := wallet.Save()
	assetError, _ := err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusWalletNoTransferRules, "Incorrect Error code")

	//Add transfers
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}
	wallet.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusErrorFailtoVerifySignature, "Incorrect Error code")

	//Non Zero balance - Negative
	payload, _ := wallet.Payload()
	payload.SpentBalance = -100
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code")
	payload.SpentBalance = 0

	//Non Zero balance - Positive
	payload.SpentBalance = 100
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code")
	payload.SpentBalance = 0

	//No Currency
	payload.Currency = 0
	txid, err = wallet.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code")
	payload.Currency = 1

	//** No Error** (Not using all Particpiants - but in reality would)
	sigP, _ := wallet.SignAsset(idP)
	sigT1, _ := wallet.SignAsset(idT1)
	sigT2, _ := wallet.SignAsset(idT2)

	signatures := []assets.SignatureID{
		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
	}
	wallet.AggregatedSign(signatures)

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

func TestMain(m *testing.M) {
	StartTestChain()
	code := m.Run()
	ShutDown()
	os.Exit(code)
}
