package qredochain

import (
	"github.com/qredo/assets/libs/assets"
	"github.com/tendermint/tendermint/abci/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/kv"
)

func (app *QredoChain) processTX(tx []byte, deliver bool) ([]abcitypes.Event, error) {
	//Decode the Asset

	txAsset, _, txHash, err := assets.BuildAssetFromTX(tx)
	//print("TX :", len(tx), "\n")
	if err != nil {
		return nil, err
	}
	assetError := txAsset.ConsensusProcess(app, tx, txHash, deliver)
	return nil, assetError
}

func processTags(tags map[string][]byte) []types.Event {
	var attributes []kv.Pair
	for key, value := range tags {
		kvpair := kv.Pair{Key: []byte(key), Value: value}
		attributes = append(attributes, kvpair)
	}
	events := []types.Event{
		{
			Type:       "tag",
			Attributes: attributes,
		},
	}
	return events
}
