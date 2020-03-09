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
	txid, err := nc.PostTx(i)
	assert.Nil(t, err, "Error should be nil")
	assert.NotEqual(t, txid, "", "TxID should not be empty")

	//Verify Good
	err = i.Verify(i)
	assert.Nil(t, err, "Error should not be nil")

	//Error Verify using Nil as verifier
	err = i.Verify(nil)
	assert.NotNil(t, err, "Error should not be nil")

	//Error Verify using not-signer as verifier
	ijunk := BuildTestIDDoc(t)
	err = i.Verify(ijunk)
	assert.NotNil(t, err, "Error should not be nil")

	//Error: Signature to Nil
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Signature = nil
	txid, err = nc.PostTx(i)
	assetError, _ := err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusSignedAssetFailtoVerify, "Incorrect Error code")

	//Error: AssetID to Nil
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.ID = nil
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeFailToRebuildAsset, "Incorrect Error code")

	//Error: Post Twice - duplicate TX block by tendermint
	i = BuildTestIDDoc(t)
	txid, err = nc.PostTx(i)
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeTypeTendermintInternalError, "Incorrect Error code")

	//Error: Missing AuthenticationReference
	i = BuildTestIDDoc(t)
	payload, _ := i.Payload()
	payload.AuthenticationReference = ""
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code")

	//Error: Missing BeneficiaryECPublicKey
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.BeneficiaryECPublicKey = nil
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code")

	//Error: Missing SikePublicKey
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.SikePublicKey = nil
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code")

	//Error: Missing BLSPublicKey
	i = BuildTestIDDoc(t)
	payload, _ = i.Payload()
	payload.BLSPublicKey = nil
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusMissingFields, "Incorrect Error code")

	//Error:  Index 2 not 1 on create
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.Index = 2
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusIndexNotZero, "Incorrect Error code")

	//Error:  Index 0 not 1 on create
	i = BuildTestIDDoc(t)
	i.CurrentAsset.Asset.Index = 0
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeConsensusIndexNotZero, "Incorrect Error code")

	//Error:  Update to immutable
	i = BuildTestIDDoc(t)
	txid, err = nc.PostTx(i)
	i.CurrentAsset.Asset.Index = 2
	txid, err = nc.PostTx(i)
	assetError, _ = err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, assetError.Code() == assets.CodeCantUpdateImmutableAsset, "Incorrect Error code")

}
