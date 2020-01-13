package assets

import (
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
)

func (i *IDDoc) Payload() *protobuffer.PBIDDoc {
	return i.PBSignedAsset.Asset.GetIddoc()
}

//NewIDDoc create a new IDDoc
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
	i.assetKeyFromPayloadHash()
	return i, nil
}

//Rebuild an existing Signed IDDoc into IDDocDeclaration object
//Seed can be manually set if known (ie. Is a local ID)
func ReBuildIDDoc(sig *protobuffer.PBSignedAsset, key []byte) (i *IDDoc, err error) {
	i = &IDDoc{}
	i.PBSignedAsset = *sig
	i.setKey(key)
	return i, nil
}
