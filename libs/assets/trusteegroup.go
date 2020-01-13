package assets

import (
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

//TrusteeGroupPayload - return the trusteeGroup Payload object
func (w *TrusteeGroup) Payload() *protobuffer.PBTrusteeGroup {
	signatureAsset := w.PBSignedAsset.Asset
	trusteeGroup := signatureAsset.GetTrusteeGroup()
	return trusteeGroup
}

//NewTrusteeGroup - Setup a new IDDoc
func NewTrusteeGroup(iddoc *IDDoc) (w *TrusteeGroup, err error) {
	w = emptyTrusteeGroup()
	w.store = iddoc.store

	trusteeGroupKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.PBSignedAsset.Asset.ID = trusteeGroupKey
	w.PBSignedAsset.Asset.Type = protobuffer.PBAssetType_trusteeGroup
	w.PBSignedAsset.Asset.Owner = iddoc.Key()
	w.assetKeyFromPayloadHash()

	return w, nil
}

//NewUpdateTrusteeGroup - Create a NewTrusteeGroup for updates/transfers based on a previous one
func NewUpdateTrusteeGroup(previousTrusteeGroup *TrusteeGroup, iddoc *IDDoc) (w *TrusteeGroup, err error) {
	w = emptyTrusteeGroup()
	if previousTrusteeGroup.store != nil {
		w.store = previousTrusteeGroup.store
	}
	w.PBSignedAsset.Asset.ID = previousTrusteeGroup.PBSignedAsset.Asset.ID
	w.PBSignedAsset.Asset.Type = protobuffer.PBAssetType_trusteeGroup
	w.PBSignedAsset.Asset.Owner = iddoc.Key() //new owner
	w.previousAsset = &previousTrusteeGroup.PBSignedAsset
	return w, nil
}

func (w *TrusteeGroup) ConfigureTrusteeGroup(expression string, participants *map[string][]byte, description string) error {
	transferRule := &protobuffer.PBTransfer{}
	transferRule.Expression = expression
	transferRule.Description = description
	transferRule.Type = protobuffer.PBTransferType_TrusteeGroupDefinition
	if transferRule.Participants == nil {
		transferRule.Participants = make(map[string][]byte)
	}
	for abbreviation, iddocID := range *participants {
		transferRule.Participants[abbreviation] = iddocID
	}
	//Cant use enum as map key, so convert to a string

	pbtrusteeGroup := &protobuffer.PBTrusteeGroup{}
	pbtrusteeGroup.TrusteeGroup = transferRule

	payload := &protobuffer.PBAsset_TrusteeGroup{}
	payload.TrusteeGroup = pbtrusteeGroup
	w.PBSignedAsset.Asset.Payload = payload

	return nil

}

//ReBuildTrusteeGroup an existing TrusteeGroup from it's on chain PBSignedAsset
func ReBuildTrusteeGroup(sig *protobuffer.PBSignedAsset) (w *TrusteeGroup, err error) {
	w = &TrusteeGroup{}
	w.PBSignedAsset = *sig
	return w, nil
}

func emptyTrusteeGroup() (w *TrusteeGroup) {
	w = &TrusteeGroup{}
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_trusteeGroup
	//TrusteeGroup
	trusteeGroup := &protobuffer.PBTrusteeGroup{}
	//Compose
	w.PBSignedAsset.Asset = asset
	payload := &protobuffer.PBAsset_TrusteeGroup{}
	payload.TrusteeGroup = trusteeGroup
	w.PBSignedAsset.Asset.Payload = payload
	return w
}
