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
			amount := float32(tx.Amount)
			txid, code, err := conn.BroadcastUnderlyingChainDeposit(tx.TxID, tx.Address, protobuffer.PBCryptoCurrency_BTC, amount)
			if code != 0 || err != nil {
				return blockhash, 0, err
			}
			fmt.Printf("Underlying ADD: Address: %v  TXID: %v \n", tx.Address, txid)
		}
	}
	return nextBlockHash, count, nil
}

func (conn *UnderlyingConnector) BroadcastUnderlyingChainDeposit(TxID string, address string, currency protobuffer.PBCryptoCurrency, amount float32) (txid string, code qredochain.TransactionCode, err error) {
	underlying, err := assets.NewUnderlying()
	if err != nil {
		return txid, code, err
	}

	payload, err := underlying.Payload()
	if err != nil {
		return txid, code, err
	}

	payload.Type = protobuffer.PBUnderlyingType_Deposit
	payload.CryptoCurrencyCode = currency
	payload.Proof = nil
	payload.Amount = amount
	payload.Address = address
	payload.TxID = TxID
	underlying.AddTag("address", []byte(address))
	return conn.NodeConnector.PostTx(underlying)
}
