package testsuite

import (
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

func checkIDDoc(t *testing.T, i *assets.IDDoc, code assets.TransactionCode) {
	txid, err := i.Save()
	assetError, _ := err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, txid == "", "TX Should be empty")
	assert.True(t, assetError.Code() == code, "Incorrect Error code - "+assetError.Code().String())

}

func Test_IDDoc(t *testing.T) {

	//Standard build
	i := BuildTestIDDoc(t)
	txid, err := i.Save()

	assert.Nil(t, err, "Error should be nil")
	assert.NotEqual(t, txid, "", "TxID should not be empty")

	//Verify Good
	err = i.Verify()
	assert.Nil(t, err, "Error should not be nil")

	//SignedAsset.VerifyImmutableCreate()
	//VerifyAll()
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.ID = nil
	checkIDDoc(t, i, assets.CodeFailToRebuildAsset)

	//Nil Payload
	i = BuildTestIDDoc(t)
	i.CurrentAsset.GetAsset().Payload = nil
	checkIDDoc(t, i, assets.CodePayloadEncodingError)

	//Error: Signature to Nil (manually, else its signed by checkIDDoc)
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Signature = nil
	txid, err = i.Save()
	assetError, _ := err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusErrorFailtoVerifySignature, "Incorrect Error code - "+assetError.Code().String())

	//Index != 1
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.Index = 2
	checkIDDoc(t, i, assets.CodeConsensusIndexNotOne)

	//Transfer has values - should be nil
	i = BuildTestIDDoc(t)
	participants := &map[string][]byte{
		"p": i.Key(),
	}
	i.AddTransfer(protobuffer.PBTransferType_SettlePush, "TEST", participants, "description")
	checkIDDoc(t, i, assets.CodeConsensusWalletHasTransferRules)

	//BeneficiaryECPublicKey is empty
	i = BuildTestIDDoc(t)
	payload, _ := i.Payload()
	payload.BeneficiaryECPublicKey = nil
	checkIDDoc(t, i, assets.CodeConsensusMissingFields)

	//SikePublicKey is empty
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.SikePublicKey = nil
	checkIDDoc(t, i, assets.CodeConsensusMissingFields)

	//SikePublicKey is empty
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.SikePublicKey = nil
	checkIDDoc(t, i, assets.CodeConsensusMissingFields)

	//Error: Missing AuthenticationReference
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.AuthenticationReference = ""
	checkIDDoc(t, i, assets.CodeConsensusMissingFields)

	//Error: AssetID to Nil
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.ID = nil
	checkIDDoc(t, i, assets.CodeFailToRebuildAsset)

	//Signature fails to BLSVBerify
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Signature = []byte("0102030405060708")
	txid, err = i.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusSignedAssetFailtoVerify, "Incorrect Error code - "+assetError.Code().String())

	//Error: Post Twice - duplicate TX block by tendermint
	i = BuildTestIDDoc(t)
	txid, err = i.Save()
	assert.Nil(t, err, "Error should  be nil")
	assert.True(t, txid != "", "TX Should NOT be empty")
	txid, err = i.Save()
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeTypeTendermintInternalError, "Incorrect Error code")

}
