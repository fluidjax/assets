package qc

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
