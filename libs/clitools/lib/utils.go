package qc

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/prettyjson"
	"github.com/qredo/assets/libs/protobuffer"
)

var (
	res = make(map[string]interface{})
)

func getEnv(name, defaultValue string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return v
}

//Use - helper to remove warnings
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func hex2base64(h string) string {
	bin, _ := hex.DecodeString(h)
	return base64.StdEncoding.EncodeToString(bin)
}

func addResultBinaryItem(key string, item []byte) {
	h := hex.EncodeToString(item)
	res[key] = base64.StdEncoding.EncodeToString([]byte(h))
}

func addResultTextItem(key string, item string) {
	res[key] = base64.StdEncoding.EncodeToString([]byte(item))
}

func addResultItem(key string, item interface{}) {
	res[key] = item
}

func ppResult() {
	pp, _ := prettyjson.Marshal(res)
	fmt.Println(string(pp))
}

func addResultSignedAsset(key string, signedAsset *protobuffer.PBSignedAsset) {
	original := reflect.ValueOf(signedAsset)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)
	res[key] = copy.Interface()
}

func prettyStringFromSignedAsset(signedAsset *protobuffer.PBSignedAsset) string {
	original := reflect.ValueOf(signedAsset)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)

	pp, _ := prettyjson.Marshal(copy.Interface())
	return string(pp)
}

//PadRight - right pad a string
func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func TranslateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		TranslateRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		TranslateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			TranslateRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:

		if original.Type() == reflect.TypeOf([]byte("")) {

			b := original.Bytes()
			h := hex.EncodeToString(b)
			hb := []byte(h)
			copy.Set(reflect.ValueOf(hb))

		} else {
			copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
			for i := 0; i < original.Len(); i += 1 {
				TranslateRecursive(copy.Index(i), original.Index(i))
			}
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			TranslateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		b64 := base64.StdEncoding.EncodeToString([]byte(original.String()))
		copy.Set(reflect.ValueOf(b64))
		//translatedString := dict[original.Interface().(string)]
		//		copy.SetString(original.String())

	// And everything else will simply be taken from the original
	default:

		copy.Set(original)
	}

}

func buildKV(data *[]KV) map[string][]byte {
	res := make(map[string][]byte)
	for _, v := range *data {
		key := v.Key
		value := []byte(v.Value)
		res[key] = value
	}
	return res
}

func (cliTool *CLITool) walletFromWalletUpdateJSON(signedUpdate *WalletUpdatePayload) (*assets.Wallet, error) {
	//Decode the JSON

	//Get the New Owner IDDoc
	idNewOwnerKey, err := hex.DecodeString(signedUpdate.Newowner)
	if err != nil {
		return nil, err
	}

	newOwnerIDDoc, err := assets.LoadIDDoc(cliTool.NodeConn, idNewOwnerKey)
	if err != nil {
		return nil, err
	}

	//Get the Existing Wallet
	existingWalletKey, err := hex.DecodeString(signedUpdate.ExistingWalletAssetID)
	if err != nil {
		return nil, err
	}

	originalWallet, err := assets.LoadWallet(cliTool.NodeConn, existingWalletKey)
	if err != nil {
		return nil, err
	}
	originalWallet.DataStore = cliTool.NodeConn
	//Make New Wallet based on Existing
	updatedWallet, err := assets.NewUpdateWallet(originalWallet, newOwnerIDDoc)
	if err != nil {
		return nil, err
	}

	var truths []string
	for _, trans := range signedUpdate.Transfer {

		binParticipants := map[string][]byte{}
		for _, v := range trans.Participants {
			binVal, err := hex.DecodeString(v.ID)
			if err != nil {
				return nil, err
			}
			binParticipants[v.Name] = binVal
		}
		transferType := protobuffer.PBTransferType(trans.TransferType)
		updatedWallet.AddTransfer(transferType, trans.Expression, &binParticipants, trans.Description)
		truthTable, err := updatedWallet.TruthTable(transferType)
		if err != nil {
			return nil, err
		}

		for _, v := range truthTable {
			x := fmt.Sprintf("%d:%s", trans.TransferType, v)
			truths = append(truths, base64.StdEncoding.EncodeToString([]byte(x)))
		}
	}

	//Add in the WalletTransfers - ie. payment destinations
	for _, wt := range signedUpdate.WalletTransfers {
		to, err := hex.DecodeString(wt.To)
		if err != nil {
			return nil, err
		}
		assetID, err := hex.DecodeString(wt.Assetid)
		if err != nil {
			return nil, err
		}

		updatedWallet.AddWalletTransfer(to, wt.Amount, assetID)
	}

	updatedWallet.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType(signedUpdate.TransferType)
	updatedWallet.DataStore = cliTool.NodeConn

	return updatedWallet, nil
}

