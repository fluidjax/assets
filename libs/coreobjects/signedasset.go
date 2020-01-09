package coreobjects

import (
	"bytes"
	"crypto/sha256"
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

//AssetKeyFromPayloadHash - set the Assets ID Key to be sha256 of the Serialized Payload
func (a *SignedAsset) AssetKeyFromPayloadHash() (err error) {
	data, err := a.SerializePayload()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	a.SetKey(res[:])
	return nil
}

//Key Return the AssetKey
func (a *SignedAsset) Key() []byte {
	return a.Asset.GetID()
}

//SetKey - set the asset Key
func (a *SignedAsset) SetKey(key []byte) {
	a.Asset.ID = key
}

//Save - write the entire Signed Asset to the store
func (a *SignedAsset) Save() error {
	store := a.store
	data, err := proto.Marshal(a)
	if err != nil {
		return err
	}
	store.Save(a.Key(), data)
	return nil
}

//Load - read a SignedAsset from the store
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

//Pretty print the Asset for debugging
func (a *SignedAsset) Dump() {
	pp, _ := prettyjson.Marshal(a)
	fmt.Printf("%v", string(pp))
}

/* AddTransfer
Add a new Transfer/Update rule to the Asset's Transferlist
transferType	enum of transfer type such as settlePush, swap, Transfer etc.
expression 		is a string containing the boolean expression such as "t1 + t2 + t3 > 1 & p"
participant		map of abbreviation:IDDocKey		eg. t1 : b51de57554c7a49004946ec56243a70e90a26fbb9457cb2e6845f5e5b3c69f6a
*/
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

/* IsValidTransfer
Calculates if the boolean expression in the asset has been satisfied by the supplied signatures
transferSignatures = array of SignatureID  - SignatureID{IDDoc: [&IDDoc{}], Abbreviation: "p", Signature: [BLSSig]}
*/
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
	expression, _ = resolveExpression(expression, transfer.Participants, transferSignatures)
	result := boolparser.BoolSolve(expression)
	return result, nil
}

/*resolveExpression
Resolve the asset's expression by substituting 0's & 1's depending on whether signatures are available for each abbreviation/participant
expression 		is a string containing the boolean expression such as "t1 + t2 + t3 > 1 & p"
participantis	map of abbreviation:IDDocKey		eg. t1 : b51de57554c7a49004946ec56243a70e90a26fbb9457cb2e6845f5e5b3c69f6a
transferSignatures = array of SignatureID  - SignatureID{IDDoc: [&IDDoc{}], Abbreviation: "p", Signature: [BLSSig]}
*/
func resolveExpression(expression string, participants map[string][]byte, transferSignatures []SignatureID) (expressionOut string, display string) {
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

/*
TruthTable
For the supplied TransferType iterate through every combination of the existence or not of a participants signature.
Every possible matchining combination is returned where that combination will result in an asset Transfer
Unecessary abbreviations(Participants) are marked as 0s
Required signatures are marked with their abbreviation
eg.  [ 0 + t2 + t3 > 1 & p] = Transfer will occur if Trustee2, Trustee3 & Principals Signatures are present
	 [t1 + 0 + t3 > 1 & p]  = Transfer will occur if Trustee1, Trustee3 & Principals Signatures are present
*/
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
		var transferSignatures []SignatureID
		for i := 0; i < totalParticipants; i++ {
			pos := int64(math.Pow(2, float64(i)))
			val := j & pos
			iddoc := participantArray[i].IDDoc

			if val == 0 {
				transferSignatures = append(transferSignatures, SignatureID{IDDoc: iddoc, Signature: nil})
			} else {
				transferSignatures = append(transferSignatures, SignatureID{IDDoc: iddoc, Signature: []byte("hello")})
			}
		}
		resolvedExpression, display := resolveExpression(expression, transfer.Participants, transferSignatures)
		result := boolparser.BoolSolve(resolvedExpression)
		if result == true {
			matchedTrue = append(matchedTrue, display)
		}
	}
	sort.Strings(matchedTrue)
	return matchedTrue, nil
}

/*
SignPayload
returns the BLS signature of the serialize payload, signed with the BLS Private key of the supplied IDDoc
note the IDDoc must contain the seed
*/
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

