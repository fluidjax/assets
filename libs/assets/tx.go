package assets

import (
	"crypto/sha256"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/protobuffer"
)

//TXAsset - generic wrapper for All Transactions
type TXAsset interface {
	ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) uint32
}

type DataSource interface {
	BatchGet(key []byte) ([]byte, error)    //Get data from the current Batch
	BatchSet(key []byte, data []byte) error //Set data in the current Batch
	Get(key []byte) ([]byte, error)         //Get data from underlying commited  database
}

//BuildAssetFromTX -
func BuildAssetFromTX(tx []byte) (txAsset TXAsset, assetID []byte, txHash []byte, err error) {
	signedAsset := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(tx, signedAsset)
	if err != nil {
		return nil, nil, nil, err
	}
	assetID = signedAsset.Asset.GetID()
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
	return txAsset, assetID, txHash, nil
}

func TxHash(rawTX []byte) []byte {
	txHashA := sha256.Sum256(rawTX)
	txHash := txHashA[:]
	return txHash
}
