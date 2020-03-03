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
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

// var (
// 	mapstore map[string]proto.Message
// )

func Test_Wallet_Create(t *testing.T) {
	StartTestChain()
	idP, idT1, idT2, idT3 := SetupIDDocs(t)

	//Standard Wallet build
	wallet, err := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
	assert.Nil(t, err, "Truth table return should be nil")
	assert.NotNil(t, wallet, "Wallet should not be nil")

	//Error no Transfer Participants
	txid, chainErr := nc.PostTx(wallet)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, txid == "", "TXID should be empty")
	assert.True(t, chainErr.Code == assets.CodeConsensusWalletNoTransferRules, "Incorrect Error code")

	//Sign with all Participants
	//TODO

	//Add Transfers  - it would now work, but we will break it for testing
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	wallet.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")

	//Error: No Payload
	tempPayload := wallet.CurrentAsset.Asset.Payload
	wallet.CurrentAsset.Asset.Payload = nil
	txid, chainErr = nc.PostTx(wallet)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusErrorEmptyPayload, "Incorrect Error code")
	assert.True(t, txid == "", "TXID should be empty")
	wallet.CurrentAsset.Asset.Payload = tempPayload

	//Error: Non zero balance
	payload, _ := wallet.Payload()
	payload.SpentBalance = -100
	txid, chainErr = nc.PostTx(wallet)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusMissingFields, "Incorrect Error code")
	assert.True(t, txid == "", "TXID should be empty")
	wallet.CurrentAsset.Asset.Payload = tempPayload

	//Error: Non zero balance
	payload, _ = wallet.Payload()
	payload.SpentBalance = 1
	txid, chainErr = nc.PostTx(wallet)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusMissingFields, "Incorrect Error code")
	assert.True(t, txid == "", "TXID should be empty")
	wallet.CurrentAsset.Asset.Payload = tempPayload
	payload.SpentBalance = 0

	//Error: No currency
	payload, _ = wallet.Payload()
	payload.Currency = 0
	txid, chainErr = nc.PostTx(wallet)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusMissingFields, "Incorrect Error code")
	assert.True(t, txid == "", "TXID should be empty")
	wallet.CurrentAsset.Asset.Payload = tempPayload
	payload.Currency = 1

	//No Error
	txid, chainErr = nc.PostTx(wallet)
	assert.Nil(t, chainErr, "Error should not nil")
	assert.True(t, txid != "", "TXID should not be empty")

}

// func Test_TruthTable(t *testing.T) {
// 	store := assets.NewMapstore()
// 	idP, idT1, idT2, idT3 := SetupIDDocs(store)
// 	expression := "t1 + t2 + t3 > 1 & p"
// 	participants := &map[string][]byte{
// 		"p":  idP.Key(),
// 		"t1": idT1.Key(),
// 		"t2": idT2.Key(),
// 		"t3": idT3.Key(),
// 	}
// 	w1, _ := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
// 	w1.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")
// 	//Create another based on previous, ie. AnUpdateWallet
// 	res, err := w1.TruthTable(protobuffer.PBTransferType_SettlePush)
// 	assert.Nil(t, err, "Truth table return should be nil")
// 	displayRes := fmt.Sprintln("[", strings.Join(res, "], ["), "]")
// 	assert.Equal(t, displayRes, "[ 0 + t2 + t3 > 1 & p], [t1 + 0 + t3 > 1 & p], [t1 + t2 + 0 > 1 & p], [t1 + t2 + t3 > 1 & p ]\n", "Truth table invalid")
// }

// func Test_RuleAdd(t *testing.T) {
// 	store := assets.NewMapstore()
// 	idP, idT1, idT2, idT3 := SetupIDDocs(store)
// 	idNewOwner, _ := assets.NewIDDoc("NewOwner")

// 	expression := "t1 + t2 + t3 > 1 & p"
// 	participants := &map[string][]byte{
// 		"p":  idP.Key(),
// 		"t1": idT1.Key(),
// 		"t2": idT2.Key(),
// 		"t3": idT3.Key(),
// 	}

// 	w1, _ := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
// 	w1.DataStore = idP.DataStore
// 	w1.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")

// 	//Create another Wallet based on previous, ie. AnUpdateWallet
// 	w2, _ := assets.NewUpdateWallet(w1, idNewOwner)

// 	//Change Payload to a SettlePush Type Transfer
// 	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_SettlePush

// 	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
// 	sigP, _ := w2.SignAsset(idP)
// 	sigT1, _ := w2.SignAsset(idT1)
// 	sigT2, _ := w2.SignAsset(idT2)
// 	sigT3, _ := w2.SignAsset(idT3)

