package qc

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/prettyjson"
)

func (cliTool *CLITool) GetIDDocForSeed(seedHex string) (iddoc *assets.IDDoc, err error) {
	seed, err := hex.DecodeString(seedHex)
	if err != nil {
		return nil, err
	}
	assetID := assets.KeyFromSeed(seed)
	assedIDHex := hex.EncodeToString(assetID)

	signedAsset, err := cliTool.NodeConn.GetAsset(assedIDHex)

	if err != nil {
		return nil, err
	}

	iddoc, err = assets.ReBuildIDDoc(signedAsset, assetID)

	if err != nil {
		return nil, err
	}
	iddoc.DataStore = cliTool.NodeConn
	iddoc.Seed = seed
	return iddoc, nil
}

func (cliTool *CLITool) PPConsensusSearch(query string, suffix string) (err error) {

	data, err := cliTool.NodeConn.ConsensusSearch(query, suffix)
	if err != nil {
		return err
	}

	if suffix != "" {
		addResultItem("suffix", suffix)
	}
	addResultItem("key", query)
	addResultItem("value", hex.EncodeToString(data))

	original := reflect.ValueOf(res)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)

	pp, _ := prettyjson.Marshal(copy.Interface())
	fmt.Println(string(pp))

	return nil
}

func (cliTool *CLITool) PPQredoChainSearch(query string) (err error) {
	results, err := cliTool.NodeConn.QredoChainSearch(query)

	original := reflect.ValueOf(results)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)

	pp, _ := prettyjson.Marshal(copy.Interface())
	fmt.Println(string(pp))
	return err
}
