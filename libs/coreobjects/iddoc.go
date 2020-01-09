package coreobjects

import (
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
)

//AuthenticatorInterface Implementations
func (i *IDDoc) PayloadSerialize() (s []byte, err error) {
	//Use Asset parent method
	return i.SerializePayload()

}

func (i *IDDoc) AssetPayload() *protobuffer.PBIDDoc {
	return i.PBSignedAsset.Asset.GetIddoc()
}

func (i *IDDoc) Verify() (bool, error) {

	//Signature
	signature := i.PBSignedAsset.Signature
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
	data, err := i.SerializePayload()
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

	i.PBSignedAsset.Signature = signature
	i.PBSignedAsset.Signers = append(i.PBSignedAsset.Signers, i.Key())
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
	asset := &protobuffer.PBAsset{}

	//IDDoc
	iddoc := &protobuffer.PBIDDoc{}
	iddoc.AuthenticationReference = authenticationReference
	iddoc.BeneficiaryECPublicKey = ecPublicKey
	iddoc.SikePublicKey = sikePublicKey
	iddoc.BLSPublicKey = blsPublicKey

	//Compose
	i.PBSignedAsset.Asset = asset
	Payload := &protobuffer.PBAsset_Iddoc{}
	Payload.Iddoc = iddoc
	i.PBSignedAsset.Asset.Payload = Payload
	i.SetTestKey()
	return i, nil
}

//Rebuild an existing Signed IDDoc into IDDocDeclaration object
//Seed can be manually set if known (ie. Is a local ID)
func ReBuildIDDoc(sig *protobuffer.PBSignedAsset, key []byte) (i *IDDoc, err error) {
	i = &IDDoc{}
	i.PBSignedAsset = *sig
	i.SetKey(key)

	return i, nil
}
