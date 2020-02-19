package qc

// createwalletjson.go
type CreateWalletJSON struct {
	TransferType int64      `json:"transferType"`
	Ownerseed    string     `json:"ownerseed"`
	Currency     string     `json:"currency"`
	Transfer     []Transfer `json:"Transfer"`
}

// transfer.go

type Transfer struct {
	TransferType int64         `json:"TransferType"`
	Expression   string        `json:"Expression"`
	Description  string        `json:"description"`
	Participants []Participant `json:"participants"`
}

// participant.go

type Participant struct {
	Name string `json:"name"`
	ID   string `json:"ID"`
}


type WalletUpdate struct {
	ExistingWalletAssetID string     `json:"ExistingWalletAssetID"`
	NewOwner              string     `json:"newownerseed"`
	TransferType          int64      `json:"transferType"`
	Currency              string     `json:"currency"`
	Transfer              []Transfer `json:"Transfer"`
}

type SignJSON struct {
	Seed string `json:"seed"`
	Msg  string `json:"msg"`
}


type AggregateSignJSON struct {
	Sigs         []Sig        `json:"Sigs"`
	WalletUpdate WalletUpdate `json:"walletUpdate"`
}

// sig.go

type Sig struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	Signature    string `json:"signature"`
}