package coreobjects

import (
	"bytes"
	"crypto/sha256"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
)

type IDDocDeclaration struct {
	Asset
}

//AuthenticatorInterface Implementations
func (i *IDDocDeclaration) Serialize() (s []byte, err error) {
	//Use Asset parent method
	return i.Asset.Serialize(i.Asset)

}

func (i *IDDocDeclaration) AssetPayload() *protobuffer.IDDoc {
	signatureAsset := i.Signature.GetDeclaration()
	iddoc := signatureAsset.GetIddoc()
	return iddoc
}

func (i *IDDocDeclaration) Verify() (bool, error) {

	//Signature
	signature := i.Signature.Signature
	if signature == nil {
		return false, errors.New("No Signature")
	}
	if len(signature) == 0 {
		return false, errors.New("Invalid Signature")
	}

	//Message
	data, err := i.Serialize()
	if err != nil {
		return false, err
	}

	//Public Key
	payload := i.AssetPayload()
	blsPK := payload.GetBLSPublicKey()

	rc := crypto.BLSVerify(data, blsPK, signature)

	if rc != 0 {
		return false, nil
	}
	return true, nil

}

func (i *IDDocDeclaration) Sign(iddoc *IDDocDeclaration) error {
	res := bytes.Compare(iddoc.key, i.key)
	if res != 0 {
		return errors.New("Only Self Sign a IDDoc")
	}
	return i.SelfSign()
}

func (i *IDDocDeclaration) SelfSign() (err error) {
	data, err := i.Serialize()
	if err != nil {
		return err
	}

	if i.seed == nil {
		return errors.New("No Seed")
	}
	_, blsSK, err := keystore.GenerateBLSKeys(i.seed)
	if err != nil {
		return err
	}
	rc, signature := crypto.BLSSign(data, blsSK)

	if rc != 0 {
		return errors.New("Failed to sign IDDoc")
	}

	i.Signature.Signature = signature
	i.Signature.Signers = append(i.Signature.Signers, i.key)
	return nil
}

//Setup a new IDDoc
func NewIDDoc(authenticationReference string) (i *IDDocDeclaration, err error) {
	//generate crypto random seed
	seed, err := cryptowallet.RandomBytes(48)
	if err != nil {
		err = errors.Wrap(err, "Failed to generate random seed")
		return nil, err
	}

	sikePublicKey, _, err := keystore.GenerateSIKEKeys(seed)
	if err != nil {
		return nil, err
	}

	blsPublicKey, _, err := keystore.GenerateBLSKeys(seed)
	if err != nil {
		return nil, err
	}

	ecPublicKey, err := keystore.GenerateECPublicKey(seed)
	if err != nil {
		return nil, err
	}

	//Main returned Object
	i = &IDDocDeclaration{}
	i.seed = seed

	// build ID Doc AssetDefinition
	idDocument := &protobuffer.AssetDeclaration{}

	assetDefinition := &protobuffer.AssetDeclaration_Iddoc{}
	assetDefinition.Iddoc = &protobuffer.IDDoc{}
	assetDefinition.Iddoc.AuthenticationReference = authenticationReference
	assetDefinition.Iddoc.BeneficiaryECPublicKey = ecPublicKey
	assetDefinition.Iddoc.SikePublicKey = sikePublicKey
	assetDefinition.Iddoc.BLSPublicKey = blsPublicKey
	idDocument.AssetDefinition = assetDefinition

	//Assign the signature wrapper
	signature := &protobuffer.Signature_Declaration{}
	signature.Declaration = idDocument

	i.Signature.SignatureAsset = signature

	return i, nil
}

//For testing only
func (i *IDDocDeclaration) SetTestKey() (err error) {
	data, err := i.Serialize()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	i.key = res[:]
	return nil
}
