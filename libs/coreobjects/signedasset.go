package coreobjects

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/hokaccha/go-prettyjson"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/boolparser"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
)

func (a *SignedAsset) Key() []byte {
	return a.Asset.GetID()
}

func (a *SignedAsset) SetKey(key []byte) {
	a.Asset.ID = key
}

func (a *SignedAsset) Save() error {
	store := a.store
	//msg := a.PBSignedAsset
	data, err := proto.Marshal(a)
	if err != nil {
		return err
	}
	store.Save(a.Key(), data)
	return nil
}

func Load(store *Mapstore, key []byte) (*protobuffer.PBSignedAsset, error) {
	val, err := store.Load(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, errors.New("Key not found")
	}
	msg := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(val, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

//For testing only
func (a *SignedAsset) SetTestKey() (err error) {
	data, err := a.SerializePayload()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	a.SetKey(res[:])
	return nil
}

//Pretty print the Asset for debugging
func (a *SignedAsset) Dump() {
	pp, _ := prettyjson.Marshal(a)
	fmt.Printf("%v", string(pp))
}

//Add a new Transfer/Update rule
//Specify the boolean expression & add list of participants
func (a *SignedAsset) AddTransfer(transferType protobuffer.PBTransferType, expression string, participants *map[string][]byte) error {
	transferRule := &protobuffer.PBTransfer{}
	transferRule.Type = transferType
	transferRule.Expression = expression
	if transferRule.Participants == nil {
		transferRule.Participants = make(map[string][]byte)
	}
	for abbreviation, iddocID := range *participants {
		transferRule.Participants[abbreviation] = iddocID
	}
	//Cant use enum as map key, so convert to a string
	transferListMapString := transferType.String()
	if a.PBSignedAsset.Asset.Transferlist == nil {
		a.PBSignedAsset.Asset.Transferlist = make(map[string]*protobuffer.PBTransfer)
	}
	a.PBSignedAsset.Asset.Transferlist[transferListMapString] = transferRule
	return nil
}

//Given a list of signature build a sig map
func (a *SignedAsset) IsValidTransfer(transferType protobuffer.PBTransferType, transferSignatures []SignatureID) (bool, error) {
	transferListMapString := transferType.String()
	previousAsset := a.previousAsset
	if previousAsset == nil {
		return false, errors.New("No Previous Asset to change")
	}
	transfer := previousAsset.Asset.Transferlist[transferListMapString]
	if transfer == nil {
		return false, errors.New("No Transfer Found")
	}
	expression := transfer.Expression
	expression, _ = ResolveExpression(expression, transfer.Participants, transferSignatures)
	result := boolparser.BoolSolve(expression)
	fmt.Printf("%v %s \n", result, expression)
	return result, nil
}

//Using the Specified participants change the abbreviations (t1, p etc) into boolean/int values
func ResolveExpression(expression string, participants map[string][]byte, transferSignatures []SignatureID) (expressionOut string, display string) {
	expressionOut = expression
	display = expression
	//Loop all the participants from previous Asset
	for abbreviation, id := range participants {
		found := false
		//Loop throught all the gathered signatures
		for _, sigID := range transferSignatures {
			res := bytes.Compare(sigID.IDDoc.Key(), id)
			if res == 0 && sigID.Signature != nil {
				//Where we have a signature for a given IDDoc, replace it with a 1
				expressionOut = strings.ReplaceAll(expressionOut, abbreviation, "1")
				found = true
				break
			}
		}
		if found == false {
			//Where we do not have signature for a given IDDoc, replace it with a 0
			display = strings.ReplaceAll(display, abbreviation, "0")
			expressionOut = strings.ReplaceAll(expressionOut, abbreviation, "0")
		}
	}
	return expressionOut, display
}

func (a *SignedAsset) TruthTable(transferType protobuffer.PBTransferType) ([]string, error) {
	transferListMapString := transferType.String()
	transfer := a.Asset.Transferlist[transferListMapString]
	if transfer == nil {
		return nil, errors.New("No Transfer Found")
	}
	expression := transfer.Expression

	totalParticipants := len(transfer.Participants)
	var participantArray []TransferParticipant
	for key, idkey := range transfer.Participants {
		idsig, err := Load(a.store, idkey)
		if err != nil {
			return nil, errors.New("Failed to load iddoc")
		}
		iddoc, err := ReBuildIDDoc(idsig, idkey)
		if err != nil {
			return nil, errors.New("Failed to Rebuild iddoc")
		}
		p := TransferParticipant{
			IDDoc:        iddoc,
			Abbreviation: key,
		}
		participantArray = append(participantArray, p)
	}

	var j int64
	var matchedTrue []string

	for j = 0; j < int64(math.Pow(2, float64(totalParticipants))); j++ {
		//fmt.Printf("%v:", j)
		var transferSignatures []SignatureID
		for i := 0; i < totalParticipants; i++ {
			pos := int64(math.Pow(2, float64(i)))
			val := j & pos
			//fmt.Printf("%v", val)
			iddoc := participantArray[i].IDDoc

			if val == 0 {
				transferSignatures = append(transferSignatures, SignatureID{IDDoc: iddoc, Signature: nil})
			} else {
				transferSignatures = append(transferSignatures, SignatureID{IDDoc: iddoc, Signature: []byte("hello")})
			}

		}
		resolvedExpression, display := ResolveExpression(expression, transfer.Participants, transferSignatures)
		result := boolparser.BoolSolve(resolvedExpression)

		if result == true {
			matchedTrue = append(matchedTrue, display)
		}
	}
	sort.Strings(matchedTrue)
	return matchedTrue, nil
}

//Payload

func (a *SignedAsset) SignPayload(i *IDDoc) (s []byte, err error) {
	data, err := a.SerializePayload()
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

func (a *SignedAsset) VerifyPayload(signature []byte, i *IDDoc) (verify bool, err error) {
	//Message
	data, err := a.SerializePayload()
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

func (a *SignedAsset) SerializePayload() (s []byte, err error) {
	s, err = proto.Marshal(a.PBSignedAsset.Asset)
	if err != nil {
		s = nil
	}
	return s, err
}

func (a *SignedAsset) Description() {
	print("Asset Description")
}

//BLS Aggregation
func (a *SignedAsset) AggregatedSign(transferSignatures []SignatureID) error {
	if transferSignatures == nil || len(transferSignatures) == 0 {
		return errors.New("Invalid transferSignatures BLS Aggregation")
	}
	var aggregatedSig []byte
	var aggregatedPublicKey []byte
	var signers []string
	for i := 0; i < len(transferSignatures); i++ {

		sig := transferSignatures[i].Signature
		pubKey := transferSignatures[i].IDDoc.GetAsset().GetIddoc().GetBLSPublicKey()
		//idDocKey := transferSignatures[i].IDDoc.Key()
		abbreviation := transferSignatures[i].Abbreviation

		signers = append(signers, abbreviation)

		if sig == nil || pubKey == nil {
			return errors.New("Invalid Sig/PubKey BLS Aggregation")
		}

		if i == 0 {
			aggregatedSig = sig
			aggregatedPublicKey = pubKey
		} else {
			_, aggregatedSig = crypto.BLSAddG1(aggregatedSig, sig)
			_, aggregatedPublicKey = crypto.BLSAddG2(aggregatedPublicKey, pubKey)
		}
	}
	a.PublicKey = aggregatedPublicKey
	a.Signature = aggregatedSig
	a.Signers = signers

	data, err := a.SerializePayload()
	print(err)

	rc := crypto.BLSVerify(data, aggregatedPublicKey, aggregatedSig)
	if rc != 0 {
		return errors.New("Signature failed to Verify")
	}

	return nil
}

//Verify the whole Asset
func (a *SignedAsset) FullVerify(previousAsset *protobuffer.PBSignedAsset) (bool, error) {
	//Get Transfer from previous
	transferType := a.Asset.TransferType
	//	aggregatedPublicKey := a.PublicKey
	//	aggregatedSig := a.Signature
	//	signers := a.Signers
	var aggregatedPublicKey []byte

	if previousAsset == nil {
		return false, errors.New("No Previous Asset supplied for Verify")
	}

	transferList := previousAsset.GetAsset().GetTransferlist()
	_ = transferList
	currentTransfer := transferList[transferType.String()]

	//For each supplied signer re-build a PublicKey
	for _, abbreviation := range a.Signers {
		participantIDDocID := currentTransfer.Participants[abbreviation]
		signedAsset, err := Load(a.store, participantIDDocID)

		if err != nil {
			return false, errors.New("Failed to retieve IDDoc")
		}

		iddoc, err := ReBuildIDDoc(signedAsset, participantIDDocID)

		pubKey := iddoc.GetAsset().GetIddoc().GetBLSPublicKey()
		if aggregatedPublicKey == nil {
			aggregatedPublicKey = pubKey
		} else {
			_, aggregatedPublicKey = crypto.BLSAddG2(aggregatedPublicKey, pubKey)
		}

	}

	//check the one used matches the one just created
	res := bytes.Compare(aggregatedPublicKey, a.GetPublicKey())

	if res != 0 {
		return false, errors.New("Generated Aggregated Public Key doesnt match the one used to sign")
	}

	//Check the Signature
	//Get Message
	data, err := a.SerializePayload()
	if err != nil {
		return false, errors.New("Failed to serialize payload")
	}
	//Sig
	aggregatedSig := a.GetSignature()

	rc := crypto.BLSVerify(data, aggregatedPublicKey, aggregatedSig)

	fmt.Println(hex.EncodeToString(data))

	if rc != 0 {
		return false, errors.New("Signature failed to Verify")
	}

	return true, nil
}
