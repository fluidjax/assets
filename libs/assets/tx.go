package assets

import (
	"crypto/sha256"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/protobuffer"
)

//TXAsset - generic wrapper for All Transactions
type TXAsset interface {
	ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) error
}

type DataSource interface {
	BatchGet(key []byte) ([]byte, error)         //Get data from the current Batch
	Set(key []byte, data []byte) (string, error) //Set data in the current Batch
	RawGet(key []byte) ([]byte, error)           //Get data from underlying commited  database
}

//BuildAssetFromTX -
func BuildAssetFromTX(tx []byte) (txAsset TXAsset, assetID []byte, txHash []byte, assetsError error) {
	signedAsset := &protobuffer.PBSignedAsset{}

	//Check 5
	err := proto.Unmarshal(tx, signedAsset)
	if err != nil {
		assetsError := NewAssetsError(CodeFailToRebuildAsset, "Consensus:Error:Check:Invalid Asset Type")
		return nil, nil, nil, assetsError
	}
	assetID = signedAsset.Asset.GetID()
	if assetID == nil {
		print("here")
	}

	txHash = TxHash(tx)

	switch signedAsset.Asset.GetType() {
	case protobuffer.PBAssetType_Wallet:
		txAsset, err = ReBuildWallet(signedAsset, assetID)
	case protobuffer.PBAssetType_Iddoc:
		txAsset, err = ReBuildIDDoc(signedAsset, assetID)
	case protobuffer.PBAssetType_Group:
		txAsset, err = ReBuildGroup(signedAsset, assetID)
	case protobuffer.PBAssetType_Underlying:
		txAsset, err = ReBuildUnderlying(signedAsset, assetID)
	case protobuffer.PBAssetType_MPC:
		txAsset, err = ReBuildMPC(signedAsset, assetID)
	case protobuffer.PBAssetType_KVAsset:
		txAsset, err = ReBuildKVAsset(signedAsset, assetID)
	}

	//Check 5
	if err != nil {
		assetsError = NewAssetsError(CodeFailToRebuildAsset, "Consensus:Error:Check:Invalid Asset Rebuild")
		return nil, nil, nil, assetsError
	}

	return txAsset, assetID, txHash, nil
}

func TxHash(rawTX []byte) []byte {
	txHashA := sha256.Sum256(rawTX)
	txHash := txHashA[:]
	return txHash
}
