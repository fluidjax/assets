package coreobjects

//Use to hold ID & Signatures for expression parsing
type SignatureID struct {
	IDDoc     *IDDoc
	Signature []byte
}

type TransferParticipant struct {
	IDDoc        *IDDoc
	Abbreviation string
}
