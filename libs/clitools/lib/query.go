package qc

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/gogo/protobuf/proto"
	"github.com/hokaccha/go-prettyjson"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

func ConsensusSearch(qredochain string, query string) error {

	tmClient, _ := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", qredochain), "/websocket")
	defer tmClient.Stop()

	if err := tmClient.Start(); err != nil {
		return errors.Wrapf(err, "Failed to start Tendermint client")
	}

	key, err := base64.StdEncoding.DecodeString(query)
	if err != nil {
		return errors.Wrapf(err, "Failed to decode Base64 Query %s", query)
	}

	result, err := tmClient.ABCIQuery("V", key)

	if err != nil {
		return errors.Wrapf(err, "Failed to run Consensus query %s", query)
	}

	data := result.Response.GetValue()
	res := make(map[string]string)
	res["type"] = "C"
	res["query"] = query
	res["hex"] = hex.EncodeToString(data)
	res["base64"] = base64.StdEncoding.EncodeToString(data)
	pp, _ := prettyjson.Marshal(res)
	fmt.Println(string(pp))

	return nil

}

func QredoChainSearch(qredochain string, query string) error {
	processedCount := 0
	currentPage := 0
	numPerPage := 30

	tmClient, _ := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", qredochain), "/websocket")
	defer tmClient.Stop()

	if err := tmClient.Start(); err != nil {
		return errors.Wrapf(err, "Failed to start Tendermint client")
	}

	for {
		result, err := tmClient.TxSearch(query, false, currentPage, numPerPage)
		if err != nil {
			return errors.Wrapf(err, "Failed to search to query %s %d %d", query, currentPage, numPerPage)
		}
		totalToProcess := result.TotalCount
		fmt.Println("Records Found:", totalToProcess)
		if totalToProcess == 0 {
			return nil
		}

		for _, chainTx := range result.Txs {
			processedCount++
			tx := chainTx.Tx
			signedAsset := &protobuffer.PBSignedAsset{}
			err := proto.Unmarshal(tx, signedAsset)

			if err != nil {
				fmt.Println("Error unmarshalling payload")
				if checkQuit(processedCount, totalToProcess) == true {
					return nil
				}
				continue
			}

			pp, _ := prettyjson.Marshal(signedAsset)
			fmt.Println(string(pp))
			if checkQuit(processedCount, totalToProcess) == true {
				return nil
			}

		}
		currentPage++
	}

}

func checkQuit(processedCount int, totalToProcess int) bool {
	return processedCount == totalToProcess
}

func getEnv(name, defaultValue string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return v
}

//Use - helper to remove warnings
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
