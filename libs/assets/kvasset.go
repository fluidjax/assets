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

/*
	A KVAsset Asset is simple KV store (KVAssetFields) together with a list of
	particpants. Primarily create for the TrusteeKVAsset Transaction type.
*/

import (
	"bytes"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

//SetKV - set the KV Field value
func (w *KVAsset) SetKV(key string, value []byte) error {
	if w.CurrentAsset == nil {
		return errors.New("KVAsset - currentAsset is nil")
	}

	signatureAsset := w.CurrentAsset.Asset
	kvAsset := signatureAsset.GetKVAsset()

	if containsString(kvAsset.Immutable, key) {
		return errors.New("KVAsset - Field is immutable")
	}

	if kvAsset.AssetFields == nil {
		kvAsset.AssetFields = make(map[string][]byte)
	}
	kvAsset.AssetFields[key] = value
	return nil
}

func (w *KVAsset) SetKeyString(key string) error {
	if w.CurrentAsset == nil {
		return errors.New("KVAsset - currentAsset is nil")
	}
	if key == "" {
		return errors.New("KVAsset - key can't be empty")
	}
	w.CurrentAsset.Asset.ID = []byte(key)
	return nil
}

func (w *KVAsset) SetImmutable(fields []string) error {
	if w.CurrentAsset == nil {
		return errors.New("KVAsset - currentAsset is nil")
	}
	if fields == nil {
		return errors.New("KVAsset - Immutable fields is is nil")
	}
	signatureAsset := w.CurrentAsset.Asset
	kvAsset := signatureAsset.GetKVAsset()
	kvAsset.Immutable = fields
	return nil
}

//SetKV - set the KV Field value
func (w *KVAsset) GetKV(key string) ([]byte, error) {
	if w.CurrentAsset == nil {
		return nil, errors.New("KVAsset - currentAsset is nil")
	}
	signatureAsset := w.CurrentAsset.Asset
	kvAsset := signatureAsset.GetKVAsset()

	if kvAsset.AssetFields == nil {
		return nil, errors.New("KVAsset.AssetFields - is nil")
	}
	return kvAsset.AssetFields[key], nil
}

//Payload - return the KVAsset Payload object
func (w *KVAsset) Payload() (*protobuffer.PBKVAsset, error) {
	if w == nil {
		return nil, errors.New("KV is nil")
	}
	if w.CurrentAsset.Asset == nil {
		return nil, errors.New("KV has no asset")
	}

	signatureAsset := w.CurrentAsset.Asset
	KVAsset := signatureAsset.GetKVAsset()
	return KVAsset, nil
}

//NewKVAsset - Setup a new KVAsset
func NewKVAsset(iddoc *IDDoc, assetType protobuffer.PBKVAssetType) (w *KVAsset, err error) {
	if iddoc == nil {
		return nil, errors.New("NewKVAsset - supplied IDDoc is nil")
	}
	w = emptyKVAsset(assetType)
	w.DataStore = iddoc.DataStore
	KVAssetKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.CurrentAsset.Asset.ID = KVAssetKey
	w.CurrentAsset.Asset.Type = protobuffer.PBAssetType_KVAsset
	w.CurrentAsset.Asset.Owner = iddoc.Key()
	w.CurrentAsset.Asset.Index = 1
	w.AssetKeyFromPayloadHash()
	return w, nil
}

//ParseChanges given as KVAsset Asset which has a previous version
//Determine how the asset ID have changed.
//The abbreviations are ignore, and a list of unchanged, added & delete IDs are returned.
/*
	old[abbreviation:ID] {
		abbrev1 : AAAA,
		abbrev2 : BBBB,
		abbrev3 : CCCC,
	}

	new[abbreviation:ID] {
		abbrev1 : BBBB
		abbrev2 : CCCC,
		abbrev3 : DDDD,
		abbrev4 : EEEE,
	}

	added := {DDDD, EEEE}
	deleted := {AAAA}
	unchanged := {BBBB,CCCC}

	Allows a GUI to easily display to the end user what changes they need to agree too.
*/
func (w *KVAsset) ParseChanges() (unchanged, added, deleted [][]byte, err error) {
	if w == nil {
		return nil, nil, nil, errors.New("ParseChanges - No Asset")
	}
	if w.PreviousAsset == nil {
		return nil, nil, nil, errors.New("ParseChanges - No previous Asset")
	}

	switch g := w.PreviousAsset.GetAsset().GetPayload().(type) {
	case *protobuffer.PBAsset_KVAsset:
		//KVAsset
		previousSet := g.KVAsset.AssetFields
		payload, err := w.Payload()
		if err != nil {
			return nil, nil, nil, errors.New("Error retrieving Payload")
		}

		currentSet := payload.AssetFields

		//Unchanged & Deleted
		for _, pID := range previousSet {
			found := false
			for _, cID := range currentSet {
				res := bytes.Compare(pID, cID)
				if res == 0 {
					found = true
				}
			}
			if found == false {
				deleted = append(deleted, pID)
			} else {
				unchanged = append(unchanged, pID)
			}
		}

		//Added
		for _, cID := range currentSet {
			found := false
			for _, pID := range previousSet {
				res := bytes.Compare(pID, cID)
				if res == 0 {
					found = true
				}
			}
			if found == false {
				added = append(added, cID)
			}
		}
	default:
		//invalid type
		return nil, nil, nil, errors.New("Previous Asset is invalid")
	}
	return unchanged, added, deleted, nil
}

//ReBuildKVAsset an existing KVAsset from it's on chain PBSignedAsset
func ReBuildKVAsset(sig *protobuffer.PBSignedAsset, key []byte) (kv *KVAsset, err error) {
	if sig == nil {
		return nil, errors.New("ReBuildIDDoc  - sig is nil")
	}
	if key == nil {
		return nil, errors.New("ReBuildIDDoc  - key is nil")
	}

	kv = &KVAsset{}
	kv.CurrentAsset = sig
	kv.setKey(key)
	return kv, nil
}

func emptyKVAsset(KVAssetType protobuffer.PBKVAssetType) (w *KVAsset) {
	w = &KVAsset{}
	w.CurrentAsset = &protobuffer.PBSignedAsset{}
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_KVAsset
	//KVAsset
	KVAsset := &protobuffer.PBKVAsset{}
	KVAsset.Type = KVAssetType
	//Compose
	w.CurrentAsset.Asset = asset
	payload := &protobuffer.PBAsset_KVAsset{}
	payload.KVAsset = KVAsset
	w.CurrentAsset.Asset.Payload = payload
	return w
}

//LoadKVAsset -
func LoadKVAsset(store DataSource, kvassetID []byte) (kv *KVAsset, err error) {
	data, err := store.RawGet(kvassetID)
	if err != nil {
		return nil, err
	}
	sa := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(data, sa)
	if err != nil {
		return nil, err
	}
	kv, err = ReBuildKVAsset(sa, kvassetID)
	if err != nil {
		return nil, err
	}

	return kv, nil

}

//NewUpdateWallet - Create a NewWallet for updates/transfers based on a previous one
func NewUpdateKVAsset(previousKVAsset *KVAsset, iddoc *IDDoc) (kv *KVAsset, err error) {

	previousType := previousKVAsset.CurrentAsset.GetAsset().Type

	kv = emptyKVAsset(protobuffer.PBKVAssetType(previousType))
	if previousKVAsset.DataStore != nil {
		kv.DataStore = previousKVAsset.DataStore
	}
	kv.CurrentAsset.Asset.ID = previousKVAsset.CurrentAsset.Asset.ID
	kv.CurrentAsset.Asset.Type = protobuffer.PBAssetType_KVAsset
	kv.CurrentAsset.Asset.Owner = iddoc.Key() //new owner
	kv.CurrentAsset.Asset.Index = previousKVAsset.CurrentAsset.Asset.Index + 1

	kv.PreviousAsset = previousKVAsset.CurrentAsset
	kv.DataStore = previousKVAsset.DataStore
	kv.DeepCopyUpdatePayload()

	return kv, nil
}

//Payload - return the wallet Previous Payload object
func (k *KVAsset) PreviousPayload() (*protobuffer.PBKVAsset, error) {
	if k == nil {
		return nil, errors.New("KVAsset is nil")
	}
	if k.CurrentAsset.Asset == nil {
		return nil, errors.New("KVAsset has no asset")
	}
	signatureAsset := k.PreviousAsset.Asset
	kv := signatureAsset.GetKVAsset()
	return kv, nil
}

func (k *KVAsset) ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) error {
	assetID := k.Key()
	exists, err := k.Exists(datasource, assetID)
	if err != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Check:Database Access")
	}
	//Wallet is mutable, if exists allow update

	if exists == false {
		//This is a new Wallet
		if deliver == true {
			//Commit
			assetsError := k.AddCoreMappings(datasource, rawTX, txHash)
			if assetsError != nil {
				return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:Add Core Mapping TxHash:RawTX")
			}
		}

	} else {
		if deliver == true {
			//Commit
			assetsError := k.AddCoreMappings(datasource, rawTX, txHash)
			if assetsError != nil {
				return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:Add Core Mapping TxHash:RawTX")
			}
		}

	}
	return nil
}
