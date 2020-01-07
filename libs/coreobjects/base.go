package coreobjects

import (
	"crypto/sha256"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

type Asset struct {
	protobuffer.Signature
	store *Mapstore
	seed  []byte
	key   []byte
}

func (a *Asset) PayloadSerialize() (s []byte, err error) {
	switch a.Signature.GetSignatureAsset().(type) {
	case *protobuffer.Signature_Declaration:
		s, err = proto.Marshal(a.Signature.GetDeclaration())
	case *protobuffer.Signature_Update:
		s, err = proto.Marshal(a.Signature.GetUpdate())
	default:
		err = errors.New("Fail to PayloadSerialize Asset")
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

func Load(store *Mapstore, key []byte) (*protobuffer.Signature, error) {
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

//For testing only
func (a *Asset) SetTestKey() (err error) {
	data, err := a.PayloadSerialize()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	a.key = res[:]
	return nil
}
