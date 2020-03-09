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

package bitcoin

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
)

var (
	TestnetHost = "34.247.28.183:18332"
	TestnetUser = "qredo"
	TestnetPass = "KcsUHi4Hn89ELZJo66vygGtGA"
)

type UnderlyingConnector struct {
	*rpcclient.Client
	*qredochain.NodeConnector
}

func NewUnderlyingConnector(host, user, pass string) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	return client, nil

}

func (conn *UnderlyingConnector) AddAddressToWallet(pubKey string) error {
	return conn.ImportPubKey(pubKey)
}

//ProcessRecentTransactions
func (conn *UnderlyingConnector) ProcessRecentTransactions(blockhash *chainhash.Hash, confirmationThreshold int) (nextBlockHash *chainhash.Hash, count int64, err error) {
	res, err := conn.ListSinceBlock(blockhash, 6, true)
	if err != nil {
		return blockhash, 0, err
	}
	nextBlockHash, err = chainhash.NewHashFromStr(res.LastBlock)

	for _, tx := range res.Transactions {
		if tx.Confirmations < int64(confirmationThreshold) {
			fmt.Println(tx.TxID, " not ready ", tx.Confirmations, " confirmations")
		} else {
			count++
			//fmt.Println(tx.TxID, " * READY *", tx.Confirmations, " confirmations")
			amount := int64(tx.Amount)
			TxID := []byte(tx.TxID)
			address := []byte(tx.Address)
			txid, err := conn.BroadcastUnderlyingChainDeposit(TxID, address, protobuffer.PBCryptoCurrency_BTC, amount)
			if err != nil {
				return blockhash, 0, err
			}
			fmt.Printf("Underlying ADD: Address: %v  TXID: %v \n", address, txid)
		}
	}
	return nextBlockHash, count, nil
}

func (conn *UnderlyingConnector) BroadcastUnderlyingChainDeposit(TxID []byte, address []byte, currency protobuffer.PBCryptoCurrency, amount int64) (txid string, err error) {
	underlying, err := assets.NewUnderlying()
	if err != nil {
		return txid, err
	}

	payload, err := underlying.Payload()
	if err != nil {
		return txid, err
	}

	payload.Type = protobuffer.PBUnderlyingType_Deposit
	payload.CryptoCurrencyCode = currency
	payload.Proof = nil
	payload.Amount = amount
	payload.Address = address
	payload.TxID = TxID
	underlying.AddTag("address", []byte(address))
	txid, err = conn.NodeConnector.PostTx(underlying)
	return txid, err
}