// 	// //Pass correct
// 	transferSignatures1 := []assets.SignatureID{
// 		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
// 		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
// 		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
// 		assets.SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: nil},
// 	}
// 	validTransfer1, _ := w2.IsValidTransfer(protobuffer.PBTransferType_SettlePush, transferSignatures1)
// 	assert.True(t, validTransfer1, "Transfer should be valid")

// 	//Fail not enough threshold
// 	transferSignatures1 = []assets.SignatureID{
// 		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
// 		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: nil},
// 		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: nil},
// 		assets.SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
// 	}
// 	validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_SettlePush, transferSignatures1)
// 	assert.False(t, validTransfer1, "Transfer should be invalid")

// 	//Fail no principal
// 	transferSignatures1 = []assets.SignatureID{
// 		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: nil},
// 		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
// 		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
// 		assets.SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
// 	}
// 	validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_SettlePush, transferSignatures1)
// 	assert.False(t, validTransfer1, "Transfer should be invalid")

// 	//Pass too many correct
// 	transferSignatures1 = []assets.SignatureID{
// 		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
// 		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
// 		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
// 		assets.SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
// 	}
// 	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_SettlePush, transferSignatures1)
// 	assert.Nil(t, err, "Error should be nil")
// 	assert.True(t, validTransfer1, "Transfer should be valid")
// }

// func Test_AggregationAndVerify(t *testing.T) {
// 	store := assets.NewMapstore()
// 	idP, idT1, idT2, idT3 := SetupIDDocs(store)
// 	idNewOwner, _ := assets.NewIDDoc("NewOwner")
// 	expression := "t1 + t2 + t3 > 1 & p"
// 	participants := &map[string][]byte{
// 		"p":  idP.Key(),
// 		"t1": idT1.Key(),
// 		"t2": idT2.Key(),
// 		"t3": idT3.Key(),
// 	}
// 	w1, _ := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
// 	w1.DataStore = idP.DataStore
// 	w1.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")

// 	//Create another Wallet based on previous, ie. AnUpdateWallet
// 	w2, _ := assets.NewUpdateWallet(w1, idNewOwner)

// 	//Change Payload to a SettlePush Type Transfer
// 	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_SettlePush

// 	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
// 	sigP, _ := w2.SignAsset(idP)
// 	sigT1, _ := w2.SignAsset(idT1)
// 	sigT2, _ := w2.SignAsset(idT2)
// 	sigT3, _ := w2.SignAsset(idT3)

// 	//Everything is sign by the s & Principal

// 	//Add sufficient Signatures
// 	transferSignatures1 := []assets.SignatureID{
// 		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
// 		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
// 		assets.SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
// 		assets.SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
// 	}

// 	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_SettlePush, transferSignatures1)
// 	assert.Nil(t, err, "Error should be nil")
// 	assert.True(t, validTransfer1, "Transfer is invalid boolean doesn't return true")

// 	//Build Aggregated Signature and insert into the releavnt Wallet Fields
// 	err = w2.AggregatedSign(transferSignatures1)
// 	assert.Nil(t, err, "Error should be nil")

// 	//Check wallet2 validatity based on previous Version
// 	verify, err := w2.FullVerify()
// 	assert.True(t, verify, "Verify should be True")
// 	assert.Nil(t, err, "Error should be nil")

// }

// func Test_AggregationAndVerifyFailingTransfer(t *testing.T) {
// 	store := assets.NewMapstore()
// 	idP, idT1, idT2, idT3 := SetupIDDocs(store)
// 	idNewOwner, _ := assets.NewIDDoc("NewOwner")
// 	expression := "t1 + t2 + t3 > 1 & p"
// 	participants := &map[string][]byte{
// 		"p":  idP.Key(),
// 		"t1": idT1.Key(),
// 		"t2": idT2.Key(),
// 		"t3": idT3.Key(),
// 	}

// 	w1, _ := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
// 	w1.DataStore = idP.DataStore
// 	w1.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")

// 	//Create another Wallet based on previous, ie. AnUpdateWallet
// 	w2, _ := assets.NewUpdateWallet(w1, idNewOwner)

// 	//Change Payload to a SettlePush Type Transfer
// 	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_SettlePush

// 	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
// 	sigP, _ := w2.SignAsset(idP)
// 	sigT1, _ := w2.SignAsset(idT1)

// 	//Everything is sign by the s & Principal

// 	//Add sufficient Signatures
// 	transferSignatures1 := []assets.SignatureID{
// 		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
// 		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
// 	}
// 	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_SettlePush, transferSignatures1)
// 	assert.Nil(t, err, "Error should be nil")
// 	assert.False(t, validTransfer1, "Transfer should be invalid")

