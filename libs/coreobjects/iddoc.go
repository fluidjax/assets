package coreobjects

import (
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
)

type IDDoc struct {
	BaseAsset
}

//AuthenticatorInterface Implementations
func (i *IDDoc) PayloadSerialize() (s []byte, err error) {
	//Use Asset parent method
	return i.BaseAsset.PayloadSerialize()

}

func (i *IDDoc) AssetPayload() *protobuffer.IDDoc {
	return i.Signature.Asset.GetIddoc()
}

func (i *IDDoc) Verify() (bool, error) {

	//Signature
	signature := i.Signature.Signature
	if signature == nil {
		return false, errors.New("No Signature")
	}
	if len(signature) == 0 {
		return false, errors.New("Invalid Signature")
	}

	//Message
	data, err := i.PayloadSerialize()
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

func (i *IDDoc) Sign() (err error) {
	data, err := i.PayloadSerialize()
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

//Create a new IDDoc
func NewIDDoc(authenticationReference string) (i *IDDoc, err error) {
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
	i = &IDDoc{}
	i.seed = seed

	//Asset
	asset := &protobuffer.Asset{}

	//IDDoc
	iddoc := &protobuffer.IDDoc{}
	iddoc.AuthenticationReference = authenticationReference
	iddoc.BeneficiaryECPublicKey = ecPublicKey
	iddoc.SikePublicKey = sikePublicKey
	iddoc.BLSPublicKey = blsPublicKey

	//Compose
	i.Signature.Asset = asset
	assetDefinition := &protobuffer.Asset_Iddoc{}
	assetDefinition.Iddoc = iddoc
	i.Signature.Asset.AssetDefinition = assetDefinition
	i.SetTestKey()
	return i, nil
}

//Rebuild an existing Signed IDDoc into IDDocDeclaration object
//Seed can be manually set if known (ie. Is a local ID)
func ReBuildIDDoc(sig *protobuffer.Signature) (i *IDDoc, err error) {
	i = &IDDoc{}
	i.Signature = *sig
	return i, nil
}
