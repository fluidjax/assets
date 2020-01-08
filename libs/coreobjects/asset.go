package coreobjects

import (
	"crypto/sha256"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hokaccha/go-prettyjson"
	"github.com/qredo/assets/libs/protobuffer"
)

//Use to hold ID & Signatures for expression parsing
type SignatureID struct {
	IDDocID   []byte
	Signature []byte
}

type Asset struct {
	protobuffer.Signature
	store *Mapstore
	seed  []byte
	key   []byte
}

func (a *Asset) PayloadSerialize() (s []byte, err error) {
	s, err = proto.Marshal(a.Signature.Asset)
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

func (a *Asset) Description() {
	print("Asset Description")
}

//Add a new Transfer/Update rule
//Specify the boolean expression & add list of participants
func (a *Asset) AddTransfer(transferType protobuffer.TransferType, expression string, participants map[string][]byte) error {
	transferRule := &protobuffer.Transfer{}
	transferRule.Type = transferType
	transferRule.Expression = expression

	if transferRule.Participants == nil {
		transferRule.Participants = make(map[string][]byte)
	}

	for abbreviation, iddocID := range participants {
		transferRule.Participants[abbreviation] = iddocID
	}
	ob := a.Signature.Asset
	ob.Transferlist = append(ob.Transferlist, transferRule)

	return nil
}

//Pretty print the Asset for debugging
func (a *Asset) Dump() {
	pp, _ := prettyjson.Marshal(a.Signature)
	fmt.Printf("%v", string(pp))
}
