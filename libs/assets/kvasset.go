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
func (w *KVAsset) Payload() *protobuffer.PBKVAsset {
	if w == nil {
		return nil
	}
	signatureAsset := w.CurrentAsset.Asset
	KVAsset := signatureAsset.GetKVAsset()
	return KVAsset
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
	w.AssetKeyFromPayloadHash()
	return w, nil
}

//NewUpdateKVAsset - Create a NewKVAsset for updates/transfers based on a previous one
func NewUpdateKVAsset(previousKVAsset *KVAsset, iddoc *IDDoc) (w *KVAsset, err error) {
	if iddoc == nil {
		return nil, errors.New("NewUpdateKVAsset - supplied IDDoc is nil")
	}
	if previousKVAsset == nil {
		return nil, errors.New("NewUpdateKVAsset - supplied previousKVAsset is nil")
	}
	p := previousKVAsset.CurrentAsset.Asset.GetKVAsset()
	previousType := p.GetType()

	w = emptyKVAsset(previousType)
	if previousKVAsset.DataStore != nil {
		w.DataStore = previousKVAsset.DataStore
	}
	w.CurrentAsset.Asset.ID = previousKVAsset.CurrentAsset.Asset.ID
	w.CurrentAsset.Asset.Type = protobuffer.PBAssetType_KVAsset
	w.CurrentAsset.Asset.Owner = iddoc.Key() //new owner
	w.PreviousAsset = previousKVAsset.CurrentAsset

	//Deep copy the old Payload to the new one
	w.DeepCopyUpdatePayload()

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
		currentSet := w.Payload().AssetFields

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
func ReBuildKVAsset(sig *protobuffer.PBSignedAsset, key []byte) (w *KVAsset, err error) {
	if sig == nil {
		return nil, errors.New("ReBuildIDDoc  - sig is nil")
	}
	if key == nil {
		return nil, errors.New("ReBuildIDDoc  - key is nil")
	}

	w = &KVAsset{}
	w.CurrentAsset = sig
	w.setKey(key)
	return w, nil
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
