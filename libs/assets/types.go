package assets

import "github.com/qredo/assets/libs/protobuffer"

//Core Heirachcy

//Group - Group Asset Type
type Group struct {
	SignedAsset
}

//Wallet - Wallet Asset Type
type Wallet struct {
	SignedAsset
}

//IDDoc - IDDoc Asset Type
type IDDoc struct {
	SignedAsset
}

//SignedAsset - Asset/Previous Asset Wrapper, holding temporary objects (seed) & previousVersions
type SignedAsset struct {
	currentAsset  *protobuffer.PBSignedAsset
	store         *Mapstore                  //Reference to object store (map or blockchain)
	seed          []byte                     //If available a seed to generate keys for object
	previousAsset *protobuffer.PBSignedAsset //Reference to (if any) previous object with the same key
}

//SignatureID - Use to hold ID & Signatures for expression parsing
type SignatureID struct {
	IDDoc        *IDDoc
	Abbreviation string
	HaveSig      bool
	Signature    []byte
}

//TransferParticipant -
type TransferParticipant struct {
	IDDoc        *IDDoc
	Abbreviation string
}
