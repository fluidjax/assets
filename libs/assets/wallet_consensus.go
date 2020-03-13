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



//Verify -
func (w *Wallet) Verify() error {
	walletUpdate, err := w.Exists(w.Key())
	if err != nil {
		return NewAssetsError(CodeDatabaseFail, "Fail to access database")
	}
	if walletUpdate == false {
		return w.VerifyCreate()
	}
	return w.VerifyUpdate()
}

//VerifyCreate -
func (w *Wallet) VerifyCreate() (err error) {
	err = w.SignedAsset.VerifyMutableCreate()
	if err != nil {
		return err
	}

	//Check Currency has been set
	payload, err := w.Payload()
	if payload.Currency == 0 {
		return NewAssetsError(CodeConsensusMissingFields, "Consensus:Error:Check:WalletCreate:Invalid Madatory Field:Currency")
	}
	//Check Balance starts at 0
	if payload.SpentBalance != 0 {
		return NewAssetsError(CodeConsensusMissingFields, "Consensus:Error:Check:Invalid Madatory Field:Balance Starts at 0")
	}

	return nil
}

//VerifyUpdate -
func (w *Wallet) VerifyUpdate() (err error) {
	assetID := w.Key()

	err = w.LoadPreviousWallet()
	if err != nil {
		return err
	}

	err = w.SignedAsset.VerifyMutableUpdate()
	if err != nil {
		return err
	}

	payload, err := w.Payload()

	currentBalance, assetsError := w.getBalanceKey(assetID)
	if assetsError != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Check:WalletUpdate:Fail to fetch Current Balance")
	}

	//Check we have enough - Pass 1
	var totalOutgoing int64
	for _, wt := range payload.WalletTransfers {
		res := bytes.Compare(wt.AssetID, assetID)
		if res == 0 {
			//this is money coming back to self, just ignore it
			continue
		}
		totalOutgoing = totalOutgoing + wt.Amount
	}

	if totalOutgoing > currentBalance {
		//println("Eject")
		return NewAssetsError(CodeConsensusInsufficientFunds, "Consensus:Error:Check:WalletUpdate:Outgoing > CurrentBalance")
	}

	return nil
}

//Deliver -
func (w *Wallet) Deliver(rawTX []byte, txHash []byte) error {
	walletUpdate, err := w.Exists(w.Key())
	if err != nil {
		return NewAssetsError(CodeDatabaseFail, "Fail to access database")
	}
	if walletUpdate == false {
		return w.DeliverCreate(rawTX, txHash)
	}
	return w.DeliverUpdate(rawTX, txHash)
}

//DeliverCreate -
func (w *Wallet) DeliverCreate(rawTX []byte, txHash []byte) (err error) {
	//New Wallet Deliver add core & set balance key
	assetsError := w.AddCoreMappings(rawTX, txHash)
	if assetsError != nil {
		return NewAssetsError(CodeDatabaseFail, "Consensus:Error:Deliver:WalletCreate:Fail to Add Core Mappings")
	}
	w.setBalanceKey(w.Key(), 0)
	return nil
}

//DeliverUpdate
func (w *Wallet) DeliverUpdate(rawTX []byte, txHash []byte) (err error) {
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
