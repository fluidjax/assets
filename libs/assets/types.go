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
	"github.com/qredo/assets/libs/protobuffer"
)

//Core Heirachcy

type MPC struct {
	SignedAsset
}

type Underlying struct {
	SignedAsset
}

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

//KVAsset - A Generic KV asset for storing key value paramaters.
type KVAsset struct {
	SignedAsset
}

//SignedAsset - Asset/Previous Asset Wrapper, holding temporary objects (seed) & previousVersions
type SignedAsset struct {
	CurrentAsset *protobuffer.PBSignedAsset
	//Store         *DataSource //Reference to object store (map or blockchain)
	DataStore     DataSource
	Seed          []byte                     //If available a seed to generate keys for object
	PreviousAsset *protobuffer.PBSignedAsset //Reference to (if any) previous object with the same key
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

//TXAsset - generic wrapper for All Transactions
//Implements the core Blockchain funcationality
type TXAsset interface {
	ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool, asset TXAsset) error
	Verify() error                             //Check if Transaction is valid
	Deliver(rawTX []byte, txHash []byte) error //Finalize the Transaction - write fields to DB etc
}

type DataSource interface {
	BatchGet(key []byte) ([]byte, error)         //Get data from the current Batch
	Set(key []byte, data []byte) (string, error) //Set data in the current Batch
	RawGet(key []byte) ([]byte, error)           //Get data from underlying commited  database
	GetAssetbyID(assetID []byte) ([]byte, error) //Get Asset using Asset ID

	CacheSet(key []byte, data []byte) //Put the Item in the cache
	CacheGet(key []byte) []byte       //Get the item from the cache

}
