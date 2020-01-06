package main

import (
	"time"

	"github.com/pkg/errors"
	"github.com/qredo/assets/cryptowallet"
	"github.com/qredo/assets/keystore"
	"github.com/qredo/assets/protobuffer"
)

func main() {
	CreateID()
	//Create an IDDoc for Principal & 3 trustees
	//Create an Wallet

	// //Get Current State of Asset
	// // [Declare/Update] Asset

	// //Create Identity Document
	// 	type = 2
	// 	ID = new ID()
	// 	owner = self ID
	// 	assetDefinition =
	// 		string AuthenticationReference = 1;
	// 		bytes BeneficiaryECPublicKey   = 2;
	// 		bytes SikePublicKey            = 3;
	// 		bytes BLSPublicKey             = 4;
	// 		int64 Timestamp                = 5;
	// 	transfers = nil
	// 	signature = sig of entire message

	// //Create Wallet
	// 	type = 2
	// 	ID = new ID()
	// 	owner = self ID
	// 	assetDefinition =
	// 		wallet =

	// 	transfers = nil
	// 	signature = sig of entire message

	// //Create Trustee Group

}

//CreateID - Create a new DocID along with keys
func CreateID() (IDDocDeclaration *protobuffer.AssetDeclaration, IDDoc []byte, seed []byte, err error) {

	//generate crypto random seed
	seed, err = cryptowallet.RandomBytes(48)
	if err != nil {
		err = errors.Wrap(err, "Failed to generate random seed")
		return nil, nil, nil, err
	}

	sikePublicKey, _, err := keystore.GenerateSIKEKeys(seed)
	if err != nil {
		return
	}

	blsPublicKey, blsSecretKey, err := keystore.GenerateBLSKeys(seed)
	if err != nil {
		return
	}

	ecPublicKey, err := keystore.GenerateECPublicKey(seed)
	if err != nil {
		return
	}

	// build ID Doc
	//idDocument := protobuffer.AssetDeclaration{}
	// idDocument.AuthenticationReference = name
	// idDocument.BeneficiaryECPublicKey = ecPublicKey
	// idDocument.SikePublicKey = sikePublicKey
	// idDocument.BLSPublicKey = blsPublicKey
	// idDocument.Timestamp = time.Now().Unix()

	// // encode ID Doc
	// rawIDDoc, err = documents.EncodeIDDocument(idDocument, blsSecretKey)
	// if err != nil {
	// 	err = errors.Wrap(err, "Failed to encode IDDocument")
	// 	return
	// }

	return nil, nil, nil, nil
}

func DeclareWallet() (*protobuffer.AssetDeclaration, error) {
	return nil, nil
}

func updateWallet() {

}

func declareTrusteeGroup() {

}

func updateTrusteeGroup() {

}
