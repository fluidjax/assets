package qc

import (
	"sync"

	"github.com/qredo/assets/libs/qredochain"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

//CLITool connector for CLI
type CLITool struct {
	NodeConn      *qredochain.NodeConnector
	ConnectString string
}

var (
	nc                *qredochain.NodeConnector
	out               <-chan ctypes.ResultEvent
	done              = make(chan struct{})
	wg                sync.WaitGroup
	mu                sync.Mutex // protects ctr
	ctr               = 0
	displayTopItem    = 0
	displayBottomItem = 0
	highlightRow      = 0
	txhistory         []ctypes.ResultEvent
)

type SignJSON struct {
	Seed string `json:"seed"`
	Msg  string `json:"msg"`
}

type CreateWalletJSON struct {
	TransferType int64      `json:"transferType"`
	Ownerseed    string     `json:"ownerseed"`
	Currency     string     `json:"currency"`
	Transfer     []Transfer `json:"Transfer"`
}