// 	//Build Aggregated Signature and insert into the releavnt Wallet Fields
// 	err = w2.AggregatedSign(transferSignatures1)
// 	assert.Nil(t, err, "Error should be nil")

// 	//Check wallet2 validatity based on previous Version
// 	verify, err := w2.FullVerify()
// 	assert.False(t, verify, "Verify should be False")
// 	assert.NotNil(t, err, "Error should describe the failure")
// }

// func Test_WalletTransfer(t *testing.T) {
// 	store := assets.NewMapstore()
// 	idP, idT1, idT2, idT3 := SetupIDDocs(store)
// 	idNewOwner, _ := assets.NewIDDoc("NewOwner")
// 	expression := "t1 + t2 + t3 > 1 & p"
// 	participants := &map[string][]byte{
// 		"p":  idP.Key(),
// 		"t1": idT1.Key(),
// 		"t2": idT2.Key(),
// 		"t3": idT3.Key(),
// 	}

// 	w1, _ := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
// 	w1.DataStore = idP.DataStore
// 	w1.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")
// 	wallet, err := w1.Payload()

// 	//Wallet has already spent 100
// 	wallet.SpentBalance = 100

// 	//Create another Wallet based on previous, ie. AnUpdateWallet
// 	w2, _ := assets.NewUpdateWallet(w1, idNewOwner)

// 	//send 30 BTC to idT3
// 	w2.AddWalletTransfer(idT3.Key(), 30, idT3.Key()) //for now we just transfer to the IDDoc Asset, but it should be a Wallet AssetID
// 	w2.AddWalletTransfer(idT2.Key(), 22, idT3.Key())

// 	//Change Payload to a SettlePush Type Transfer
// 	w2.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_SettlePush

// 	//Generate Signatures for each Participant - note they are signing the new Wallet with the TransferType set!
// 	sigP, _ := w2.SignAsset(idP)
// 	sigT1, _ := w2.SignAsset(idT1)

// 	//Everything is sign by the s & Principal

// 	//Add sufficient Signatures
// 	transferSignatures1 := []assets.SignatureID{
// 		assets.SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
// 		assets.SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
// 	}
// 	validTransfer1, err := w2.IsValidTransfer(protobuffer.PBTransferType_SettlePush, transferSignatures1)
// 	assert.Nil(t, err, "Error should be nil")
// 	assert.False(t, validTransfer1, "Transfer should be invalid")

// 	//Build Aggregated Signature and insert into the releavnt Wallet Fields
// 	err = w2.AggregatedSign(transferSignatures1)
// 	assert.Nil(t, err, "Error should be nil")

// 	//Check wallet2 validatity based on previous Version
// 	verify, err := w2.FullVerify()
// 	assert.False(t, verify, "Verify should be False")
// 	assert.NotNil(t, err, "Error should describe the failure")

// 	payload, _ := w2.Payload()
// 	assert.True(t, payload.SpentBalance == 152, "Invalid total spent")

// }

// func Test_ScriptTesting(t *testing.T) {
// 	seed1, _ := hex.DecodeString("af37ab062cae50f77d2a33ff8361671e80451460b1133613c797e0d743897638a7900f17a3e4cb66c2d88cb86bc49a73")
// 	seed2, _ := hex.DecodeString("b442ca9a5e487f6a2fb5bbb43d0205f2ae6f3624c5eff7952005ee592fe5fee433e0b7b59f3065217b26399547b58399")
// 	seed3, _ := hex.DecodeString("a59ed2f6e3666f86fa8593cdc06727a23cbc9f36f6b7e9ca38ece954614f65175be8d027e2bc24dd7b6df71ed93c9d00")
// 	seed4, _ := hex.DecodeString("92c0aad97d8c03cae7ffd4d415e04c3f3115aace732f00e035f9d996498ebf9a6ff8a8053c81683d860aaaea31d86360")

