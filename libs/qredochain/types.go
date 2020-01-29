package qredochain

import (
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
)

type TransactionCode uint32

const (
	CodeTypeOK            TransactionCode = 0
	CodeTypeEncodingError                 = 1
	CodeTypeBadNonce                      = 2
	CodeTypeUnauthorized                  = 3
	CodeAlreadyExists                     = 4
	CodeDatabaseFail                      = 5
	CodeFailVerfication                   = 6
	CodeTypeHTTPError                     = 7
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")

	ProtocolVersion version.Protocol  = 0x1
	_               types.Application = (*Qredochain)(nil)
)

const (
	ValidatorSetChangePrefix string = "val:"
)

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

type Qredochain struct {
	types.BaseApplication
	state              State
	ValUpdates         []types.ValidatorUpdate
	valAddrToPubKeyMap map[string]types.PubKey
	logger             log.Logger
}
