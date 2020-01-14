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

import "github.com/qredo/assets/libs/protobuffer"

//Core Heirachcy

//Group - Group Asset Type
type Group struct {
	SignedAsset
}

//Wallet - Wallet Asset Type
type Wallet struct {
	SignedAsset
}

//IDDoc - IDDoc Asset Type
type IDDoc struct {
	SignedAsset
}

//SignedAsset - Asset/Previous Asset Wrapper, holding temporary objects (seed) & previousVersions
type SignedAsset struct {
	currentAsset  *protobuffer.PBSignedAsset
	store         *Mapstore                  //Reference to object store (map or blockchain)
	seed          []byte                     //If available a seed to generate keys for object
	previousAsset *protobuffer.PBSignedAsset //Reference to (if any) previous object with the same key
}

//SignatureID - Use to hold ID & Signatures for expression parsing
type SignatureID struct {
	IDDoc        *IDDoc
	Abbreviation string
	HaveSig      bool
	Signature    []byte
}

//TransferParticipant -
type TransferParticipant struct {
	IDDoc        *IDDoc
	Abbreviation string
}
