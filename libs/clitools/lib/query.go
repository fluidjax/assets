package qc

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/clitools/lib/prettyjson"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

func GetAsset(qredochain string, assetID string) error {
	//Get TX for Asset ID
	// key, err := hex.DecodeString(assetID)
	// query := fmt.Sprintf("tx.hash='%s'", key)

	// tmClient, _ := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", qredochain), "/websocket")
	// defer tmClient.Stop()

	// if err := tmClient.Start(); err != nil {
	// 	return errors.Wrapf(err, "Failed to start Tendermint client")
	// }
	// key, err = hex.DecodeString(query)
	// if err != nil {
	// 	return errors.Wrapf(err, "Failed to decode Base64 Query %s", query)
	// }

	// result, err := tmClient.ABCIQuery("V", key)

	// if err != nil {
	// 	return errors.Wrapf(err, "Failed to run Consensus query %s", query)
	// }

	// assetID = result.Response.GetValue()

	// //Get TX for Asset ID
	return nil
}

func PPConsensusSearch(qredochain string, query string) (err error) {
	data, err := ConsensusSearch(qredochain, query)
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

func ConsensusSearch(qredochain string, query string) (data []byte, err error) {

	tmClient, _ := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", qredochain), "/websocket")
	defer tmClient.Stop()

	if err := tmClient.Start(); err != nil {
		return nil, errors.Wrapf(err, "Failed to start Tendermint client")
	}

	//key, err := base64.StdEncoding.DecodeString(query)
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

func PPQredoChainSearch(qredochain string, query string) (err error) {
	results, err := QredoChainSearch(qredochain, query)

	original := reflect.ValueOf(results)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)

	pp, _ := prettyjson.Marshal(copy.Interface())
	fmt.Println(string(pp))
	return err
}

func QredoChainSearch(qredochain string, query string) (results []*protobuffer.PBSignedAsset, err error) {
	processedCount := 0
	currentPage := 0
	numPerPage := 30

	tmClient, _ := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", qredochain), "/websocket")
	defer tmClient.Stop()

	if err := tmClient.Start(); err != nil {
		return nil, errors.Wrapf(err, "Failed to start Tendermint client")
	}

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
				if checkQuit(processedCount, totalToProcess) == true {
					return results, nil
				}
				continue
			}
			results = append(results, signedAsset)

			if checkQuit(processedCount, totalToProcess) == true {
				return results, nil
			}
		}
		currentPage++
	}
}

func checkQuit(processedCount int, totalToProcess int) bool {
	return processedCount == totalToProcess
}
