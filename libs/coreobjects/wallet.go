package coreobjects

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/protobuffer"
)

//WalletPayload - return the wallet Payload object
func (w *Wallet) WalletPayload() *protobuffer.PBWallet {
	signatureAsset := w.PBSignedAsset.Asset
	wallet := signatureAsset.GetWallet()
	return wallet
}

//Verify - Verify a Wallet signature with supplied ID
func (w *Wallet) Verify(i *IDDoc) (bool, error) {
	//Signature
	signature := w.PBSignedAsset.Signature
	if signature == nil {
		return false, errors.New("No Signature")
	}
	if len(signature) == 0 {
		return false, errors.New("Invalid Signature")
	}
	//Message
	data, err := w.SerializePayload()
	if err != nil {
		return false, err
	}
	//Public Key
	payload := i.IDDocPayload()
	blsPK := payload.GetBLSPublicKey()

	rc := crypto.BLSVerify(data, blsPK, signature)
	if rc == 0 {
		return true, nil
	}
	return false, nil
}

//Sign a wallet with the supplied IDDoc - who must be decalred as the wallet owner
func (w *Wallet) Sign(i *IDDoc) (err error) {
	walletOwner := w.Asset.GetOwner()
	signer := i.Key()

	res := bytes.Compare(walletOwner, signer)
	if res != 0 {
		return errors.New("Only the Owner can self sign")
	}

	signature, err := w.SignPayload(i)
	if err != nil {
		return err
	}
	w.PBSignedAsset.Signature = signature
	w.PBSignedAsset.Signers = append(w.PBSignedAsset.Signers, "self")
	return nil
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
