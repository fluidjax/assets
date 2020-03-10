/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package assets

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/boolparser"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/prettyjson"
	"github.com/qredo/assets/libs/protobuffer"
)

func (a *SignedAsset) DeepCopyUpdatePayload() {
	//Deep copy the old Payload to the new one
	copier.Copy(a.CurrentAsset.Asset.Payload, a.PreviousAsset.Asset.Payload)
}

// Sign - this Signs the Asset including the Payload
func (a *SignedAsset) Sign(iddoc *IDDoc) error {
	if a == nil {
		return errors.New("SignedAsset is nil")
	}
	if iddoc == nil {
		return errors.New("Sign - supplied IDDoc is nil")
	}
	msg, err := a.SerializeAsset()
	if err != nil {
		return errors.Wrap(err, "Failed to Marshall Asset in Sign")
	}
	if iddoc.Seed == nil {
		return errors.New("Unable to Sign IDDoc - No Seed")
	}
	_, blsSK, err := keystore.GenerateBLSKeys(iddoc.Seed)
	if err != nil {
		return err
	}
	rc, signature := crypto.BLSSign(msg, blsSK)
	if rc != 0 {
		return errors.New("Failed to Sign Asset")
	}
	a.CurrentAsset.Signature = signature

	return nil
}

// Verify the Signature of the Asset (including the Payload)
func (a *SignedAsset) Verify(iddoc *IDDoc) error {

	//Check 2
	if a == nil {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:VerifySignedAsset:Signed Asset is nil")
	}
	if a.CurrentAsset == nil {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:VerifySignedAsset:Current Asset is nil")
	}
	if a.CurrentAsset.Signature == nil {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:VerifySignedAsset:Invalid Signature")
	}
	if iddoc == nil {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:VerifySignedAsset:IDDoc is nil")
	}
	msg, err := a.SerializeAsset()
	if err != nil {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:VerifySignedAsset:Fail to Serialize Asset")
	}
	payload, err := iddoc.Payload()
	if err != nil {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:VerifySignedAsset:Fail to Retrieve Payload")
	}

	//Check 3
	blsPK := payload.GetBLSPublicKey()
	rc := crypto.BLSVerify(msg, blsPK, a.CurrentAsset.Signature)
	if rc != 0 {
		return NewAssetsError(CodeConsensusErrorFailtoVerifySignature, "Consensus:Error:Check:VerifySignedAsset:Invalid Signature")
	}
	return nil
}

// Key Return the AssetKey
func (a *SignedAsset) Key() []byte {
	if a == nil {
		return nil
	}
	return a.CurrentAsset.Asset.GetID()
}

// Save - write the entire Signed Asset to the store
func (a *SignedAsset) Save() (string, error) {
	if a == nil {
		return "", errors.New("SignedAsset is nil")
	}
	store := a.DataStore
	data, err := a.SerializeSignedAsset()
	if err != nil {
		return "", err
	}
	return store.Set(a.Key(), data)

}

