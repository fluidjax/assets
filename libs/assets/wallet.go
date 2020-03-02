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
	"bytes"

	"math"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

//NewWallet - Setup a new Wallet
func NewWallet(iddoc *IDDoc, currency string) (w *Wallet, err error) {
	if iddoc == nil {
		return nil, errors.New("Sign - supplied IDDoc is nil")
	}
	w = emptyWallet()
	w.DataStore = iddoc.DataStore

	walletKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.CurrentAsset.Asset.ID = walletKey
	w.CurrentAsset.Asset.Type = protobuffer.PBAssetType_Wallet
	w.CurrentAsset.Asset.Owner = iddoc.Key()
	w.CurrentAsset.Asset.Index = 1
	w.CurrentAsset.Asset.TransferType = protobuffer.PBTransferType_None
	w.AssetKeyFromPayloadHash()

	currentPayload, err := w.Payload()
	if err != nil {
		return nil, err
	}
	currentPayload.Currency = currency

	return w, nil
}

//NewUpdateWallet - Create a NewWallet for updates/transfers based on a previous one
func NewUpdateWallet(previousWallet *Wallet, iddoc *IDDoc) (w *Wallet, err error) {
	w = emptyWallet()
	if previousWallet.DataStore != nil {
		w.DataStore = previousWallet.DataStore
	}
	w.CurrentAsset.Asset.ID = previousWallet.CurrentAsset.Asset.ID
	w.CurrentAsset.Asset.Type = protobuffer.PBAssetType_Wallet
	w.CurrentAsset.Asset.Owner = iddoc.Key() //new owner
	w.CurrentAsset.Asset.Index = previousWallet.CurrentAsset.Asset.Index + 1
	w.DataStore = previousWallet.DataStore
	w.PreviousAsset = previousWallet.CurrentAsset
	previousPayload, err := w.PreviousPayload()
	if err != nil {
		return nil, err
	}
	currentPayload, err := w.Payload()
	if err != nil {
		return nil, err
	}
	currentPayload.SpentBalance = previousPayload.SpentBalance
	w.DeepCopyUpdatePayload()
	currentPayload.WalletTransfers = nil
	return w, nil
}

//ReBuildWallet an existing Wallet from it's on chain PBSignedAsset
func ReBuildWallet(sig *protobuffer.PBSignedAsset, key []byte) (w *Wallet, err error) {
	if sig == nil {
		return nil, errors.New("ReBuildIDDoc  - sig is nil")
	}
	if key == nil {
		return nil, errors.New("ReBuildIDDoc  - key is nil")
	}
	w = &Wallet{}
	w.CurrentAsset = sig
	w.setKey(key)
	return w, nil
}

//Payload - return the wallet Payload object
func (w *Wallet) Payload() (*protobuffer.PBWallet, error) {
	if w == nil {
		return nil, errors.New("Wallet is nil")
	}
	if w.CurrentAsset.Asset == nil {
		return nil, errors.New("Wallet has no asset")
	}
	signatureAsset := w.CurrentAsset.Asset
	wallet := signatureAsset.GetWallet()
	return wallet, nil
}

//Payload - return the wallet Previous Payload object
func (w *Wallet) PreviousPayload() (*protobuffer.PBWallet, error) {
	if w == nil {
		return nil, errors.New("Wallet is nil")
	}
	if w.CurrentAsset.Asset == nil {
		return nil, errors.New("Wallet has no asset")
	}
	signatureAsset := w.PreviousAsset.Asset
	wallet := signatureAsset.GetWallet()
	return wallet, nil
}

func (w *Wallet) AddWalletTransfer(to []byte, amount int64, assetid []byte) (err error) {
	if to == nil {
		return errors.New("Transfer to is nil")
	}

	if assetid == nil {
		return errors.New("Transfer assetid is nil")
	}
	if amount == 0 {
		return errors.New("Can't transfer zero amount")
	}

	if amount < 0 {
		return errors.New("Can't transfer negative amount")
	}
	if amount >= math.MaxInt64 {
		return errors.New("Invalid Amount")
	}
	currentPayload, err := w.Payload()

	if err != nil {
		return errors.Wrap(err, "Fail to retrieve Payload in AddWalletTransfer")
	}

	currentPayload.SpentBalance += amount

	currentPayload.WalletTransfers = append(currentPayload.WalletTransfers,
		&protobuffer.PBWalletTransfer{To: to, Amount: amount, AssetID: assetid})
	return nil
}

