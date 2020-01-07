package coreobjects

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

type Asset struct {
	protobuffer.Signature
	MapDataStore
	Description
	Authenticator
	seed []byte
	key  []byte
}

//Authenticator
type Authenticator struct {
}

type AuthenticatorInterface interface {
	Serialize(as interface{}) ([]byte, error)
	Sign() error
	Verify() (bool, error)
}

//MapDataStore
type MapDataStore struct {
}

//Description
type Description struct {
}

func (a *Authenticator) Serialize(as interface{}) (s []byte, err error) {
	i := as.(Asset)
	switch i.Signature.GetSignatureAsset().(type) {
	case *protobuffer.Signature_Declaration:
		s, err = proto.Marshal(i.Signature.GetDeclaration())
	case *protobuffer.Signature_Update:
		s, err = proto.Marshal(i.Signature.GetUpdate())
	default:
		err = errors.New("Fail to Serialize Asset")
	}

	if err != nil {
		s = nil
	}

	return s, err
}
