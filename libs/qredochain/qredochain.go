package qredochain

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/kv"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

func (app *Qredochain) SetLogger(l log.Logger) {
	app.logger = l
}

func (app *Qredochain) Info(req types.RequestInfo) types.ResponseInfo {
	res := app.app.Info(req)
	res.LastBlockHeight = app.app.state.Height
	res.LastBlockAppHash = app.app.state.AppHash
	return res
}

func (app *Qredochain) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return app.app.SetOption(req)
}

// tx is either "val:pubkey!power" or "key=value" or just arbitrary bytes
func (app *Qredochain) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	// if it starts with "val:", update the validator set
	// format is "val:pubkey!power"
	if isValidatorTx(req.Tx) {
		// update validators in the merkle tree
		// and in app.ValUpdates
		return app.execValidatorTx(req.Tx)
	}

	// otherwise, update the key-value store
	return app.app.DeliverTx(req)
}

func (app *Qredochain) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	return app.app.CheckTx(req)
}

// Commit will panic if InitChain was not called
func (app *Qredochain) Commit() types.ResponseCommit {
	return app.app.Commit()
}

// When path=/val and data={validator address}, returns the validator update (types.ValidatorUpdate) varint encoded.
// For any other path, returns an associated value or nil if missing.
func (app *Qredochain) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	switch reqQuery.Path {
	case "/val":
		key := []byte("val:" + string(reqQuery.Data))
		value, err := app.app.state.db.Get(key)
		if err != nil {
			panic(err)
		}

		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		return
	default:
		return app.app.Query(reqQuery)
	}
}

// Save the validators in the merkle tree
func (app *Qredochain) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	for _, v := range req.Validators {
		r := app.updateValidator(v)
		if r.IsErr() {
			app.logger.Error("Error updating validators", "r", r)
		}
	}
	return types.ResponseInitChain{}
}

// Track the block hash and header information
func (app *Qredochain) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	// reset valset changes
	app.ValUpdates = make([]types.ValidatorUpdate, 0)

	for _, ev := range req.ByzantineValidators {
		if ev.Type == tmtypes.ABCIEvidenceTypeDuplicateVote {
			// decrease voting power by 1
			if ev.TotalVotingPower == 0 {
				continue
			}
			app.updateValidator(types.ValidatorUpdate{
				PubKey: app.valAddrToPubKeyMap[string(ev.Validator.Address)],
				Power:  ev.TotalVotingPower - 1,
			})
		}
	}
	return types.ResponseBeginBlock{}
}

// Update the validator set
func (app *Qredochain) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

//---------------------------------------------
// update validators

func (app *Qredochain) Validators() (validators []types.ValidatorUpdate) {
	itr, err := app.app.state.db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	for ; itr.Valid(); itr.Next() {
		if isValidatorTx(itr.Key()) {
			validator := new(types.ValidatorUpdate)
			err := types.ReadMessage(bytes.NewBuffer(itr.Value()), validator)
			if err != nil {
				panic(err)
			}
			validators = append(validators, *validator)
		}
	}
	return
}

func (app *Application) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       ProtocolVersion.Uint64(),
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.AppHash,
	}
}

// tx is either "key=value" or just arbitrary bytes
func (app *Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	var key, value []byte
	parts := bytes.Split(req.Tx, []byte("="))
	if len(parts) == 2 {
		key, value = parts[0], parts[1]
	} else {
		key, value = req.Tx, req.Tx
	}

	app.state.db.Set(prefixKey(key), value)
	app.state.Size++

	events := []types.Event{
		{
			Type: "app",
			Attributes: []kv.Pair{
				{Key: []byte("creator"), Value: []byte("Cosmoshi Netowoko")},
				{Key: []byte("key"), Value: key},
			},
		},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: events}
}

func (app *Application) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (app *Application) Commit() types.ResponseCommit {
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height++
	saveState(app.state)
	return types.ResponseCommit{Data: appHash}
}

// Returns an associated value or nil if missing.
func (app *Application) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	if reqQuery.Prove {
		value, err := app.state.db.Get(prefixKey(reqQuery.Data))
		if err != nil {
			panic(err)
		}
		if value == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value

		return
	}

	resQuery.Key = reqQuery.Data
	value, err := app.state.db.Get(prefixKey(reqQuery.Data))
	if err != nil {
		panic(err)
	}
	if value == nil {
		resQuery.Log = "does not exist"
	} else {
		resQuery.Log = "exists"
	}
	resQuery.Value = value

	return resQuery
}