func (w *Wallet) FullVerify() (bool, error) {
	payload, err := w.Payload()
	if err != nil {
		return false, errors.Wrap(err, "Fail to retrieve Payload in FullVerify")
	}
	previousPayload, err := w.PreviousPayload()

	incomingSpend := previousPayload.SpentBalance
	if incomingSpend < 0 {
		return false, errors.New("Spend less than Zero")
	}
	finalSpend := payload.SpentBalance
	if finalSpend < 0 {
		return false, errors.New("Spend less than Zero")
	}
	spent := finalSpend - incomingSpend

	var calculatedSpent int64
	for _, wt := range payload.WalletTransfers {
		if wt.Amount < 0 {
			return false, errors.New("Spend less than Zero")
		}
		calculatedSpent += wt.Amount

	}

	if calculatedSpent != spent {
		return false, errors.New("Spend Invalid : Previous != Current + Transfers")
	}

	return w.SignedAsset.FullVerify()
}

func emptyWallet() (w *Wallet) {
	w = &Wallet{}
	w.CurrentAsset = &protobuffer.PBSignedAsset{}
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_Wallet
	//Wallet
	wallet := &protobuffer.PBWallet{}
	//Compose
	w.CurrentAsset.Asset = asset
	payload := &protobuffer.PBAsset_Wallet{}
	payload.Wallet = wallet
	w.CurrentAsset.Asset.Payload = payload
	return w
}

//LoadWallet -
func LoadWallet(store DataSource, walletID []byte) (w *Wallet, err error) {
	data, err := store.RawGet(walletID)
	if err != nil {
		return nil, err
	}
	sa := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(data, sa)
	if err != nil {
		return nil, err
	}
	wallet, err := ReBuildWallet(sa, walletID)
	if err != nil {
		return nil, err
	}

	return wallet, nil

}

func (w *Wallet) ConsensusProcess(datasource DataSource, rawTX []byte, txHash []byte, deliver bool) TransactionCode {
	assetID := w.Key()

	exists, err := w.Exists(datasource, assetID)
	if err != nil {
		return CodeDatabaseFail
	}

	// if exists == true {
	// 	println("EXISTS-", hex.EncodeToString(assetID))
	// } else {
	// 	println("NOTEXISTS-", hex.EncodeToString(assetID))
	// }

	//Wallet is mutable, if exists allow update

	if exists == false {
		//This is a new Wallet
		//println("New Wallet", hex.EncodeToString(w.Key()))
		if deliver == true {
			//Commit
			code := w.AddCoreMappings(datasource, rawTX, txHash)
			if code != 0 {
				return CodeDatabaseFail
			}
			//Set Balance to 0
			w.setBalanceKey(datasource, w.Key(), 0)
		}
	} else {
		//This is a wallet update
		//println("Wallet Update")

		//events = processTags(wallet.CurrentAsset.Asset.Tags)

		//Loop through all the transfers out and update their destinations
		payload, err := w.Payload()
		if err != nil {
			return CodeFailVerfication
		}

		currentBalance, code := w.getBalanceKey(datasource, assetID)
		if code != 0 {
			return code
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
			return CodeInsufficientFunds
		}

		//We have enough funds, do the database updates for transfer Pass 2
		if deliver == true {
			//println("Deliver")
			code := w.AddCoreMappings(datasource, rawTX, txHash)
			if code != 0 {
				return CodeDatabaseFail
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
				w.addToBalanceKey(datasource, destinationAssetID, amount)
				totalToSubtract += amount
			}
			w.subtractFromBalanceKey(datasource, assetID, totalToSubtract)
		}
	}
	return CodeTypeOK
}
