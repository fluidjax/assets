package lib

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
	truths            []string
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

type CreateGroupJSON struct {
	TransferType int64      `json:"transferType"`
	Ownerseed    string     `json:"ownerseed"`
	Group        Group      `json:"group"`
	Transfer     []Transfer `json:"Transfer"`
}

type GroupUpdate struct {
	Sigs               []Sig              `json:"sigs"`
	GroupUpdatePayload GroupUpdatePayload `json:"GroupUpdatePayload"`
}

type GroupUpdatePayload struct {
	ExistingGroupAssetID string     `json:"existingGroupAssetID"`
	Newowner             string     `json:"newowner"`
	TransferType         int64      `json:"transferType"`
	Group                Group      `json:"group"`
	Transfer             []Transfer `json:"transfer"`
}

type Group struct {
	Type         int64  `json:"type"`
	Description  string `json:"description"`
	GroupFields  []KV   `json:"groupfields"`
	Participants []KV   `json:"Participants"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KVUpdate struct {
	Sigs            []Sig           `json:"sigs"`
	KVUpdatePayload KVUpdatePayload `json:"kvUpdatePayload"`
}

type KVUpdatePayload struct {
	ExistingKVAssetID string     `json:"existingWalletAssetID"`
	Newowner          string     `json:"newowner"`
	TransferType      int64      `json:"transferType"`
	KV                []KV       `json:"kv"`
	Transfer          []Transfer `json:"transfer"`
}

type CreateKVJSON struct {
	TransferType int64  `json:"transferType"`
	KVAssetType  int64  `json:"kvAssetType"`
	Ownerseed    string `json:"ownerseed"`
	AssetID      string `json:"assetID"`

	KV       []KV       `json:"kv"`
	Transfer []Transfer `json:"Transfer"`
}

type WalletUpdate struct {
	Sigs                []Sig               `json:"sigs"`
	WalletUpdatePayload WalletUpdatePayload `json:"walletUpdatePayload"`
}

type Sig struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	Signature    string `json:"signature"`
}

type WalletUpdatePayload struct {
	ExistingWalletAssetID string           `json:"existingWalletAssetID"`
	Newowner              string           `json:"newowner"`
	TransferType          int64            `json:"transferType"`
	Currency              string           `json:"currency"`
	WalletTransfers       []WalletTransfer `json:"walletTransfers"`
	Transfer              []Transfer       `json:"transfer"`
}

// transfer.go

type Transfer struct {
	TransferType int64         `json:"transferType"`
	Expression   string        `json:"expression"`
	Description  string        `json:"description"`
	Participants []Participant `json:"participants"`
}

// participant.go

type Participant struct {
	Name string `json:"name"`
	ID   string `json:"ID"`
}

// wallettransfer.go

type WalletTransfer struct {
	To      string `json:"to"`
	Amount  int64  `json:"amount"`
	Assetid string `json:"assetid"`
}
