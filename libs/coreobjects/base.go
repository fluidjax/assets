package coreobjects

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

type Asset struct {
	protobuffer.Signature
	Description
	Authenticator
	store *Mapstore
	seed  []byte
	key   []byte
}

//Authenticator
type Authenticator struct {
}

type AuthenticatorInterface interface {
	Serialize(as interface{}) ([]byte, error)
	Sign(iddoc IDDoc) error
	SelfSign() error
	Verify() (bool, error)
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

func (a *Asset) Save() error {
	store := a.store
	msg := a.Signature
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	store.Save(a.key, data)
	return nil
}

func (a Asset) Load(key []byte) (*protobuffer.Signature, error) {
	store := a.store
	val, err := store.Load(key)
	if err != nil {
		return nil, err
	}
	msg := &protobuffer.Signature{}
	err = proto.Unmarshal(val, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