// 	msg, _ := hex.DecodeString("0804122065613663373836653964313465336662323633636666386566303835646666381a20e4cfb76d0c25c746c7bf8957e88b4ce23ab1a1d9147c7e4005d1b84db6dc36c42002280132de010a044e6f6e6512d501080112157435202b207436202b207437203e203120262070321a260a0270321220da8e7b8af72751ac651df4c2681addc451a8435d845721b7c50e9839ef38e3c91a260a0274351220e8eaa3a434e8cd3a26f1b08dc1ae9394791f6e67fef383b91af712121377c4a21a260a0274361220d5f05cd6eb4ac291d8c84495637f95a3eb3ff659dc63ba648b5d498e4b1955101a260a0274371220b642bf9658b458753562adf75263a24ebacd01e785b1d45f36452a15451506ca221a736f6d65206465736372697074696f6e20676f657320686572659a013b080112110a046b65793112096e657776616c75653312110a046b65793212096e657776616c75653212110a046b65793312096e657776616c756531")
// 	iddoc1, _ := assets.NewIDDocWithSeed(seed1, "p1")
// 	iddoc2, _ := assets.NewIDDocWithSeed(seed2, "p1")
// 	iddoc3, _ := assets.NewIDDocWithSeed(seed3, "p1")
// 	iddoc4, _ := assets.NewIDDocWithSeed(seed4, "p1")

// 	signature1, _ := assets.Sign(msg, iddoc1)
// 	fmt.Println(hex.EncodeToString(signature1))
// 	signature2, _ := assets.Sign(msg, iddoc2)
// 	fmt.Println(hex.EncodeToString(signature2))
// 	signature3, _ := assets.Sign(msg, iddoc3)
// 	fmt.Println(hex.EncodeToString(signature3))
// 	signature4, _ := assets.Sign(msg, iddoc4)
// 	fmt.Println(hex.EncodeToString(signature4))

// 	aggregatedSig := signature1

// 	_, aggregatedSig = crypto.BLSAddG1(aggregatedSig, signature2)
// 	_, aggregatedSig = crypto.BLSAddG1(aggregatedSig, signature3)
// 	_, aggregatedSig = crypto.BLSAddG1(aggregatedSig, signature4)

// 	fmt.Println(hex.EncodeToString(aggregatedSig))

// 	aggregatedPublicKey := iddoc1.CurrentAsset.GetAsset().GetIddoc().BLSPublicKey
// 	_, aggregatedPublicKey = crypto.BLSAddG2(aggregatedPublicKey, iddoc2.CurrentAsset.GetAsset().GetIddoc().BLSPublicKey)
// 	_, aggregatedPublicKey = crypto.BLSAddG2(aggregatedPublicKey, iddoc3.CurrentAsset.GetAsset().GetIddoc().BLSPublicKey)
// 	_, aggregatedPublicKey = crypto.BLSAddG2(aggregatedPublicKey, iddoc4.CurrentAsset.GetAsset().GetIddoc().BLSPublicKey)

// 	fmt.Println(hex.EncodeToString(aggregatedPublicKey))

// 	rc := crypto.BLSVerify(msg, aggregatedPublicKey, aggregatedSig)
// 	assert.True(t, rc == 0, "Return should be 0")

// 	m, _ := hex.DecodeString("0804122036323161353261613062393930386335323962636432623336343037303739381a20e4cfb76d0c25c746c7bf8957e88b4ce23ab1a1d9147c7e4005d1b84db6dc36c42003280132de010a044e6f6e6512d501080112157435202b207436202b207437203e203120262070321a260a0270321220da8e7b8af72751ac651df4c2681addc451a8435d845721b7c50e9839ef38e3c91a260a0274351220e8eaa3a434e8cd3a26f1b08dc1ae9394791f6e67fef383b91af712121377c4a21a260a0274361220d5f05cd6eb4ac291d8c84495637f95a3eb3ff659dc63ba648b5d498e4b1955101a260a0274371220b642bf9658b458753562adf75263a24ebacd01e785b1d45f36452a15451506ca221a736f6d65206465736372697074696f6e20676f657320686572659a013b080112110a046b65793112096e657776616c75653312110a046b65793212096e657776616c75653212110a046b65793312096e657776616c756531")
// 	p, _ := hex.DecodeString("0c57ff98f4125a8ef403cd5b974aaf6f8450436186c8a6934f9575c2d4fc4a199ce4e4836803870abb1272c4d2cfc5600640736edb4cf7941c3ffd7ac8bdc96226af193a1f75502cefcfedce7b2ca8cbc2b1aa239352f5b71164a73db4e46e3802b9e88865eea94ace75260ffa8a785eb22434c7fa8c41caaec2fb0e52b26115b4e03bd406faed802ec60eb82094c6980bcf0922f7dc10f31718e9b5b4bc6eb0101f607b06f487af55cf0d77c6c4b210f1a8d74fde16d3aeb3cbe805916e9910")
// 	s, _ := hex.DecodeString("030ea6112061f860fb4ca1b1fac4c09205b339ebebeeb6f548db5e1d4aa4e444e9cc0b8a4b0015962182482fb696bb4708")
// 	rc = crypto.BLSVerify(m, p, s)
// 	//	assert.True(t, rc == 0, "Return should be 0")

// }
