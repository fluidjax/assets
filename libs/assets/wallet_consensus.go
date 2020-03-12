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

import "bytes"

//ConsensusProcess - this is the  Verification for the Consensus Rules.
//Different rules/processes depending on whether its an Update or Create and tendermint Check_TX or Deliver_TX
func (w *Wallet) ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) error {
	w.DataStore = datasource

	//Check if Asset Exists - ie, if update or create
	walletUpdate, err := w.Exists(w.Key())
	if err != nil {
		return NewAssetsError(CodeDatabaseFail, "Fail to access database")
	}

	if walletUpdate == false {
		//This is a new Wallet - CREATE
		assetError := w.VerifyWalletCreate()
		if assetError != nil {
			return assetError
		}
		if deliver == true {
			return w.DeliverWalletCreate(rawTX, txHash)
		}
	} else {
		//This is a wallet UPDATE

		assetError := w.VerifyWalletUpdate()
		if assetError != nil {
			return assetError
		}
		if deliver == true {
			w.DeliverWalletUpdate(rawTX, txHash)
		}
	}
	return nil
}






func (w *Wallet) DeliverWalletCreate(rawTX []byte, txHash []byte) (err error) {
	//New Wallet Deliver
	assetsError := w.AddCoreMappings(rawTX, txHash)
	if assetsError != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:WalletCreate:Fail to Add Core Mappings")
	}
	w.setBalanceKey(w.Key(), 0)
	return nil
}

//Deliver - the transaction is being committed/dewlivered so update the Consensus Database with changes
func (w *Wallet) DeliverWalletUpdate(rawTX []byte, txHash []byte) (err error) {
	assetID := w.Key()
	payload, err := w.Payload()
	if err != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:WalletUpdate:Fail to determine Wallet Payload")
	}

	assetsError := w.AddCoreMappings(rawTX, txHash)
	if assetsError != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:WalletUpdate:Fail to Add Core Mappings")
	}
	var totalToSubtract int64

	for _, wt := range payload.WalletTransfers {
		res := bytes.Compare(wt.AssetID, assetID)
		if res == 0 {
			//this is money coming back to self, just ignore it
			continue
		}
		amount := wt.Amount
		destinationAssetID := wt.AssetID
		w.addToBalanceKey(destinationAssetID, amount)
		totalToSubtract += amount
	}
	w.subtractFromBalanceKey(assetID, totalToSubtract)
	return nil
}
