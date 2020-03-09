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
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

//NewIDDoc create a new IDDoc
func NewMPC() (m *MPC, err error) {
	//generate crypto random seed

	//Main returned Object
	m = &MPC{}

	//Signed Asset
	x := protobuffer.PBSignedAsset{}
	m.CurrentAsset = &x
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_MPC

	//IDDoc
	mpc := &protobuffer.PBMPC{}

	//Compose
	m.CurrentAsset.Asset = asset
	asset.Tags = make(map[string][]byte)
	Payload := &protobuffer.PBAsset_MPC{}
	Payload.MPC = mpc

	m.CurrentAsset.Asset.Payload = Payload

	return m, nil
}

//ReBuildIDDoc rebuild an existing Signed IDDoc into IDDocDeclaration object
//Seed can be manually set if known (ie. Is a local ID)
func ReBuildMPC(sig *protobuffer.PBSignedAsset, key []byte) (m *MPC, err error) {
	if sig == nil {
		return nil, errors.New("ReBuildMPC - sig is nil")
	}
	if key == nil {
		return nil, errors.New("ReBuildMPC  - key is nil")
	}
	m = &MPC{}
	m.CurrentAsset = sig
	m.setKey(key)
	return m, nil
}

//Payload - return the IDDoc payload
func (m *MPC) Payload() (*protobuffer.PBMPC, error) {
	if m == nil {
		return nil, errors.New("MPC is nil")
	}
	if m.CurrentAsset.Asset == nil {
		return nil, errors.New("MPC has no asset")
	}
	return m.CurrentAsset.Asset.GetMPC(), nil
}

func (m *MPC) ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) error {
	assetID := m.Key()
	exists, err := m.Exists(datasource, assetID)
	if err != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Check:Database Access")
	}
	if exists == true {
		//MPC is immutable so if this AssetID already has a value we can't update it.
		return NewAssetsError(CodeCantUpdateImmutableAsset, "Consensus:Error:Check:Immutable Asset")
	}

	payload, err := m.Payload()
	if err != nil {
		return NewAssetsError(CodePayloadEncodingError, "Consensus:Error:Check:Invalid Payload Encoding")
	}

	address := payload.Address
	walletAssetID := payload.AssetID

	if deliver == true {
		//Commit
		err = m.SetWithSuffix(datasource, assetID, ".as2ad", address)
		if err != nil {
			return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Check:Add mappings as2ad & ad2as")
		}

		err = m.SetWithSuffix(datasource, address, ".ad2as", walletAssetID)
		if err != nil {
			return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Check:Add mappings as2ad & ad2as")
		}

	}
	return nil
}
