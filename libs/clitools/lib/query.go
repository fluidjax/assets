package qc

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/clitools/lib/prettyjson"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
)

func (cliTool *CLITool) GetIDDocForSeed(seedHex string) (iddoc *assets.IDDoc, err error) {
	seed, err := hex.DecodeString(seedHex)
	if err != nil {
		return nil, err
	}
	assetID := assets.KeyFromSeed(seed)
	assedIDHex := hex.EncodeToString(assetID)

	signedAsset, err := cliTool.GetAsset(assedIDHex)

	if err != nil {
		return nil, err
	}

	iddoc, err = assets.ReBuildIDDoc(signedAsset, assetID)

	if err != nil {
		return nil, err
	}

	iddoc.Seed = seed
	return iddoc, nil
}

func (cliTool *CLITool) GetAsset(assetID string) (*protobuffer.PBSignedAsset, error) {
	//Get TX for Asset ID
	txid, err := cliTool.ConsensusSearch(assetID)
	if err != nil {
		return nil, err
	}
	query := "tx.hash='" + hex.EncodeToString(txid) + "'"
	result, err := cliTool.QredoChainSearch(query)
	if len(result) != 1 {
		return nil, errors.New("Incorrect number of responses")
	}

	return result[0], nil
}

func (cliTool *CLITool) PPConsensusSearch(query string) (err error) {
	data, err := cliTool.ConsensusSearch(query)
	if err != nil {
		return err
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

func (cliTool *CLITool) ConsensusSearch(query string) (data []byte, err error) {

	tmClient := cliTool.NodeConn.TmClient
	key, err := hex.DecodeString(query)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to decode Base64 Query %s", query)
	}

	result, err := tmClient.ABCIQuery("V", key)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to run Consensus query %s", query)
	}

	data = result.Response.GetValue()
	return data, nil

}

func (cliTool *CLITool) PPQredoChainSearch(query string) (err error) {
	results, err := cliTool.QredoChainSearch(query)

	original := reflect.ValueOf(results)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)

	pp, _ := prettyjson.Marshal(copy.Interface())
	fmt.Println(string(pp))
	return err
}

func (cliTool *CLITool) QredoChainSearch(query string) (results []*protobuffer.PBSignedAsset, err error) {
	processedCount := 0
	currentPage := 0
	numPerPage := 30

	tmClient := cliTool.NodeConn.TmClient

	for {
		result, err := tmClient.TxSearch(query, false, currentPage, numPerPage)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to search to query %s %d %d", query, currentPage, numPerPage)
		}
		totalToProcess := result.TotalCount
		fmt.Println("Records Found:", totalToProcess)
		if totalToProcess == 0 {
			return nil, nil
		}

		for _, chainTx := range result.Txs {
			processedCount++
			tx := chainTx.Tx
			signedAsset := &protobuffer.PBSignedAsset{}
			err := proto.Unmarshal(tx, signedAsset)

			if err != nil {
				fmt.Println("Error unmarshalling payload")
				if cliTool.checkQuit(processedCount, totalToProcess) == true {
					return results, nil
				}
				continue
			}
			results = append(results, signedAsset)

			if cliTool.checkQuit(processedCount, totalToProcess) == true {
				return results, nil
			}
		}
		currentPage++
	}
}

func (cliTool *CLITool) checkQuit(processedCount int, totalToProcess int) bool {
	return processedCount == totalToProcess
}
