package coreobjects

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/hokaccha/go-prettyjson"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/boolparser"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
)

//Use to hold ID & Signatures for expression parsing
type SignatureID struct {
	IDDoc     *IDDoc
	Signature []byte
}

type BaseAsset struct {
	protobuffer.Signature
	store         *Mapstore
	seed          []byte
	key           []byte
	previousAsset *BaseAsset
}

func (a *BaseAsset) SignPayload(i *IDDoc) (s []byte, err error) {
	data, err := a.PayloadSerialize()
	if err != nil {
		return nil, errors.New("Failed to serialize payload")
	}
	if i.seed == nil {
		return nil, errors.New("No Seed in Supplied IDDoc")
	}
	_, blsSK, err := keystore.GenerateBLSKeys(i.seed)
	if err != nil {
		return nil, err
	}
	rc, signature := crypto.BLSSign(data, blsSK)
	if rc != 0 {
		return nil, errors.New("Failed to sign IDDoc")
	}
	return signature, nil
}

func (a *BaseAsset) VerifyPayload(signature []byte, i *IDDoc) (verify bool, err error) {
	//Message
	data, err := a.PayloadSerialize()
	if err != nil {
		return false, errors.New("Failed to serialize payload")
	}
	//Public Key
	payload := i.AssetPayload()
	blsPK := payload.GetBLSPublicKey()

	rc := crypto.BLSVerify(data, blsPK, signature)
	if rc == 0 {
		return true, nil
	}
	return false, nil
}

func (a *BaseAsset) PayloadSerialize() (s []byte, err error) {
	s, err = proto.Marshal(a.Signature.Asset)
	if err != nil {
		s = nil
	}
	return s, err
}

func (a *BaseAsset) Save() error {
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
func (a *BaseAsset) SetTestKey() (err error) {
	data, err := a.PayloadSerialize()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	a.key = res[:]
	return nil
}

func (a *BaseAsset) Description() {
	print("Asset Description")
}

//Add a new Transfer/Update rule
//Specify the boolean expression & add list of participants
func (a *BaseAsset) AddTransfer(transferType protobuffer.TransferType, expression string, participants map[string][]byte) error {
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
	//Cant use enum as map key, so convert to a string
	transferListMapString := transferType.String()
	if ob.Transferlist == nil {
		ob.Transferlist = make(map[string]*protobuffer.Transfer)
	}
	ob.Transferlist[transferListMapString] = transferRule
	return nil
}

//Pretty print the Asset for debugging
func (a *BaseAsset) Dump() {
	pp, _ := prettyjson.Marshal(a.Signature)
	fmt.Printf("%v", string(pp))
}

//Given a list of signature build a sig map
func (a *BaseAsset) IsValidTransfer(transferType protobuffer.TransferType, transferSignatures []SignatureID) (bool, error) {
	transferListMapString := transferType.String()
	//sigmap := make(map[string][]byte)

	previousAsset := a.previousAsset
	if previousAsset == nil {
		return false, errors.New("No Previous Asset to change")
	}

	transfer := previousAsset.Asset.Transferlist[transferListMapString]
	if transfer == nil {
		return false, errors.New("No Transfer Found")
	}

	expression := transfer.Expression

	//Loop All Participants
	for abbreviation, id := range transfer.Participants {
		//Loop through transfer Signatures
		found := false
		for _, sigID := range transferSignatures {
			res := bytes.Compare(sigID.IDDoc.key, id)

			//fmt.Printf("TRY:  %v\n  %v\n  %v\n\n", res, sigID.IDDoc.key, id)

			if res == 0 && sigID.Signature != nil {
				//fmt.Printf("replace %v with 1\n", abbreviation)
				expression = strings.ReplaceAll(expression, abbreviation, "1")
				found = true
				break
			}
		}
		if found == false {
			//fmt.Printf("replace %v with 0\n", abbreviation)
			expression = strings.ReplaceAll(expression, abbreviation, "0")
		}
	}

	result := boolparser.BoolSolve(expression)

	fmt.Printf("%v %s \n", result, expression)

	return result, nil

}
