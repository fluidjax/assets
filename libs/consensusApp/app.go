package main

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/dgraph-io/badger"
	"github.com/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"

	"github.com/qredo/assets/libs/assets"
	"github.com/tendermint/tendermint/libs/kv"
)

//KVStoreApplication -
type KVStoreApplication struct {
	db           *badger.DB
	currentBatch *badger.Txn
}

var _ abcitypes.Application = (*KVStoreApplication)(nil)

//NewKVStoreApplication -
func NewKVStoreApplication(db *badger.DB) *KVStoreApplication {
	return &KVStoreApplication{
		db: db,
	}
}

//Info -
func (KVStoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

//SetOption -
func (KVStoreApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

//DeliverTx -
func (app *KVStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	// payload, err := decodeTX(req.Tx)

	print("hello")
	calcTXHash := sha256.Sum256(req.Tx)
	print("Hash:", hex.EncodeToString(calcTXHash[:]))

	events := []abcitypes.Event{
		{
			Type: "transfer",
			Attributes: kv.Pairs{
				kv.Pair{Key: []byte("sender"), Value: []byte("Chris")},
				kv.Pair{Key: []byte("recipient"), Value: []byte("Alice")},
				kv.Pair{Key: []byte("balance"), Value: []byte("101")},
			},
		},
	}
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: events}

	// if err != nil {
	// 	print("Invalid Transaction - ignore")
	// 	return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError, Events: nil}
	// }

	// calcTXHash := sha256.Sum256(req.Tx)
	// calcTXHashBase58 := base58.Encode(calcTXHash[:])

	// print("Calc hash = ", calcTXHashBase58)

	// var atts []cmn.Pair

	// // //Add the tags
	// // for k, v := range payload.Key() {
	// // 	atts = append(atts, cmn.Pair{Key: []byte(k), Value: []byte(v)})
	// // }

	// // //Add all the recipients
	// // for _, v := range payload.AdditionalRecipientCID {
	// // 	atts = append(atts, cmn.Pair{Key: []byte("recipient"), Value: []byte(v)})
	// // }

	// //Add the TX hash
	// // atts = append(atts, cmn.Pair{Key: []byte("sender"), Value: []byte(payload.SenderCID)})
	// // atts = append(atts, cmn.Pair{Key: []byte("recipient"), Value: []byte(payload.RecipientCID)})
	// // atts = append(atts, cmn.Pair{Key: []byte("txhash"), Value: []byte(TXHash)})
	// // atts = append(atts, cmn.Pair{Key: []byte("txtype"), Value: []byte(strconv.Itoa(int(payload.TXType)))})
	// atts = append(atts, cmn.Pair{Key: []byte("key"), Value: []byte(calcTXHash[:])})
	// atts = append(atts, cmn.Pair{Key: []byte("key58"), Value: []byte(calcTXHashBase58)})
	// atts = append(atts, cmn.Pair{Key: []byte("hello"), Value: []byte("bye")})
	// atts = append(atts, cmn.Pair{Key: []byte("A"), Value: []byte("B")})

	// events := []types.Event{
	// 	{
	// 		Type:       "tag", // curl "localhost:26657/tx_search?query=\"tag.key58='9Hi8MpLNNQiha7eH6bejKs6HhdvKyc9Mt7yMxw4bP5rP'\""
	// 		Attributes: atts,
	// 	},
	// }

	// //fmt.Printf("\n\n****** BLOCK %v %v\n", payload.Processor, payload.RecipientCID)

	// return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: events}
	// //return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: nil}

}

//Commit -
func (app *KVStoreApplication) Commit() abcitypes.ResponseCommit {
	app.currentBatch.Commit()
	return abcitypes.ResponseCommit{Data: []byte{}}
}

//Query -
func (app *KVStoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	// resQuery.Log = "exists"
	// resQuery.Value = []byte("helloo")
	// return

	print("\nXXXX", reqQuery.Data)

	resQuery.Key = reqQuery.Data
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(reqQuery.Data)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == badger.ErrKeyNotFound {
			resQuery.Log = "does not exist"
		} else {
			return item.Value(func(val []byte) error {
				resQuery.Log = "exists"
				resQuery.Value = val
				return nil
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

//InitChain -
func (KVStoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

//BeginBlock -
func (app *KVStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

//EndBlock -
func (KVStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (app *KVStoreApplication) isValid(tx []byte) (code uint32) {
	print("is valid")
	return 0
	_, err := decodeTX(tx)
	if err != nil {
		return 1
	}
	return 0
	// return 0
	// print("check is valid")
	// // check format
	// parts := bytes.Split(tx, []byte("="))
	// if len(parts) != 2 {
	// 	return 1
	// }

	// key, value := parts[0], parts[1]

	// // check if the same key=value already exists
	// err := app.db.View(func(txn *badger.Txn) error {
	// 	item, err := txn.Get(key)
	// 	if err != nil && err != badger.ErrKeyNotFound {
	// 		return err
	// 	}
	// 	if err == nil {
	// 		return item.Value(func(val []byte) error {
	// 			if bytes.Equal(val, value) {
	// 				code = 2
	// 			}
	// 			return nil
	// 		})
	// 	}

	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// return code
}

//CheckTx -
func (app *KVStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	//code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: 0, GasWanted: 0}
	//	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 0}
}

func decodeTX(data []byte) (assets.SignedAsset, error) {
	payload := assets.SignedAsset{}
	return payload, nil
}