// Load - read a SignedAsset from the store
func Load(store DataSource, key []byte) (*protobuffer.PBSignedAsset, error) {
	val, err := store.RawGet(key)
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







// Dump - Pretty print the Asset for debugging
func (a *SignedAsset) Dump() {
	if a == nil {
		return
	}
	pp, _ := prettyjson.Marshal(a)
	fmt.Printf("%v", string(pp))
}

// AddTransfer - Add a new Transfer/Update rule to the Asset's Transferlist
// transferType	enum of transfer type such as settlePush, swap, Transfer etc.
// expression 		is a string containing the boolean expression such as "t1 + t2 + t3 > 1 & p"
// participant		map of abbreviation:IDDocKey		eg. t1 : b51de57554c7a49004946ec56243a70e90a26fbb9457cb2e6845f5e5b3c69f6a
func (a *SignedAsset) AddTransfer(transferType protobuffer.PBTransferType, expression string, participants *map[string][]byte, description string) error {
	if a == nil {
		return errors.New("AddTransfer is nil")
	}
	if expression == "" {
		return errors.New("AddTransfer - expression is empty")
	}
	if participants == nil {
		return errors.New("AddTransfer - participants is empty")
	}
	if len(*participants) == 0 {
		return errors.New("AddTransfer - participants is 0 length")
	}
	transferRule := &protobuffer.PBTransfer{}
	transferRule.Type = transferType
	transferRule.Expression = expression
	transferRule.Description = description
	if transferRule.Participants == nil {
		transferRule.Participants = make(map[string][]byte)
	}
	for abbreviation, iddocID := range *participants {
		transferRule.Participants[abbreviation] = iddocID
	}
	//Cant use enum as map key, so convert to a string
	transferListMapString := transferType.String()
	if a.CurrentAsset.Asset.Transferlist == nil {
		a.CurrentAsset.Asset.Transferlist = make(map[string]*protobuffer.PBTransfer)
	}
	a.CurrentAsset.Asset.Transferlist[transferListMapString] = transferRule
	return nil
}

// IsValidTransfer - Calculates if the boolean expression in the asset has been satisfied by the supplied signatures
// transferSignatures = array of SignatureID  - SignatureID{IDDoc: [&IDDoc{}], Abbreviation: "p", Signature: [BLSSig]}
func (a *SignedAsset) IsValidTransfer(transferType protobuffer.PBTransferType, transferSignatures []SignatureID) (bool, error) {
	if a == nil {
		return false, errors.New("IsValidTransfer - SignedAsset is nil")
	}
	if transferSignatures == nil {
		return false, errors.New("IsValidTransfer - transferSignatures is empty")
	}
	if len(transferSignatures) == 0 {
		return false, errors.New("IsValidTransfer - transferSignatures is 0 length")
	}
	transferListMapString := transferType.String()
	PreviousAsset := a.PreviousAsset
	if PreviousAsset == nil {
		return false, errors.New("No Previous Asset to change")
	}
	transfer := PreviousAsset.Asset.Transferlist[transferListMapString]
	if transfer == nil {
		return false, errors.New("No Transfer Found")
	}
	expression := transfer.Expression
	expression, _, err := resolveExpression(a.DataStore, expression, transfer.Participants, transferSignatures, "")
	if err != nil {
		return false, err
	}
	result := boolparser.BoolSolve(expression)
	return result, nil
}

// resolveExpression - Resolve the asset's expression by substituting 0's & 1's depending on whether signatures are available for each abbreviation/participant
// expression 		is a string containing the boolean expression such as "t1 + t2 + t3 > 1 & p"
// participantis	map of abbreviation:IDDocKey		eg. t1 : b51de57554c7a49004946ec56243a70e90a26fbb9457cb2e6845f5e5b3c69f6a
// transferSignatures = array of SignatureID  - SignatureID{IDDoc: [&IDDoc{}], Abbreviation: "p", Signature: [BLSSig]}
// recursionPrefix    = initally called empty, recursion appends sub objects eg. "tg1.x1" for participant x1 in tg1 Group
func resolveExpression(store DataSource, expression string, participants map[string][]byte, transferSignatures []SignatureID, prefix string) (expressionOut string, display string, err error) {
	if expression == "" {
		return "", "", errors.New("resolveExpression - expression is empty")
	}
	if participants == nil {
		return "", "", errors.New("resolveExpression - participants is empty")
	}
	if len(participants) == 0 {
		return "", "", errors.New("resolveExpression - participants is 0 length")
	}
	if transferSignatures == nil {
		return "", "", errors.New("resolveExpression - transferSignatures is empty")
	}
	if len(transferSignatures) == 0 {
		return "", "", errors.New("resolveExpression - transferSignatures is 0 length")
	}
	expressionOut = expression
	display = expression
	//Loop all the participants from previous Asset
	for abbreviation, id := range participants {
		found := false
		participant, err := Load(store, id)
		if err != nil {
			return "", "", errors.New("Fail to retrieve participant")
		}
		switch participant.GetAsset().GetPayload().(type) {
		case *protobuffer.PBAsset_Group:
			if participant.Asset.Type != protobuffer.PBAssetType_Group {
				break
			}
			Group, err := ReBuildGroup(participant, id)
			if err != nil {
				return "", "", errors.Wrap(err, "Failed to Rebuild  Group in resolve Expression")
			}
			recursionExpression := string(Group.Payload().GetGroupFields()["expression"])
			recursionParticipants := Group.Payload().Participants
			recursionPrefix := prefix + abbreviation + "."
			subExpression, subDisplay, err := resolveExpression(store, recursionExpression, recursionParticipants, transferSignatures, recursionPrefix)
			if err != nil {
				return "", "", errors.Wrap(err, "Failed to Resolve Recursive Expression")
			}
			expressionOut = strings.ReplaceAll(expressionOut, abbreviation, " ( "+subExpression+" ) ")
			display = strings.ReplaceAll(display, abbreviation, " ( "+subDisplay+" ) ")

		case *protobuffer.PBAsset_Iddoc:
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
	}
	return expressionOut, display, nil
}

// TruthTable - For the supplied TransferType iterate through every combination of the existence or not of a participants signature.
// Every possible matchining combination is returned where that combination will result in an asset Transfer
// Unecessary abbreviations(Participants) are marked as 0s
// Required signatures are marked with their abbreviation
// eg.  [ 0 + t2 + t3 > 1 & p] = Transfer will occur if 2, 3 & Principals Signatures are present
// [t1 + 0 + t3 > 1 & p]  = Transfer will occur if 1, 3 & Principals Signatures are present
func (a *SignedAsset) TruthTable(transferType protobuffer.PBTransferType) ([]string, error) {
	if a == nil {
		return nil, errors.New("TruthTable - SignedAsset is nil")
	}
	transferListMapString := transferType.String()
	transfer := a.CurrentAsset.Asset.Transferlist[transferListMapString]
	if transfer == nil {
		return nil, errors.New("No Transfer Found")
	}
	expression := transfer.Expression
	totalParticipants := len(transfer.Participants)
	var participantArray []TransferParticipant
	for key, idkey := range transfer.Participants {
		idsig, err := Load(a.DataStore, idkey)
		if err != nil {
			return nil, errors.New("Failed to load iddoc " + hex.EncodeToString(idkey))
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
		resolvedExpression, display, err := resolveExpression(a.DataStore, expression, transfer.Participants, transferSignatures, "")
		if err != nil {
			return nil, err
		}
		result := boolparser.BoolSolve(resolvedExpression)
		if result == true {
			matchedTrue = append(matchedTrue, display)
		}
	}
	sort.Strings(matchedTrue)
	return matchedTrue, nil
}

// Sign - generic Sign Function
func Sign(msg []byte, iddoc *IDDoc) (signature []byte, err error) {
	if iddoc == nil {
		return nil, errors.New("Sign - supplied IDDoc is nil")
	}
	if iddoc.Seed == nil {
		return nil, errors.New("Unable to Sign IDDoc - No Seed")
	}
	_, blsSK, err := keystore.GenerateBLSKeys(iddoc.Seed)
	if err != nil {
		return nil, err
	}
	rc, signature := crypto.BLSSign(msg, blsSK)
	if rc != 0 {
		return nil, errors.New("Failed to Sign Asset")
	}
	return signature, nil
}

// SignAsset - returns the BLS signature of the serialize payload, signed with the BLS Private key of the supplied IDDoc
// note the IDDoc must contain the seed
func (a *SignedAsset) SignAsset(i *IDDoc) (s []byte, err error) {
	if a == nil {
		return nil, errors.New("SignAsset - SignedAsset is nil")
	}
	if i == nil {
		return nil, errors.New("Verify - supplied IDDoc is nil")
	}
	msg, err := a.SerializeAsset()
	if err != nil {
		return nil, errors.New("Failed to serialize payload")
	}
	signature, err := Sign(msg, i)
	return signature, err
}

// Verify - generic verify function
func Verify(msg []byte, signature []byte, iddoc *IDDoc) (bool, error) {
	idDocPayload, err := iddoc.Payload()
	if err != nil {
		return false, err
	}
	blsPK := idDocPayload.GetBLSPublicKey()
	rc := crypto.BLSVerify(msg, blsPK, signature)
	if rc == 0 {
		return true, nil
	}
	return false, nil
}

// VerifyAsset - verifies the supplied signature with supplied IDDoc's BLS Public Key
// note the IDDoc seed is NOT required
func (a *SignedAsset) VerifyAsset(signature []byte, i *IDDoc) (verify bool, err error) {
	if a == nil {
		return false, errors.New("VerifyAsset - SignedAsset is nil")
	}
	if i == nil {
		return false, errors.New("VerifyAsset - supplied IDDoc is nil")
	}
	if signature == nil {
		return false, errors.New("VerifyAsset - signature is nil")
	}
	//Message
	msg, err := a.SerializeAsset()
	if err != nil {
		return false, errors.New("Failed to serialize payload")
	}
	return Verify(msg, signature, i)
}

// AggregatedSign  - Aggregates BLSPubKeys and BLSSignatures from supplied array of SignatureIDs
// Results are inserted into the object
// only error is returned
func (a *SignedAsset) AggregatedSign(transferSignatures []SignatureID) error {
	if a == nil {
		return errors.New("AggregatedSign - SignedAsset is nil")
	}
	if transferSignatures == nil {
		return errors.New("AggregatedSign - signature is nil")
	}
	if transferSignatures == nil || len(transferSignatures) == 0 {
		return errors.New("Invalid transferSignatures BLS Aggregation")
	}
	var aggregatedSig []byte
	rc := 0
	var aggregatedPublicKey []byte
	signers := make(map[string][]byte)
	for i := 0; i < len(transferSignatures); i++ {
		sig := transferSignatures[i].Signature
		pubKey := transferSignatures[i].IDDoc.CurrentAsset.GetAsset().GetIddoc().GetBLSPublicKey()
		idkey := transferSignatures[i].IDDoc.Key()
		abbreviation := transferSignatures[i].Abbreviation
		signers[abbreviation] = idkey
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
	a.CurrentAsset.PublicKey = aggregatedPublicKey
	a.CurrentAsset.Signature = aggregatedSig
	a.CurrentAsset.Signers = signers
	data, err := a.SerializeAsset()
	if err != nil {
		return errors.Wrap(err, "Fail to Aggregated Signatures")
	}
	// fmt.Println("DATA:", hex.EncodeToString(data))
	// fmt.Println("KEY:", hex.EncodeToString(aggregatedPublicKey))
	// fmt.Println("SIG:", hex.EncodeToString(aggregatedSig))

	rc = crypto.BLSVerify(data, aggregatedPublicKey, aggregatedSig)
	if rc != 0 {
		return errors.New("Signature failed to Verify")
	}
	return nil
}

// buildSigKeys - Aggregated the signatures and public keys for all Participants
func buildSigKeys(store *DataSource, signers []string, currentTransfer *protobuffer.PBTransfer, aggregatedPublicKey []byte, transferSignatures []SignatureID) ([]SignatureID, []byte, error) {
	//For each supplied signer re-build a PublicKey
	for _, abbreviation := range signers {
		participantID := currentTransfer.Participants[abbreviation]
		signedAsset, err := Load(*store, participantID)
		if err != nil {
			return nil, nil, errors.New("Failed to retieve IDDoc")
		}
		switch signedAsset.GetAsset().GetPayload().(type) {
		case *protobuffer.PBAsset_Group:
			//print("recurse")
		case *protobuffer.PBAsset_Iddoc:
			iddoc, err := ReBuildIDDoc(signedAsset, participantID)
			if err != nil {
				return nil, nil, errors.Wrap(err, "Fail to obtain public Key in FullVerify")
			}
			pubKey := iddoc.CurrentAsset.GetAsset().GetIddoc().GetBLSPublicKey()
			if aggregatedPublicKey == nil {
				aggregatedPublicKey = pubKey
			} else {
				_, aggregatedPublicKey = crypto.BLSAddG2(aggregatedPublicKey, pubKey)
			}
			sigID := SignatureID{IDDoc: iddoc, Abbreviation: abbreviation, Signature: []byte("UNKNOWN")}
			transferSignatures = append(transferSignatures, sigID)
		}
	}
	return transferSignatures, aggregatedPublicKey, nil
}

// FullVerify - Based on the previous Asset state, retrieve the IDDocs of all signers.
// publickeys = Aggregated the signers BLS Public Keys
// message = Create a Serialized Payload
// Using these fields verify the Signature in the transfer.
func (a *SignedAsset) FullVerify() (bool, error) {
	previousAsset := a.PreviousAsset
	if a == nil {
		return false, errors.New("FullVerify - SignAsset is nil")
	}
	transferType := a.CurrentAsset.Asset.TransferType
	var transferSignatures []SignatureID
	var aggregatedPublicKey []byte
	if previousAsset == nil {
		return false, errors.New("No Previous Asset supplied for Verify")
	}
	transferList := previousAsset.GetAsset().GetTransferlist()
	_ = transferList

	//For each supplied signer re-build a PublicKey
	for abbreviation, participantID := range a.CurrentAsset.Signers {
		signedAsset, err := Load(a.DataStore, participantID)
		if err != nil {
			return false, errors.New("Failed to retieve IDDoc")
		}
		iddoc, err := ReBuildIDDoc(signedAsset, participantID)
		if err != nil {
			return false, errors.Wrap(err, "Fail to obtain public Key in FullVerify")
		}
		pubKey := iddoc.CurrentAsset.GetAsset().GetIddoc().GetBLSPublicKey()
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
	res := bytes.Compare(aggregatedPublicKey, a.CurrentAsset.GetPublicKey())
	if res != 0 {
		return false, errors.New("Generated Aggregated Public Key doesnt match the one used to sign")
	}
	//Get Message
	data, err := a.SerializeAsset()
	if err != nil {
		return false, errors.New("Failed to serialize payload")
	}
	//Retrieve the Sig
	aggregatedSig := a.CurrentAsset.GetSignature()
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

// assetKeyFromPayloadHash - set the Assets ID Key to be sha256 of the Serialized Payload
func (a *SignedAsset) AssetKeyFromPayloadHash() (err error) {
	if a == nil {
		return errors.New("AssetKeyFromPayloadHash - SignAsset is nil")
	}
	data, err := a.SerializeAsset()
	if err != nil {
		return err
	}
	res := sha256.Sum256(data)
	a.setKey(res[:])
	return nil
}

// SetKey - set the asset Key
func (a *SignedAsset) setKey(key []byte) {
	if a == nil {
		return
	}
	a.CurrentAsset.Asset.ID = key
}

// SerializeAsset - Serialize the Asset (PBAsset, not PBSignedAsset)
func (a *SignedAsset) SerializeAsset() (s []byte, err error) {
	if a == nil {
		return nil, errors.New("SerializeAsset - SignAsset is nil")
	}
	if a.CurrentAsset.Asset == nil {
		return nil, errors.New("Can't serialize nil payload")
	}

	s, err = proto.Marshal(a.CurrentAsset.Asset)
	if err != nil {
		s = nil
	}
	return s, err
}

// SerializeSignedAsset - Serialize the entire PBSignedAsset
func (a *SignedAsset) SerializeSignedAsset() (s []byte, err error) {
	if a == nil {
		return nil, errors.New("SerializeSignedAsset - SignAsset is nil")
	}
	s, err = proto.Marshal(a.CurrentAsset)
	if err != nil {
		s = nil
	}

	return s, err
}

// Hash is for debugging
func (a *SignedAsset) Hash() string {
	if a == nil {
		return ""
	}
	data, err := a.SerializeAsset()
	if err != nil {
		return ""
	}
	result := sha256.Sum256(data)
	return hex.EncodeToString(result[:])
}

func (a *SignedAsset) AddTag(key string, value []byte) {
	if a.CurrentAsset.Asset.Tags == nil {
		a.CurrentAsset.Asset.Tags = map[string][]byte{}
	}
	a.CurrentAsset.Asset.Tags[key] = value

}

func (a *SignedAsset) GetWithSuffix(datasource DataSource, key []byte, suffix string) ([]byte, error) {
	fullSuffix := []byte(suffix)
	key = append(key[:], fullSuffix[:]...)
	return datasource.RawGet(key)
}

func (a *SignedAsset) SetWithSuffix(datasource DataSource, key []byte, suffix string, data []byte) error {
	suffixBytes := []byte(suffix)
	fullkey := append(key[:], suffixBytes[:]...)
	// println("SET1 ", hex.EncodeToString(key))
	// println("SET2 ", hex.EncodeToString(data))
	// println("SET3 ", hex.EncodeToString(fullkey))
	_, err := datasource.Set(fullkey, data)
	return err
}

func (a *SignedAsset) Exists(datasource DataSource, key []byte) (bool, error) {
	item, err := datasource.RawGet(key)
	return item != nil, err
}

func (a *SignedAsset) BatchExists(datasource DataSource, key []byte) (bool, error) {
	item, err := datasource.BatchGet(key)
	return item != nil, err
}

func (a *SignedAsset) AddCoreMappings(datasource DataSource, rawTX []byte, txHash []byte) (err error) {
	_, err = datasource.Set(txHash, rawTX)
	if err != nil {
		return err
	}
	_, err = datasource.Set(a.Key(), txHash)
	if err != nil {
		return err
	}
	return nil
}

func (a *SignedAsset) subtractFromBalanceKey(datasource DataSource, assetID []byte, amount int64) error {
	currentBalance, assetsError := a.getBalanceKey(datasource, assetID)
	if assetsError != nil {
		return assetsError
	}

	newBalance := currentBalance - amount
	if newBalance < 0 {
		return NewAssetsError(CodeConsensusInsufficientFunds, "Consensus:Error:Check:Balance:Newbalance is less than Zero")
	}
	return a.setBalanceKey(datasource, assetID, newBalance)
}

func (a *SignedAsset) addToBalanceKey(datasource DataSource, assetID []byte, amount int64) error {
	currentBalance, assetsError := a.getBalanceKey(datasource, assetID)
	if assetsError != nil {
		return assetsError
	}
	newBalance := currentBalance + amount
	return a.setBalanceKey(datasource, assetID, newBalance)
}
func (a *SignedAsset) setBalanceKey(datasource DataSource, assetID []byte, newBalance int64) error {
	//Convert new balance to bytes and save for AssetID
	newBalanceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(newBalanceBytes, uint64(newBalance))
	err := a.SetWithSuffix(datasource, assetID, ".balance", newBalanceBytes)
	if err != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Check:Balance:Fail to Set Balance Key")
	}
	return nil
}

func (a *SignedAsset) getBalanceKey(datasource DataSource, assetID []byte) (amount int64, assetError error) {
	currentBalanceBytes, err := a.GetWithSuffix(datasource, assetID, ".balance")
	if currentBalanceBytes == nil || err != nil {
		return 0, NewAssetsError(CodeDatabaseFail, "Consensus:Error:Check:Balance:Fail to Get Balance Key")
	}
	currentBalance := int64(binary.LittleEndian.Uint64(currentBalanceBytes))
	return currentBalance, nil
}