func (cliTool *CLITool) kVFromKVUpdateJSON(signedUpdate *KVUpdatePayload) (*assets.KVAsset, error) {
	//Decode the JSON

	//Get the New Owner IDDoc
	idNewOwnerKey, err := hex.DecodeString(signedUpdate.Newowner)
	if err != nil {
		return nil, err
	}

	newOwnerIDDoc, err := assets.LoadIDDoc(cliTool.NodeConn, idNewOwnerKey)
	if err != nil {
		return nil, err
	}

	//Get the Existing KV
	existingKVKey, err := hex.DecodeString(signedUpdate.ExistingKVAssetID)
	if err != nil {
		return nil, err
	}

	originalKV, err := assets.LoadKVAsset(cliTool.NodeConn, existingKVKey)
	if err != nil {
		return nil, err
	}

	originalKV.DataStore = cliTool.NodeConn

	//Make New KV based on Existing
	updatedKV, err := assets.NewUpdateKVAsset(originalKV, newOwnerIDDoc)
	if err != nil {
		return nil, err
	}

	//add keys
	for _, pair := range signedUpdate.KV {
		key := pair.Key
		value := pair.Value
		updatedKV.SetKV(key, []byte(value))
	}

	var truths []string
	for _, trans := range signedUpdate.Transfer {

		binParticipants := map[string][]byte{}
		for _, v := range trans.Participants {
			binVal, err := hex.DecodeString(v.ID)
			if err != nil {
				return nil, err
			}
			binParticipants[v.Name] = binVal
		}
		transferType := protobuffer.PBTransferType(trans.TransferType)
		updatedKV.AddTransfer(transferType, trans.Expression, &binParticipants, trans.Description)
		truthTable, err := updatedKV.TruthTable(transferType)
		if err != nil {
			return nil, err
		}

		for _, v := range truthTable {
			x := fmt.Sprintf("%d:%s", trans.TransferType, v)
			truths = append(truths, base64.StdEncoding.EncodeToString([]byte(x)))
		}
	}

	updatedKV.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType(signedUpdate.TransferType)

	return updatedKV, nil
}

func (cliTool *CLITool) groupFromGroupUpdateJSON(signedUpdate *GroupUpdatePayload) (*assets.Group, error) {
	//Decode the JSON

	//Get the New Owner IDDoc
	idNewOwnerKey, err := hex.DecodeString(signedUpdate.Newowner)
	if err != nil {
		return nil, err
	}

	newOwnerIDDoc, err := assets.LoadIDDoc(cliTool.NodeConn, idNewOwnerKey)
	if err != nil {
		return nil, err
	}

	//Get the Existing Group
	existingGroupKey, err := hex.DecodeString(signedUpdate.ExistingGroupAssetID)
	if err != nil {
		return nil, err
	}

	originalGroup, err := assets.LoadGroup(cliTool.NodeConn, existingGroupKey)
	if err != nil {
		return nil, err
	}

	originalGroup.DataStore = cliTool.NodeConn

	//Make New Group based on Existing
	updatedGroup, err := assets.NewUpdateGroup(originalGroup, newOwnerIDDoc)
	if err != nil {
		return nil, err
	}

	//add keys
	payload := updatedGroup.Payload()
	payload.Type = protobuffer.PBGroupType(signedUpdate.Group.Type)
	payload.Description = signedUpdate.Group.Description

	payload.GroupFields = buildKV(&signedUpdate.Group.GroupFields)
	payload.Participants = buildKV(&signedUpdate.Group.Participants)

	var truths []string
	for _, trans := range signedUpdate.Transfer {

		binParticipants := map[string][]byte{}
		for _, v := range trans.Participants {
			binVal, err := hex.DecodeString(v.ID)
			if err != nil {
				return nil, err
			}
			binParticipants[v.Name] = binVal
		}
		transferType := protobuffer.PBTransferType(trans.TransferType)
		updatedGroup.AddTransfer(transferType, trans.Expression, &binParticipants, trans.Description)
		truthTable, err := updatedGroup.TruthTable(transferType)
		if err != nil {
			return nil, err
		}

		for _, v := range truthTable {
			x := fmt.Sprintf("%d:%s", trans.TransferType, v)
			truths = append(truths, base64.StdEncoding.EncodeToString([]byte(x)))
		}
	}

	updatedGroup.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType(signedUpdate.TransferType)

	return updatedGroup, nil
}
