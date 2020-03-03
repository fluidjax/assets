package testsuite

import (
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/stretchr/testify/assert"
)

func Test_IDDoc(t *testing.T) {
	StartTestChain()

	//Standard build
	i := BuildTestIDDoc(t)
	txid, chainErr := nc.PostTx(i)
	assert.Nil(t, chainErr, "Error should be nil")
	assert.NotEqual(t, txid, "", "TxID should not be empty")

	//Verify Good
	assetError := i.Verify(i)
	assert.Nil(t, assetError, "Error should not be nil")

	//Error Verify using Nil as verifier
	assetError = i.Verify(nil)
	assert.NotNil(t, assetError, "Error should not be nil")

	//Error Verify using not-signer as verifier
	ijunk := BuildTestIDDoc(t)
	assetError = i.Verify(ijunk)
	assert.NotNil(t, assetError, "Error should not be nil")

	//Error: Signature to Nil
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Signature = nil
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusSignedAssetFailtoVerify, "Incorrect Error code")

	//Error: AssetID to Nil
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.ID = nil
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeFailToRebuildAsset, "Incorrect Error code")

	//Error: Post Twice - Fail to rebuild
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.ID = nil
	txid, chainErr = nc.PostTx(i)
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeFailToRebuildAsset, "Incorrect Error code")

	//Error: Missing AuthenticationReference
	i = BuildTestIDDoc(t)
	payload, _ := i.Payload()
	payload.AuthenticationReference = ""
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusMissingFields, "Incorrect Error code")

	//Error: Missing BeneficiaryECPublicKey
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.BeneficiaryECPublicKey = nil
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusMissingFields, "Incorrect Error code")

	//Error: Missing SikePublicKey
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.SikePublicKey = nil
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusMissingFields, "Incorrect Error code")

	//Error: Missing BLSPublicKey
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.BLSPublicKey = nil
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusMissingFields, "Incorrect Error code")

	//Error:  Index 2 not 1 on create
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.Index = 2
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusIndexNotZero, "Incorrect Error code")

	//Error:  Index 0 not 1 on create
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.Index = 0
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeConsensusIndexNotZero, "Incorrect Error code")

	//Error:  Update to immutable
	i = BuildTestIDDoc(t)
	txid, chainErr = nc.PostTx(i)
	i.CurrentAsset.Asset.Index = 2
	txid, chainErr = nc.PostTx(i)
	assert.NotNil(t, chainErr, "Error should not be nil")
	assert.True(t, chainErr.Code == assets.CodeCantUpdateImmutableAsset, "Incorrect Error code")
}
