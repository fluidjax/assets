package qredochain

import (
	"crypto/sha256"

	"github.com/dgraph-io/badger"
	"github.com/qredo/assets/libs/assets"
	"github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

//https://k0d.su/dkodnik/minter-go-node/src/commit/14439248926def5f6f1d334dc9a84f7d0d2e5a0a/core/minter/minter.go?lang=pl-PL

//QredoChain -
type QredoChain struct {
	DB           *badger.DB
	CurrentBatch *badger.Txn
	Height       uint64
}

var _ abcitypes.Application = (*QredoChain)(nil)

//NewQredoChain -
func NewQredoChain(db *badger.DB) *QredoChain {

	kv := &QredoChain{
		DB:           db,
		CurrentBatch: db.NewTransaction(true),
	}
	return kv
}

//Info -
func (app *QredoChain) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	lh, err := app.GetLastHeight()
	if err != nil {
		return abcitypes.ResponseInfo{}
	}
	lastHeight := int64(lh)

	lastBlockHash, _ := app.GetLastBlockHash()
	if lastHeight == 0 {
		return abcitypes.ResponseInfo{}
	}

	return abcitypes.ResponseInfo{
		LastBlockHeight:  lastHeight,
		LastBlockAppHash: lastBlockHash,
	}
}

//SetOption -
func (app *QredoChain) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

//Commit -
func (app *QredoChain) Commit() abcitypes.ResponseCommit {
	// Persist the application state.
	// Return an (optional) Merkle root hash of the application state
	// ResponseCommit.Data is included as the Header.AppHash in the next block
	// it may be empty
	// Later calls to Query can return proofs about the application state anchored in this Merkle root hash
	// Note developers can return whatever they want here (could be nothing, or a constant string, etc.), so
	//		 long as it is deterministic - it must not be a function of anything that did not come from the
	//		 BeginBlock/DeliverTx/EndBlock methods.

	app.CurrentBatch.Commit()

	hash := sha256.Sum256([]byte("TEST"))

	app.SetLastBlockHash(hash[:])
	app.SetLastHeight(app.Height)

	return abcitypes.ResponseCommit{Data: hash[:]}
}

//Query -
func (app *QredoChain) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	// Query for data from the application at current or past height.
	// Optionally return Merkle proof.
	// Merkle proof includes self-describing type field to support many types of Merkle trees
	//	and encoding formats.

	switch reqQuery.Path {
	case "V": //V = get Value
		err := app.DB.View(func(txn *badger.Txn) error {
			item, err := txn.Get(reqQuery.Data)

			if err != nil && err != badger.ErrKeyNotFound {
				return err
			}
			if err == badger.ErrKeyNotFound {
				resQuery.Log = "does not exist"
			} else {
				return item.Value(func(val []byte) error {
					resQuery.Log = "exists"
					//valHex := hex.EncodeToString(val)
					resQuery.Value = []byte(val)
					return nil
				})
			}
			return nil
		})
		if err != nil {
			resQuery.Code = 1
			return
		}
	case "I": //I = indirect value

		err := app.DB.View(func(txn *badger.Txn) error {
			//item, err := txn.Get(reqQuery.Data)
			txid, err := txn.Get(reqQuery.Data)
			if err != nil && err != badger.ErrKeyNotFound {
				resQuery.Code = 1
				return err
			}

			err2 := txid.Value(func(val []byte) error {
				indirectValue, err3 := txn.Get(val)
				if err3 != nil && err3 != badger.ErrKeyNotFound {
					return err3
				}
				return indirectValue.Value(func(ival []byte) error {
					resQuery.Log = "exists"
					resQuery.Value = []byte(ival)
					return nil
				})

			})
			if err2 != nil {
				resQuery.Code = 1
				return err2
			}
			return nil
		})
		if err != nil {
			resQuery.Code = 1
			return
		}

	}
	return
}

//InitChain -
func (QredoChain) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

//BeginBlock -
func (app *QredoChain) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	// 	Signals the beginning of a new block. Called prior to any DeliverTxs.
	// The header contains the height, timestamp, and more - it exactly matches the Tendermint block header.
	//			 We may seek to generalize this in the future.
	// The LastCommitInfo and ByzantineValidators can be used to determine rewards and punishments for the validators.
	//			 NOTE validators here do not include pubkeys.
	app.Height = uint64(req.Header.Height)
	//app.db.SetLastHeight(height)
	//fmt.Printf("Current block is %d", app.db.GetLastHeight())

	app.CurrentBatch = app.DB.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

//EndBlock -
func (QredoChain) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

//CheckTx -
func (app *QredoChain) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	// Technically optional - not involved in processing blocks.
	// Guardian of the mempool: every node runs CheckTx before letting a transaction into its local mempool.
	// The transaction may come from an external user or another node
	// CheckTx need not execute the transaction in full, but rather a light-weight yet stateful validation,
	//					like checking signatures and account balances, but not running code in a virtual machine.
	// Transactions where ResponseCheckTx.Code != 0 will be rejected - they will not be broadcast to other nodes or included in a proposal block.
	// Tendermint attributes no other value to the response code
	events, err := app.processTX(req.Tx, false)
	var code uint32
	var data []byte

	if err != nil {
		if err, ok := err.(*assets.AssetsError); ok {
			data = []byte(err.Error())
			code = uint32(err.Code())
		}
	}

	return abcitypes.ResponseCheckTx{Code: code, Data: data, Events: events}
}

//DeliverTx -
func (app *QredoChain) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	events, err := app.processTX(req.Tx, true)
	var code uint32
	var data []byte
	if err != nil {
		if err, ok := err.(*assets.AssetsError); ok {
			data = []byte(err.Error())
			code = uint32(err.Code())
		}
	}
	return types.ResponseDeliverTx{Code: code, Data: data, Events: events}
}
