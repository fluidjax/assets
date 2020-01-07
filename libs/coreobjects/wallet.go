package coreobjects

import (
	"crypto/sha256"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
)

type WalletDeclaration struct {
	Asset
}

//AuthenticatorInterface Implementations
func (w *WalletDeclaration) Serialize() (s []byte, err error) {
	//Use Asset parent method
	return w.Asset.Serialize(w.Asset)

}

func (w *WalletDeclaration) assetPayload() *protobuffer.Wallet {
	signatureAsset := w.Signature.GetDeclaration()
	wallet := signatureAsset.GetWallet()
	return wallet
}

func (w *WalletDeclaration) Verify(i *IDDocDeclaration) (bool, error) {

	//Signature
	signature := w.Signature.Signature
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

func (w *WalletDeclaration) Sign(i *IDDocDeclaration) (err error) {
	data, err := w.Serialize()

	if err != nil {
		return err
	}

	if i.seed == nil {
		return errors.New("No Seed in Supplied IDDoc")
	}
	_, blsSK, err := keystore.GenerateBLSKeys(i.seed)
	if err != nil {
		return err
	}
	rc, signature := crypto.BLSSign(data, blsSK)

	if rc != 0 {
		return errors.New("Failed to sign IDDoc")
	}

	w.Signature.Signature = signature
	w.Signature.Signers = append(w.Signature.Signers, i.key)
	return nil
}

//Setup a new IDDoc
func NewWallet(idkey []byte) (w *WalletDeclaration, err error) {
	//generate crypto random seed

	//Main returned Object
	w = &WalletDeclaration{}

	// build ID Doc AssetDefinition
	assetDeclaration := &protobuffer.AssetDeclaration{}
	assetDefinition := &protobuffer.AssetDeclaration_Wallet{}
	assetDefinition.Wallet = &protobuffer.Wallet{}
	assetDeclaration.AssetDefinition = assetDefinition

	//Assign the signature wrapper
	signature := &protobuffer.Signature_Declaration{}
	signature.Declaration = assetDeclaration

	w.Signature.SignatureAsset = signature
	return w, nil
}

//For testing only
func (i *WalletDeclaration) SetTestKey() (err error) {
	data, err := i.Serialize()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	i.key = res[:]
	return nil
}
