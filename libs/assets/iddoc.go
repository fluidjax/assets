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
	"crypto/sha256"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/keystore"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/store"
)

//NewIDDocWithSeed - generate new IDDoc with supplied seed & Auth ref
func NewIDDocWithSeed(seed []byte, authenticationReference string) (i *IDDoc, err error) {
	if seed == nil {
		return nil, errors.Wrap(err, "No seed provided")
	}
	sikePublicKey, _, err := keystore.GenerateSIKEKeys(seed)
	if err != nil {
		return nil, err
	}

	blsPublicKey, _, err := keystore.GenerateBLSKeys(seed)
	if err != nil {
		return nil, err
	}

	ecPublicKey, err := keystore.GenerateECPublicKey(seed)
	if err != nil {
		return nil, err
	}

	//Main returned Object
	i = &IDDoc{}
	i.Seed = seed

	//Signed Asset
	x := protobuffer.PBSignedAsset{}
	i.CurrentAsset = &x
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_Iddoc

	//IDDoc
	iddoc := &protobuffer.PBIDDoc{}
	iddoc.AuthenticationReference = authenticationReference
	iddoc.BeneficiaryECPublicKey = ecPublicKey
	iddoc.SikePublicKey = sikePublicKey
	iddoc.BLSPublicKey = blsPublicKey

	//Compose
	i.CurrentAsset.Asset = asset
	asset.Tags = make(map[string][]byte)
	Payload := &protobuffer.PBAsset_Iddoc{}
	Payload.Iddoc = iddoc

	i.CurrentAsset.Asset.Payload = Payload
	i.setKey(KeyFromSeed(seed))
	return i, nil
}

//KeyFromSeed double the hash the seed to make a unique key
func KeyFromSeed(seed []byte) []byte {
	onehash := sha256.Sum256(seed)
	key := sha256.Sum256(onehash[:])
	return key[:]

}

//NewIDDoc create a new IDDoc with a random seed
func NewIDDoc(authenticationReference string) (i *IDDoc, err error) {
	//generate crypto random seed
	seed, err := cryptowallet.RandomBytes(48)
	if err != nil {
		err = errors.Wrap(err, "Failed to generate random seed")
		return nil, err
	}

	return NewIDDocWithSeed(seed, authenticationReference)

}

//ReBuildIDDoc rebuild an existing Signed IDDoc into IDDocDeclaration object
//Seed can be manually set if known (ie. Is a local ID)
func ReBuildIDDoc(sig *protobuffer.PBSignedAsset, key []byte) (i *IDDoc, err error) {
	if sig == nil {
		return nil, errors.New("ReBuildIDDoc  - sig is nil")
	}
	if key == nil {
		return nil, errors.New("ReBuildIDDoc  - key is nil")
	}
	i = &IDDoc{}
	i.CurrentAsset = sig
	i.setKey(key)
	return i, nil
}

//Payload - return the IDDoc payload
func (i *IDDoc) Payload() (*protobuffer.PBIDDoc, error) {
	if i == nil {
		return nil, errors.New("IDDoc is nil")
	}
	if i.CurrentAsset.Asset == nil {
		return nil, errors.New("IDDoc has no asset")
	}
	return i.CurrentAsset.Asset.GetIddoc(), nil
}

//LoadIDDoc -
func LoadIDDoc(store store.StoreInterface, iddocID []byte) (i *IDDoc, err error) {
	data, err := store.Load(iddocID)
	if err != nil {
		return nil, err
	}
	sa := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(data, sa)
	if err != nil {
		return nil, err
	}
	iddoc, err := ReBuildIDDoc(sa, iddocID)
	if err != nil {
		return nil, err
	}

	return iddoc, nil

}

func (i *IDDoc) ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) uint32 {
	assetID := i.Key()
	exists, err := i.Exists(datasource, assetID)
	if err != nil {
		return CodeDatabaseFail
	}
	if exists == true {
		//IDDoc is immutable so if this AssetID already has a value we can't update it.
		return CodeAlreadyExists
	}

	//Check the IDDoc is valid
	if i.VerifyIDDoc() == false {
		return CodeFailVerfication
	}

	//Add pointer from AssetID to the txHash of the Object
	if deliver == true {

		code := i.AddCoreMappings(datasource, rawTX, txHash)
		if code != 0 {
			return CodeDatabaseFail
		}
		//events = processTags(iddoc.CurrentAsset.Asset.Tags)
	}
	return CodeTypeOK
}

func (i *IDDoc) VerifyIDDoc() bool {
	//Check signature
	verify, err := i.Verify(i)
	if err != nil {
		return false
	}
	if verify == false {
		return false
	}

	//Check Payload fields
	payload, err := i.Payload()
	if err != nil {
		return false
	}
	if payload == nil {
		return false
	}

	if payload.AuthenticationReference == "" ||
		payload.BeneficiaryECPublicKey == nil ||
		payload.SikePublicKey == nil ||
		payload.BLSPublicKey == nil {
		return false
	}

	if i.CurrentAsset.Asset.Index != 0 {
		return false
	}

	return true
}
