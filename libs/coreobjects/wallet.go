package coreobjects

import (
	"crypto/sha256"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/protobuffer"
)

type Wallet struct {
	BaseAsset
}

//AuthenticatorInterface Implementations
func (w *Wallet) Serialize() (s []byte, err error) {
	//Use Asset parent method
	return w.BaseAsset.PayloadSerialize()

}

func (w *Wallet) AssetPayload() *protobuffer.Wallet {
	signatureAsset := w.SignedAsset.Asset
	wallet := signatureAsset.GetWallet()
	return wallet
}

func (w *Wallet) Verify(i *IDDoc) (bool, error) {

	//Signature
	signature := w.SignedAsset.Signature
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

	signature, err := w.BaseAsset.SignPayload(i)
	if err != nil {
		return err
	}
	// data, err := w.Serialize()

	// if err != nil {
	// 	return err
	// }

	// if i.seed == nil {
	// 	return errors.New("No Seed in Supplied IDDoc")
	// }
	// _, blsSK, err := keystore.GenerateBLSKeys(i.seed)
	// if err != nil {
	// 	return err
	// }
	// rc, signature := crypto.BLSSign(data, blsSK)

	// if rc != 0 {
	// 	return errors.New("Failed to sign IDDoc")
	// }

	w.SignedAsset.Signature = signature
	w.SignedAsset.Signers = append(w.SignedAsset.Signers, i.key)
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
	w.SignedAsset.Asset.AsssetID = walletKey
	w.SignedAsset.Asset.Type = protobuffer.AssetType_wallet
	w.SignedAsset.Asset.Owner = iddoc.key
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
	// assetDefinition := &protobuffer.Asset_Wallet{}
	// assetDefinition.Wallet = wallet
	// w.Signature.Asset.AssetDefinition = assetDefinition
	// w.store = iddoc.store
	// return w, nil

}

func NewUpdateWallet(previousWallet *Wallet, iddoc *IDDoc) (w *Wallet, err error) {
	w = EmptyWallet()
	w.SignedAsset.Asset.AsssetID = previousWallet.SignedAsset.Asset.AsssetID
	w.SignedAsset.Asset.Type = protobuffer.AssetType_wallet
	w.SignedAsset.Asset.Owner = iddoc.key //new owner
	w.previousAsset = &previousWallet.BaseAsset
	return w, nil
}

func EmptyWallet() (w *Wallet) {
	w = &Wallet{}
	//Asset
	asset := &protobuffer.Asset{}
	asset.Type = protobuffer.AssetType_wallet
	//Wallet
	wallet := &protobuffer.Wallet{}
	//Compose
	w.SignedAsset.Asset = asset
	assetDefinition := &protobuffer.Asset_Wallet{}
	assetDefinition.Wallet = wallet
	w.SignedAsset.Asset.AssetDefinition = assetDefinition
	return w
}

//Rebuild an existing Signed Wallet into WalletDeclaration object
func ReBuildWallet(sig *protobuffer.SignedAsset) (w *Wallet, err error) {
	w = &Wallet{}
	w.SignedAsset = *sig
	return w, nil
}

//For testing only
func (i *Wallet) SetTestKey() (err error) {
	data, err := i.Serialize()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	i.key = res[:]
	return nil
}
