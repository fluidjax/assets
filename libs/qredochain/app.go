package qredochain

import (
	"github.com/dgraph-io/badger"
	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
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

func decodeTX(data []byte) (*protobuffer.PBSignedAsset, error) {
	signedAsset := &protobuffer.PBSignedAsset{}

	err := proto.Unmarshal(data, signedAsset)
	if err != nil {
		return nil, err
	}
	return signedAsset, nil
}

//DeliverTx -
func (app *KVStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := app.processTX(req.Tx, false)
	return types.ResponseDeliverTx{Code: code, Events: nil}
}

//Commit -
func (app *KVStoreApplication) Commit() abcitypes.ResponseCommit {
	// Persist the application state.
	// Return an (optional) Merkle root hash of the application state
	// ResponseCommit.Data is included as the Header.AppHash in the next block
	// it may be empty
	// Later calls to Query can return proofs about the application state anchored in this Merkle root hash
	// Note developers can return whatever they want here (could be nothing, or a constant string, etc.), so
	//		 long as it is deterministic - it must not be a function of anything that did not come from the
	//		 BeginBlock/DeliverTx/EndBlock methods.

	app.currentBatch.Commit()
	return abcitypes.ResponseCommit{Data: []byte{}}
}

//Query -
func (app *KVStoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	// Query for data from the application at current or past height.
	// Optionally return Merkle proof.
	// Merkle proof includes self-describing type field to support many types of Merkle trees and encoding formats.

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
	// 	Signals the beginning of a new block. Called prior to any DeliverTxs.
	// The header contains the height, timestamp, and more - it exactly matches the Tendermint block header.
	//			 We may seek to generalize this in the future.
	// The LastCommitInfo and ByzantineValidators can be used to determine rewards and punishments for the validators.
	//			 NOTE validators here do not include pubkeys.
	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

//EndBlock -
func (KVStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

//CheckTx -
func (app *KVStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	// Technically optional - not involved in processing blocks.
	// Guardian of the mempool: every node runs CheckTx before letting a transaction into its local mempool.
	// The transaction may come from an external user or another node
	// CheckTx need not execute the transaction in full, but rather a light-weight yet stateful validation,
	//					like checking signatures and account balances, but not running code in a virtual machine.
	// Transactions where ResponseCheckTx.Code != 0 will be rejected - they will not be broadcast to other nodes or included in a proposal block.
	// Tendermint attributes no other value to the response code
	code := app.processTX(req.Tx, true)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 0}
}
