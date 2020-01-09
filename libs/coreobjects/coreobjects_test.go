package coreobjects

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

func Test_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")
	i.SetTestKey()
	i.Sign()
	i.Description()
	res, err := i.Verify()
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	i.Dump()
}

func Test_Serialize_IDDoc(t *testing.T) {
	i, err := NewIDDoc("chris")
	assert.Nil(t, err, "Error should be nil")

	data, err := i.PayloadSerialize()
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, data, "Result should not be nil")

	i.PBSignedAsset.Asset = nil
	data, err = i.PayloadSerialize()
	assert.NotNil(t, err, "Error should not be nil")
}

func Test_Save_Load(t *testing.T) {
	testName := "ABC!23"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.SetTestKey()
	i.Sign()
	i.store = NewMapstore()
	key := i.Key()
	i.Save()
	retrieved, err := Load(i.store, key)
	assert.Nil(t, err, "Error should be nil")
	print(retrieved)
	iddoc := retrieved.Asset.GetIddoc()
	assert.Equal(t, testName, iddoc.AuthenticationReference, "Load/Save failed")
}

func Test_Wallet(t *testing.T) {
	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")
	i.SetTestKey()
	i.Sign()
	i.store = NewMapstore()
	i.Save()

	w, err := NewWallet(i)
	walletContents := w.AssetPayload()
	walletContents.Description = testDescription
	w.SetTestKey()
	w.Sign(i)
	assert.NotNil(t, w.PBSignedAsset.Signature, "Signature is empty")
	res, err := w.Verify(i)
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res, "Verify should be true")
	w.Save()

	retrieved, err := Load(i.store, w.Key())
	retrievedWallet := retrieved.Asset.GetWallet()

	assert.Equal(t, testDescription, retrievedWallet.Description, "Load/Save failed")
}

func Test_TruthTable(t *testing.T) {

	store := NewMapstore()

	idP, _ := NewIDDoc("Primary")
	idP.store = store
	idP.Save()

	idT1, _ := NewIDDoc("trustee1")
	idT1.store = store
	idT1.Save()

	idT2, _ := NewIDDoc("trustee2")
	idT2.store = store
	idT2.Save()

	idT3, _ := NewIDDoc("trustee3")
	idT3.store = store
	idT3.Save()

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.SetTestKey()
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)
	//w1.Dump()

	//Create another based on previous, ie. AnUpdateWallet

	res, err := w1.TruthTable(protobuffer.PBTransferType_settlePush)
	assert.Nil(t, err, "Truth table return should be nil")

	displayRes := fmt.Sprintln("[", strings.Join(res, "], ["), "]")
	assert.Equal(t, displayRes, "[ 0 + t2 + t3 > 1 & p], [t1 + 0 + t3 > 1 & p], [t1 + t2 + 0 > 1 & p], [t1 + t2 + t3 > 1 & p ]\n", "Truth table invalid")

}

