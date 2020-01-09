package coreobjects

import "github.com/qredo/assets/libs/protobuffer"

//Core Heirachcy
type Wallet struct {
	SignedAsset
}

type IDDoc struct {
	SignedAsset
}

type SignedAsset struct {
	protobuffer.PBSignedAsset
	store *Mapstore //Reference to object store (map or blockchain)
	seed  []byte    //If available a seed to generate keys for object
	//key           []byte                     //
	previousAsset *protobuffer.PBSignedAsset //Reference to (if any) previous object with the same key
}

//Use to hold ID & Signatures for expression parsing
type SignatureID struct {
	IDDoc        *IDDoc
	Abbreviation string
	HaveSig      bool
	Signature    []byte
}

type TransferParticipant struct {
	IDDoc        *IDDoc
	Abbreviation string
}
