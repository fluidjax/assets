package qc

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

// kv.go

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
