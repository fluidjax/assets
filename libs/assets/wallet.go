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
	"math"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

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

func (w *Wallet) AddWalletTransfer(to []byte, amount int64) (err error) {
	if to == nil {
		return errors.New("Transfer to is nil")
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
		&protobuffer.PBWalletTransfer{To: to, Amount: amount})
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

//NewWallet - Setup a new IDDoc
func NewWallet(iddoc *IDDoc, currency string) (w *Wallet, err error) {
	if iddoc == nil {
		return nil, errors.New("Sign - supplied IDDoc is nil")
	}
	w = emptyWallet()
	w.Store = iddoc.Store

	walletKey, err := RandomBytes(32)
	if err != nil {
		return nil, errors.New("Fail to generate random key")
	}
	w.CurrentAsset.Asset.ID = walletKey
	w.CurrentAsset.Asset.Type = protobuffer.PBAssetType_wallet
	w.CurrentAsset.Asset.Owner = iddoc.Key()
	w.assetKeyFromPayloadHash()

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
	if previousWallet.Store != nil {
		w.Store = previousWallet.Store
	}
	w.CurrentAsset.Asset.ID = previousWallet.CurrentAsset.Asset.ID
	w.CurrentAsset.Asset.Type = protobuffer.PBAssetType_wallet
	w.CurrentAsset.Asset.Owner = iddoc.Key() //new owner
	w.CurrentAsset.Asset.Index = previousWallet.CurrentAsset.Asset.Index + 1

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

func emptyWallet() (w *Wallet) {
	w = &Wallet{}
	w.CurrentAsset = &protobuffer.PBSignedAsset{}
	//Asset
	asset := &protobuffer.PBAsset{}
	asset.Type = protobuffer.PBAssetType_wallet
	//Wallet
	wallet := &protobuffer.PBWallet{}
	//Compose
	w.CurrentAsset.Asset = asset
	payload := &protobuffer.PBAsset_Wallet{}
	payload.Wallet = wallet
	w.CurrentAsset.Asset.Payload = payload
	return w
}
