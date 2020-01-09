package coreobjects

import (
	"bytes"
	"crypto/sha256"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/protobuffer"
)

//AuthenticatorInterface Implementations
func (w *Wallet) Serialize() (s []byte, err error) {
	//Use Asset parent method
	return w.SerializePayload()

}

func (w *Wallet) AssetPayload() *protobuffer.PBWallet {
	signatureAsset := w.PBSignedAsset.Asset
	wallet := signatureAsset.GetWallet()
	return wallet
}

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
	data, err := w.Serialize()
	if err != nil {
		return false, err
	}

	//Public Key
	payload := i.AssetPayload()
	blsPK := payload.GetBLSPublicKey()

	rc := crypto.BLSVerify(data, blsPK, signature)

	if rc == 0 {
		return true, nil
	}
	return false, nil

}

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

//Setup a new IDDoc
func NewWallet(iddoc *IDDoc) (w *Wallet, err error) {
	w = EmptyWallet()
	w.store = iddoc.store

	walletKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.PBSignedAsset.Asset.ID = walletKey
	w.PBSignedAsset.Asset.Type = protobuffer.PBAssetType_wallet
	w.PBSignedAsset.Asset.Owner = iddoc.Key()
	return w, nil

	// //Asset
	// asset := &protobuffer.Asset{}
	// asset.Type = protobuffer.AssetType_wallet

	// walletKey, err := RandomBytes(32)
	// if err != nil {
	// 	return nil, errors.New("Fail to generate random key")
	// }
	// asset.ID = walletKey
	// asset.Owner = iddoc.key

	// //Wallet
	// wallet := &protobuffer.Wallet{}

	// //Compose
	// w.Signature.Asset = asset
	// payload := &protobuffer.Asset_Wallet{}
	// payload.Wallet = wallet
	// w.Signature.Asset.Payload = payload
	// w.store = iddoc.store
	// return w, nil

}

func NewUpdateWallet(previousWallet *Wallet, iddoc *IDDoc) (w *Wallet, err error) {
	w = EmptyWallet()
	if previousWallet.store != nil {
		w.store = previousWallet.store
	}
	w.PBSignedAsset.Asset.ID = previousWallet.PBSignedAsset.Asset.ID
	w.PBSignedAsset.Asset.Type = protobuffer.PBAssetType_wallet
	w.PBSignedAsset.Asset.Owner = iddoc.Key() //new owner
	w.previousAsset = &previousWallet.PBSignedAsset
	return w, nil
}

func EmptyWallet() (w *Wallet) {
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

//Rebuild an existing Signed Wallet into WalletDeclaration object
func ReBuildWallet(sig *protobuffer.PBSignedAsset) (w *Wallet, err error) {
	w = &Wallet{}
	w.PBSignedAsset = *sig
	return w, nil
}

//For testing only
func (i *Wallet) SetTestKey() (err error) {
	data, err := i.Serialize()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	i.SetKey(res[:])
	return nil
}
