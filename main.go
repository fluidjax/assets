package main

import (
	"github.com/gogo/protobuf/proto"
)

var store = make(map[string]proto.Message)

// func main() {

// 	//Create an IDDoc
// 	iddoc1, _, _ := CreateID("Chris")
// 	iddoc1Key := Store(iddoc1)

// 	//Create an WalletDeclaration
// 	wallet1, _, _ := CreateWallet(iddoc1Key)

// 	DumpStore()
// 	//Create an IDDoc for Principal & 3 trustees
// 	//Create an Wallet

// 	// //Get Current State of Asset
// 	// // [Declare/Update] Asset

// 	// //Create Identity Document
// 	// 	type = 2
// 	// 	ID = new ID()
// 	// 	owner = self ID
// 	// 	assetDefinition =
// 	// 		string AuthenticationReference = 1;
// 	// 		bytes BeneficiaryECPublicKey   = 2;
// 	// 		bytes SikePublicKey            = 3;
// 	// 		bytes BLSPublicKey             = 4;
// 	// 		int64 Timestamp                = 5;
// 	// 	transfers = nil
// 	// 	signature = sig of entire message

// 	// //Create Wallet
// 	// 	type = 2
// 	// 	ID = new ID()
// 	// 	owner = self ID
// 	// 	assetDefinition =
// 	// 		wallet =

// 	// 	transfers = nil
// 	// 	signature = sig of entire message

// 	// //Create Trustee Group
// 	print("DONE")
// }

// func DumpStore() {
// 	for key, value := range store {
// 		pp, _ := prettyjson.Marshal(value)
// 		fmt.Printf("%v - %v", key, string(pp))
// 	}
// }

// func Store(msg proto.Message) string {
// 	key := makeKeyString(msg)
// 	store[key] = msg
// 	return key
// }

// //This is design to create a unique key for each object for testing purposes
// //It emulates the TX id when placed in a blockchain
// func makeKeyString(msg proto.Message) string {
// 	key, _ := proto.Marshal(msg)
// 	res := sha256.Sum256(key)
// 	return hex.EncodeToString(res[:])
// }

// //CreateID - Create a new DocID along with keys
// func CreateID(authenticationReference string) (IDDocDeclaration *protobuffer.AssetDeclaration, seed []byte, err error) {

// 	//generate crypto random seed
// 	seed, err = cryptowallet.RandomBytes(48)
// 	if err != nil {
// 		err = errors.Wrap(err, "Failed to generate random seed")
// 		return nil, nil, err
// 	}

// 	sikePublicKey, _, err := keystore.GenerateSIKEKeys(seed)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	blsPublicKey, _, err := keystore.GenerateBLSKeys(seed)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	ecPublicKey, err := keystore.GenerateECPublicKey(seed)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// build ID Doc
// 	idDocument := &protobuffer.PBAssetDeclaration{}
// 	asset := &protobuffer.AssetDeclaration_Iddoc{}
// 	asset.Iddoc = &protobuffer.PBIDDoc{}
// 	asset.Iddoc.AuthenticationReference = authenticationReference
// 	asset.Iddoc.BeneficiaryECPublicKey = ecPublicKey
// 	asset.Iddoc.SikePublicKey = sikePublicKey
// 	asset.Iddoc.BLSPublicKey = blsPublicKey
// 	idDocument.Payload = asset
// 	return idDocument, seed, nil
// }

// func CreateWallet(iddoc string) (*protobuffer.AssetDeclaration, error) {
// 	wallet := &protobuffer.AssetDeclaration{}
// 	asset := &protobuffer.AssetDeclaration_Wallet{}
// 	asset.Wallet = &protobuffer.Wallet{}
// 	wallet.AssetDefinition = asset

// 	return wallet, nil
// }

// // func UpdateWallet() {
// // }

// // func DeclareTrusteeGroup() {
// // }

// // func UpdateTrusteeGroup() {
// // }

// // func RetrieveIDDoc(key string) {}
// // func RetrieveWallet(key string) {}
// // func RetrieveTrusteeGroup(key string){}

// // func VerifyTransaction(key string) (bool, error){}
// // func VerifyIDDoc(key string) (bool, error){}