func Test_RuleAdd(t *testing.T) {
	store := NewMapstore()

	idP, _ := NewIDDoc("Primary")
	idP.store = store
	idP.Save()

	idT1, _ := NewIDDoc("trustee1")
	idT1.store = store
	idT1.Save()

	idT2, _ := NewIDDoc("trustee2")
	idT2.store = store
	idT2.Save()

	idT3, _ := NewIDDoc("trustee3")
	idT3.store = store
	idT3.Save()

	idNewOwner, _ := NewIDDoc("NewOwner")

	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w1, _ := NewWallet(idP)
	w1.store = store
	w1.SetTestKey()
	w1.AddTransfer(protobuffer.PBTransferType_settlePush, expression, participants)
	//w1.Dump()

	//Create another based on previous, ie. AnUpdateWallet
	w2, _ := NewUpdateWallet(w1, idNewOwner)

	//Change Payload to a SettlePush Type Transfer
	w2.Asset.TransferType = protobuffer.PBTransferType_settlePush

	//Generate Signatures for each Participant
	sigP, _ := w2.SignPayload(idP)
	sigT1, _ := w2.SignPayload(idT1)
	sigT2, _ := w2.SignPayload(idT2)
	sigT3, _ := w2.SignPayload(idT3)

	// res0, _ := w2.VerifyPayload(sigP, idP)
	// assert.True(t, res0, "Sig fails to verify")
	// res1, _ := w2.VerifyPayload(sigT1, idT1)
	// assert.True(t, res1, "Sig fails to verify")
	// res2, _ := w2.VerifyPayload(sigT2, idT2)
	// assert.True(t, res2, "Sig fails to verify")
	// res3, _ := w2.VerifyPayload(sigT3, idT3)
	// assert.True(t, res3, "Sig fails to verify")

	// rc := 0

	// aggsig := sigP
	// rc, aggsig = crypto.BLSAddG1(aggsig, sigT1)
	// rc, aggsig = crypto.BLSAddG1(aggsig, sigT2)
	// rc, aggsig = crypto.BLSAddG1(aggsig, sigT3)

	// addpk := idP.GetAsset().GetIddoc().GetBLSPublicKey()
	// rc, addpk = crypto.BLSAddG2(addpk, idT1.GetAsset().GetIddoc().GetBLSPublicKey())
	// rc, addpk = crypto.BLSAddG2(addpk, idT2.GetAsset().GetIddoc().GetBLSPublicKey())
	// rc, addpk = crypto.BLSAddG2(addpk, idT3.GetAsset().GetIddoc().GetBLSPublicKey())

	// rc = crypto.BLSVerify(data, idP.GetAsset().GetIddoc().GetBLSPublicKey(), sigP)
	// assert.Equal(t, rc, 0, "Should be zero")
	// rc = crypto.BLSVerify(data, idT1.GetAsset().GetIddoc().GetBLSPublicKey(), sigT1)
	// assert.Equal(t, rc, 0, "Should be zero")
	// rc = crypto.BLSVerify(data, idT2.GetAsset().GetIddoc().GetBLSPublicKey(), sigT2)
	// assert.Equal(t, rc, 0, "Should be zero")
	// rc = crypto.BLSVerify(data, idT3.GetAsset().GetIddoc().GetBLSPublicKey(), sigT3)
	// assert.Equal(t, rc, 0, "Should be zero")

	// rc = crypto.BLSVerify(data, addpk, aggsig)
	// fmt.Println(hex.EncodeToString(data))
	// //fmt.Println(hex.EncodeToString(addpk))
	// //fmt.Println(hex.EncodeToString(aggsig))

	// //assert.Equal(t, rc, 0, "Should be zero")

	//print("done")
	// //We should have 3 valid sigs.

	// //Pass correct
	// transferSignatures1 := []SignatureID{
	// 	SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
	// 	SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
	// 	SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
	// 	SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: nil},
	// }
	// validTransfer1, _ := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	// assert.True(t, validTransfer1, "Transfer should be valid")

	// //Fail not enough threshold
	// transferSignatures1 = []SignatureID{
	// 	SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
	// 	SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: nil},
	// 	SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: nil},
	// 	SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	// }
	// validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	// assert.False(t, validTransfer1, "Transfer should be invalid")

	// //Fail no principal
	// transferSignatures1 = []SignatureID{
	// 	SignatureID{IDDoc: idP, Abbreviation: "p", Signature: nil},
	// 	SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
	// 	SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
	// 	SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	// }
	// validTransfer1, _ = w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	// assert.False(t, validTransfer1, "Transfer should be invalid")

	//Pass too many correct
	transferSignatures1 := []SignatureID{
		SignatureID{IDDoc: idP, Abbreviation: "p", Signature: sigP},
		SignatureID{IDDoc: idT1, Abbreviation: "t1", Signature: sigT1},
		SignatureID{IDDoc: idT2, Abbreviation: "t2", Signature: sigT2},
		SignatureID{IDDoc: idT3, Abbreviation: "t3", Signature: sigT3},
	}
	validTransfer1, _ := w2.IsValidTransfer(protobuffer.PBTransferType_settlePush, transferSignatures1)
	assert.True(t, validTransfer1, "Transfer should be valid")

	// data, _ = w2.SerializePayload()
	// fmt.Println(hex.EncodeToString(data))

	//Build Aggregated Signature and insert into the releavnt Wallet Fields
	w2.AggregatedSign(transferSignatures1)

	// data, _ = w2.SerializePayload()
	// fmt.Println(hex.EncodeToString(data))

	//Check wallet2 validatity based on previous Version
	verify, err := w2.FullVerify(w2.previousAsset)

	assert.True(t, verify, "Verify should be True")
	assert.Nil(t, err, "Error should be nil")

}
