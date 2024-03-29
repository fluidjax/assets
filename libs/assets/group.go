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
	A group Asset is simple KV store (GroupFields) together with a list of
	particpants. Primarily create for the TrusteeGroup Transaction type.
*/

import (
	"bytes"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

//Payload - return the Group Payload object
func (w *Group) Payload() *protobuffer.PBGroup {
	if w == nil {
		return nil
	}
	signatureAsset := w.CurrentAsset.Asset
	Group := signatureAsset.GetGroup()
	return Group
}

//NewGroup - Setup a new Group
func NewGroup(iddoc *IDDoc, groupType protobuffer.PBGroupType) (w *Group, err error) {
	if iddoc == nil {
		return nil, errors.New("NewGroup - supplied IDDoc is nil")
	}
	w = EmptyGroup(groupType)
	w.DataStore = iddoc.DataStore
	GroupKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.CurrentAsset.Asset.ID = GroupKey
	w.CurrentAsset.Asset.Type = protobuffer.PBAssetType_Group
	w.CurrentAsset.Asset.Owner = iddoc.Key()
	w.CurrentAsset.Asset.Index = 1
	w.AssetKeyFromPayloadHash()
	return w, nil
}

//NewUpdateGroup - Create a NewGroup for updates/transfers based on a previous one
func NewUpdateGroup(previousGroup *Group, iddoc *IDDoc) (w *Group, err error) {
	if iddoc == nil {
		return nil, errors.New("NewUpdateGroup - supplied IDDoc is nil")
	}
	if previousGroup == nil {
		return nil, errors.New("NewUpdateGroup - supplied previousGroup is nil")
	}
	p := previousGroup.CurrentAsset.Asset.GetGroup()
	previousType := p.GetType()

	w = EmptyGroup(previousType)
	if previousGroup.DataStore != nil {
		w.DataStore = previousGroup.DataStore
	}
	w.CurrentAsset.Asset.ID = previousGroup.CurrentAsset.Asset.ID
	w.CurrentAsset.Asset.Type = protobuffer.PBAssetType_Group
	w.CurrentAsset.Asset.Owner = iddoc.Key() //new owner
	w.CurrentAsset.Asset.Index = previousGroup.CurrentAsset.Asset.Index + 1

	w.PreviousAsset = previousGroup.CurrentAsset
	w.DataStore = previousGroup.DataStore
	w.DeepCopyUpdatePayload()
	return w, nil
}

// ConfigureGroup - configure the Group
func (w *Group) ConfigureGroup(expression string, participants *map[string][]byte, description string) error {
	if w == nil {
		return errors.New("ConfigureGroup - group is nil")
	}
	pbGroup := &protobuffer.PBGroup{}
	if pbGroup.Participants == nil {
		pbGroup.Participants = make(map[string][]byte)
	}
	for abbreviation, iddocID := range *participants {
		pbGroup.Participants[abbreviation] = iddocID
	}
	if pbGroup.GroupFields == nil {
		pbGroup.GroupFields = make(map[string][]byte)
	}
	pbGroup.Description = description
	expressionBytes := []byte(expression)
	pbGroup.GroupFields["expression"] = expressionBytes
	payload := &protobuffer.PBAsset_Group{}
	payload.Group = pbGroup
	w.CurrentAsset.Asset.Payload = payload
	return nil
}

//ParseChanges given as Group Asset which has a previous version
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
func (w *Group) ParseChanges() (unchanged, added, deleted [][]byte, err error) {
	if w == nil {
		return nil, nil, nil, errors.New("ParseChanges - No Asset")
	}
	if w.PreviousAsset == nil {
		return nil, nil, nil, errors.New("ParseChanges - No previous Asset")
	}

	switch g := w.PreviousAsset.GetAsset().GetPayload().(type) {
	case *protobuffer.PBAsset_Group:
		//group
		previousSet := g.Group.Participants
		currentSet := w.Payload().Participants

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

//ReBuildGroup an existing Group from it's on chain PBSignedAsset
func ReBuildGroup(sig *protobuffer.PBSignedAsset, key []byte) (w *Group, err error) {
	if sig == nil {
		return nil, errors.New("ReBuildIDDoc  - sig is nil")
	}
	if key == nil {
		return nil, errors.New("ReBuildIDDoc  - key is nil")
	}

	w = &Group{}
	w.CurrentAsset = sig
	w.setKey(key)
	return w, nil
}

func EmptyGroup(groupType protobuffer.PBGroupType) (w *Group) {
	w = &Group{}
	w.CurrentAsset = &protobuffer.PBSignedAsset{}
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_Group
	//Group
	group := &protobuffer.PBGroup{}
	group.Type = groupType
	//Compose
	w.CurrentAsset.Asset = asset
	payload := &protobuffer.PBAsset_Group{}
	payload.Group = group
	w.CurrentAsset.Asset.Payload = payload
	return w
}

//LoadGroup -
func LoadGroup(store DataSource, groupAssetID []byte) (g *Group, err error) {
	data, err := store.RawGet(groupAssetID)
	if err != nil {
		return nil, err
	}
	sa := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(data, sa)
	if err != nil {
		return nil, err
	}
	g, err = ReBuildGroup(sa, groupAssetID)
	if err != nil {
		return nil, err
	}

	return g, nil

}

// func (g *Group) ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) error {
// 	assetID := g.Key()
// 	g.DataStore = datasource
// 	exists, err := g.Exists(assetID)
// 	if err != nil {
// 		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Check:Database Access")
// 	}
// 	//Wallet is mutable, if exists allow update

// 	if exists == false {
// 		//This is a new Wallet
// 		if deliver == true {
// 			//Commit
// 			err := g.AddCoreMappings(rawTX, txHash)
// 			if err != nil {
// 				return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:Add Core Mapping TxHash:RawTX")
// 			}
// 		} else {
// 			//Check
// 			err := g.VerifyGroup()
// 			if err != nil {
// 				return err
// 			}
// 		}

// 	} else {
// 		if deliver == true {
// 			//Commit
// 			err := g.AddCoreMappings(rawTX, txHash)
// 			if err != nil {
// 				return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:Add Core Mapping TxHash:RawTX")
// 			}
// 		}
// 	}
// 	return nil
// }

func (g *Group) VerifyGroup() error {
	//fy, err := g.Verify()
	return nil
}

func (g *Group) Verify() (err error) {
	return nil
}
func (g *Group) Deliver(rawTX []byte, txHash []byte) (err error) {
	return nil
}
