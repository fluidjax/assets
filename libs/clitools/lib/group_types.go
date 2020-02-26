package qc

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
