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

import "github.com/qredo/assets/libs/crypto"

//ConsensusProcess - this is the  Verification for the Consensus Rules.
func (i *IDDoc) ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) error {
	i.DataStore = datasource

	err := i.Verify()
	if err != nil {
		return err
	}
	if deliver == true {
		return i.Deliver(rawTX, txHash)
	}
	return nil
}

func (i *IDDoc) Verify() error {
	assetError := i.SignedAsset.VerifyImmutableCreate()
	if assetError != nil {
		return assetError
	}

	//Check IDDoc specific Fields
	payload, _ := i.Payload()
	if payload.AuthenticationReference == "" {
		return NewAssetsError(CodeConsensusMissingFields, "Consensus:Error:Check:Invalid Madatory Field:AuthenticationReference")
	}
	//check 11
	if payload.BeneficiaryECPublicKey == nil {
		return NewAssetsError(CodeConsensusMissingFields, "Consensus:Error:Check:Invalid Madatory Field:BeneficiaryECPublicKey")
	}
	//check 11
	if payload.SikePublicKey == nil {
		return NewAssetsError(CodeConsensusMissingFields, "Consensus:Error:Check:Invalid Madatory Field:SikePublicKey")
	}
	//check 11
	if payload.BLSPublicKey == nil {
		return NewAssetsError(CodeConsensusMissingFields, "Consensus:Error:Check:Invalid Madatory Field:BLSPublicKey")
	}
	//Check 7
	if i.CurrentAsset.Asset.Index != 1 {
		return NewAssetsError(CodeConsensusIndexNotZero, "Consensus:Error:Check:Invalid Index")
	}

	//Self signed so simply Check the Signature Message
	msg, err := i.SerializeAsset()
	if err != nil {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:Fail to Serialize Asset")
	}
	idDocPayload, err := i.Payload()
	if err != nil {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:Fail to Parse Payload")
	}
	blsPK := idDocPayload.GetBLSPublicKey()
	rc := crypto.BLSVerify(msg, blsPK, i.CurrentAsset.Signature)
	if rc != 0 {
		return NewAssetsError(CodeConsensusSignedAssetFailtoVerify, "Consensus:Error:Check:BLSVerify fails")
	}
	return nil
}

func (i *IDDoc) Deliver(rawTX []byte, txHash []byte) error {
	assetsError := i.AddCoreMappings(rawTX, txHash)
	if assetsError != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:Add Core Mapping TxHash:RawTX")
	}
	return nil
}
