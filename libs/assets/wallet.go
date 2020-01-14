package assets

import (
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

//Payload - return the wallet Payload object
func (w *Wallet) Payload() (*protobuffer.PBWallet, error) {
	if w == nil {
		return nil, errors.New("Wallet is nil")
	}
	if w.currentAsset.Asset == nil {
		return nil, errors.New("Wallet has no asset")
	}
	signatureAsset := w.currentAsset.Asset
	wallet := signatureAsset.GetWallet()
	return wallet, nil
}

//NewWallet - Setup a new IDDoc
func NewWallet(iddoc *IDDoc) (w *Wallet, err error) {
	if iddoc == nil {
		return nil, errors.New("Sign - supplied IDDoc is nil")
	}
	w = emptyWallet()
	w.store = iddoc.store

	walletKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.currentAsset.Asset.ID = walletKey
	w.currentAsset.Asset.Type = protobuffer.PBAssetType_wallet
	w.currentAsset.Asset.Owner = iddoc.Key()
	w.assetKeyFromPayloadHash()
	return w, nil
}

//NewUpdateWallet - Create a NewWallet for updates/transfers based on a previous one
func NewUpdateWallet(previousWallet *Wallet, iddoc *IDDoc) (w *Wallet, err error) {
	w = emptyWallet()
	if previousWallet.store != nil {
		w.store = previousWallet.store
	}
	w.currentAsset.Asset.ID = previousWallet.currentAsset.Asset.ID
	w.currentAsset.Asset.Type = protobuffer.PBAssetType_wallet
	w.currentAsset.Asset.Owner = iddoc.Key() //new owner
	w.previousAsset = previousWallet.currentAsset
	return w, nil
}

//ReBuildWallet an existing Wallet from it's on chain PBSignedAsset
func ReBuildWallet(sig *protobuffer.PBSignedAsset, key []byte) (w *Wallet, err error) {
	if sig == nil {
		return nil, errors.New("ReBuildIDDoc  - sig is nil")
	}
	if key == nil {
		return nil, errors.New("ReBuildIDDoc  - key is nil")
	}
	w = &Wallet{}
	w.currentAsset = sig
	w.setKey(key)
	return w, nil
}

func emptyWallet() (w *Wallet) {
	w = &Wallet{}
	w.currentAsset = &protobuffer.PBSignedAsset{}
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_wallet
	//Wallet
	wallet := &protobuffer.PBWallet{}
	//Compose
	w.currentAsset.Asset = asset
	payload := &protobuffer.PBAsset_Wallet{}
	payload.Wallet = wallet
	w.currentAsset.Asset.Payload = payload
	return w
}