/*
VerifyPayload
verifies the supplied signature with supplied IDDoc's BLS Public Key
note the IDDoc seed is NOT required
*/
func (a *SignedAsset) VerifyPayload(signature []byte, i *IDDoc) (verify bool, err error) {
	//Message
	data, err := a.SerializePayload()
	if err != nil {
		return false, errors.New("Failed to serialize payload")
	}
	//Public Key
	payload := i.IDDocPayload()
	blsPK := payload.GetBLSPublicKey()

	rc := crypto.BLSVerify(data, blsPK, signature)
	if rc == 0 {
		return true, nil
	}
	return false, nil
}

/*
SerializePayload
serialize the Assets payload (oneof) into a byte
*/
func (a *SignedAsset) SerializePayload() (s []byte, err error) {
	s, err = proto.Marshal(a.PBSignedAsset.Asset)
	if err != nil {
		s = nil
	}
	return s, err
}

/*
AggregatedSign
Aggregates BLSPubKeys and BLSSignatures from supplied array of SignatureIDs
Results are inserted into the object
only error is returned
*/
func (a *SignedAsset) AggregatedSign(transferSignatures []SignatureID) error {
	if transferSignatures == nil || len(transferSignatures) == 0 {
		return errors.New("Invalid transferSignatures BLS Aggregation")
	}
	var aggregatedSig []byte
	rc := 0
	var aggregatedPublicKey []byte
	var signers []string
	for i := 0; i < len(transferSignatures); i++ {
		sig := transferSignatures[i].Signature
		pubKey := transferSignatures[i].IDDoc.GetAsset().GetIddoc().GetBLSPublicKey()
		abbreviation := transferSignatures[i].Abbreviation
		signers = append(signers, abbreviation)
		if sig == nil || pubKey == nil {
			continue
		}
		if i == 0 {
			aggregatedSig = sig
			aggregatedPublicKey = pubKey
		} else {
			rc, aggregatedSig = crypto.BLSAddG1(aggregatedSig, sig)
			if rc != 0 {
				return errors.New("BLSAddG1 failed in AggregatedSign")
			}
			rc, aggregatedPublicKey = crypto.BLSAddG2(aggregatedPublicKey, pubKey)
			if rc != 0 {
				return errors.New("BLSAddG2 failed in AggregatedSign")
			}
		}
	}
	a.PublicKey = aggregatedPublicKey
	a.Signature = aggregatedSig
	a.Signers = signers
	data, err := a.SerializePayload()
	if err != nil {
		return errors.Wrap(err, "Fail to Aggregated Signatures")
	}
	rc = crypto.BLSVerify(data, aggregatedPublicKey, aggregatedSig)
	if rc != 0 {
		return errors.New("Signature failed to Verify")
	}
	return nil
}

/*
FullVerify
Based on the previous Asset state, retrieve the IDDocs of all signers.
publickeys = Aggregated the signers BLS Public Keys
message = Create a Serialized Payload

Using these fields verify the Signature in the transfer.

*/
func (a *SignedAsset) FullVerify(previousAsset *protobuffer.PBSignedAsset) (bool, error) {
	transferType := a.Asset.TransferType
	var transferSignatures []SignatureID

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

		sigID := SignatureID{IDDoc: iddoc, Abbreviation: abbreviation, Signature: []byte("UNKNOWN")}
		transferSignatures = append(transferSignatures, sigID)
	}

	//check the one in the object matches the one just created
	//Todo: We could probably remove the one in the object?
	res := bytes.Compare(aggregatedPublicKey, a.GetPublicKey())
	if res != 0 {
		return false, errors.New("Generated Aggregated Public Key doesnt match the one used to sign")
	}

	//Get Message
	data, err := a.SerializePayload()
	if err != nil {
		return false, errors.New("Failed to serialize payload")
	}
	//Retrieve the Sig
	aggregatedSig := a.GetSignature()

	//Verify
	rc := crypto.BLSVerify(data, aggregatedPublicKey, aggregatedSig)

	if rc != 0 {
		return false, errors.New("Signature failed to Verify")
	}

	//As the Signature Verified, we know that each  Signature:[]byte("UNKNOWN") which was used to generate the aggregate Signature
	//is valid, but currently unknown, we have enough info to do an IsValidTransfer

	validTransfer, err := a.IsValidTransfer(transferType, transferSignatures)
	if err != nil {
		return false, errors.Wrap(err, "Fail to fullVerify asset - checking IsValidTransfer failed")
	}
	if validTransfer == false {
		return false, errors.New("Invalid transfer - insufficient signatures")
	}

	return true, nil
}
