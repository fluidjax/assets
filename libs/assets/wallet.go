package assets

import (
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

//WalletPayload - return the wallet Payload object
func (w *Wallet) Payload() *protobuffer.PBWallet {
	signatureAsset := w.PBSignedAsset.Asset
	wallet := signatureAsset.GetWallet()
	return wallet
}

//NewWallet - Setup a new IDDoc
func NewWallet(iddoc *IDDoc) (w *Wallet, err error) {
	w = emptyWallet()
	w.store = iddoc.store

	walletKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.PBSignedAsset.Asset.ID = walletKey
	w.PBSignedAsset.Asset.Type = protobuffer.PBAssetType_wallet
	w.PBSignedAsset.Asset.Owner = iddoc.Key()
	w.assetKeyFromPayloadHash()
	return w, nil
}

//NewUpdateWallet - Create a NewWallet for updates/transfers based on a previous one
func NewUpdateWallet(previousWallet *Wallet, iddoc *IDDoc) (w *Wallet, err error) {
	w = emptyWallet()
	if previousWallet.store != nil {
		w.store = previousWallet.store
	}
	w.PBSignedAsset.Asset.ID = previousWallet.PBSignedAsset.Asset.ID
	w.PBSignedAsset.Asset.Type = protobuffer.PBAssetType_wallet
	w.PBSignedAsset.Asset.Owner = iddoc.Key() //new owner
	w.previousAsset = &previousWallet.PBSignedAsset
	return w, nil
}

//ReBuildWallet an existing Wallet from it's on chain PBSignedAsset
func ReBuildWallet(sig *protobuffer.PBSignedAsset) (w *Wallet, err error) {
	w = &Wallet{}
	w.PBSignedAsset = *sig
	return w, nil
}

func emptyWallet() (w *Wallet) {
	w = &Wallet{}
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_wallet
	//Wallet
	wallet := &protobuffer.PBWallet{}
	//Compose
	w.PBSignedAsset.Asset = asset
	payload := &protobuffer.PBAsset_Wallet{}
	payload.Wallet = wallet
	w.PBSignedAsset.Asset.Payload = payload
	return w
}
