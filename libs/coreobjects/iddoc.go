package coreobjects

import (
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
)

func (i *IDDoc) IDDocPayload() *protobuffer.PBIDDoc {
	return i.PBSignedAsset.Asset.GetIddoc()
}

//Verify - verify IDDoc with its own BLSPublicKey
func (i *IDDoc) Verify() (bool, error) {

	if len(i.Signers) != 1 {
		return false, errors.New("Signer not specified")
	}
	if i.Signers[0] != "self" {
		return false, errors.New("IDDocs can only be self signed")
	}

	//Signature
	signature := i.PBSignedAsset.Signature
	if signature == nil {
		return false, errors.New("No Signature")
	}
	if len(signature) == 0 {
		return false, errors.New("Invalid Signature")
	}

	//Message
	data, err := i.SerializePayload()
	if err != nil {
		return false, err
	}

	//Public Key
	payload := i.IDDocPayload()
	blsPK := payload.GetBLSPublicKey()

	rc := crypto.BLSVerify(data, blsPK, signature)

	if rc != 0 {
		return false, nil
	}
	return true, nil
}

//Sign an IDDoc with its own BLS Private Key, signer is set to self
func (i *IDDoc) Sign() (err error) {
	data, err := i.SerializePayload()
	if err != nil {
		return err
	}

	if i.seed == nil {
		return errors.New("Unable to Sign IDDoc - No Seed")
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
	i.PBSignedAsset.Signers = append(i.PBSignedAsset.Signers, "self")
	return nil
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
	i.AssetKeyFromPayloadHash()
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
