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
func NewUnderlying() (u *Underlying, err error) {
	//generate crypto random seed

	//Main returned Object
	u = &Underlying{}

	//Signed Asset
	x := protobuffer.PBSignedAsset{}
	u.CurrentAsset = &x
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_Underlying

	underlying := &protobuffer.PBUnderlying{}

	//Compose
	u.CurrentAsset.Asset = asset
	asset.Tags = make(map[string][]byte)
	Payload := &protobuffer.PBAsset_Underlying{}
	Payload.Underlying = underlying

	u.CurrentAsset.Asset.Payload = Payload
	return u, nil
}

//ReBuildUnderlying -
func ReBuildUnderlying(sig *protobuffer.PBSignedAsset, key []byte) (u *Underlying, err error) {
	if sig == nil {
		return nil, errors.New("ReBuildUnderlying  - sig is nil")
	}
	if key == nil {
		return nil, errors.New("ReBuildUnderlying  - key is nil")
	}
	u = &Underlying{}
	u.CurrentAsset = sig
	u.setKey(key)
	return u, nil
}

//Payload - return the IDDoc payload
func (u *Underlying) Payload() (*protobuffer.PBUnderlying, error) {
	if u == nil {
		return nil, errors.New("Underlying is nil")
	}
	if u.CurrentAsset.Asset == nil {
		return nil, errors.New("Underlying has no asset")
	}
	return u.CurrentAsset.Asset.GetUnderlying(), nil
}
