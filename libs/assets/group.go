package assets

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

//GroupPayload - return the Group Payload object
func (w *Group) Payload() *protobuffer.PBGroup {
	if w == nil {
		return nil
	}
	signatureAsset := w.PBSignedAsset.Asset
	Group := signatureAsset.GetGroup()
	return Group
}

//NewGroup - Setup a new IDDoc
func NewGroup(iddoc *IDDoc, groupType protobuffer.PBGroupType) (w *Group, err error) {
	if iddoc == nil {
		return nil, errors.New("NewGroup - supplied IDDoc is nil")
	}
	w = emptyGroup(groupType)
	w.store = iddoc.store

	GroupKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.PBSignedAsset.Asset.ID = GroupKey
	w.PBSignedAsset.Asset.Type = protobuffer.PBAssetType_Group
	w.PBSignedAsset.Asset.Owner = iddoc.Key()
	w.assetKeyFromPayloadHash()

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

	p := previousGroup.PBSignedAsset.Asset.GetGroup()
	previousType := p.GetType()

	w = emptyGroup(previousType)
	if previousGroup.store != nil {
		w.store = previousGroup.store
	}
	w.PBSignedAsset.Asset.ID = previousGroup.PBSignedAsset.Asset.ID
	w.PBSignedAsset.Asset.Type = protobuffer.PBAssetType_Group
	w.PBSignedAsset.Asset.Owner = iddoc.Key() //new owner
	w.previousAsset = &previousGroup.PBSignedAsset
	return w, nil
}

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
	w.PBSignedAsset.Asset.Payload = payload

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
	deleyed := {AAAA}
	unchanged := {BBBB,CCCC}

	Allows a GUI to easily display to the end user what changes they need to agree too.
*/
func (w *Group) ParseChanges() (unchanged, added, deleted [][]byte, err error) {
	if w == nil {
		return nil, nil, nil, errors.New("ParseChanges - No Asset")
	}
	if w.previousAsset == nil {
		return nil, nil, nil, errors.New("ParseChanges - No previous Asset")
	}

	switch g := w.previousAsset.GetAsset().GetPayload().(type) {
	case *protobuffer.PBAsset_Group:
		//group
		previousSet := g.Group.Participants
		currentSet := w.Payload().Participants

		//Unchanged & Deleted
		for _, pId := range previousSet {
			found := false
			for _, cId := range currentSet {
				res := bytes.Compare(pId, cId)
				if res == 0 {
					found = true
				}
			}
			if found == false {
				deleted = append(deleted, pId)
			} else {
				unchanged = append(unchanged, pId)
			}
		}

		//Added
		//Unchanged & Deleted
		for _, cId := range currentSet {
			found := false
			for _, pId := range previousSet {
				res := bytes.Compare(pId, cId)
				if res == 0 {
					found = true
				}
			}
			if found == false {
				added = append(added, cId)
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
	w.PBSignedAsset = *sig
	w.setKey(key)
	return w, nil
}

func emptyGroup(groupType protobuffer.PBGroupType) (w *Group) {
	w = &Group{}
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_Group
	//Group
	group := &protobuffer.PBGroup{}
	group.Type = groupType
	//Compose
	w.PBSignedAsset.Asset = asset
	payload := &protobuffer.PBAsset_Group{}
	payload.Group = group
	w.PBSignedAsset.Asset.Payload = payload
	return w
}
