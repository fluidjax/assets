package qredochain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/datastore"
	"github.com/qredo/assets/libs/logger"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	nodeConnectionTimeout = time.Second * 10
	txChanSize            = 1000
)

type ChainPostable interface {
	SerializeSignedAsset() ([]byte, error)
	Key() []byte
}

type NodeConnector struct {
	nodeID     string
	tmNodeAddr string
	httpClient *http.Client
	tmClient   *tmclient.HTTP
	log        *logger.Logger
	store      *datastore.Store
}

// NewNodeConnector constructs a new Tendermint NodeConnector
func NewNodeConnector(tmNodeAddr string, nodeID string, store *datastore.Store, log *logger.Logger) (conn *NodeConnector, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("Initialize tendermint node connector: %v", r)
		}
	}()
	tmNodeAddr = strings.TrimRight(tmNodeAddr, "/")
	tmClient, err := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", tmNodeAddr), "/websocket")
	if err := tmClient.Start(); err != nil {
		return nil, errors.Wrap(err, "Start tendermint client")
	}
	return &NodeConnector{
		tmNodeAddr: tmNodeAddr,
		nodeID:     nodeID,
		log:        log,
		store:      store,
		httpClient: &http.Client{
			Timeout: nodeConnectionTimeout,
		},
		tmClient: tmClient,
	}, nil
}

func (nc *NodeConnector) TxSearch(query string, prove bool, currentPage int, numPerPage int) (resultRaw *ctypes.ResultTxSearch, err error) {
	resultRaw, err = nc.tmClient.TxSearch(query, prove, currentPage, numPerPage)
	return resultRaw, err
}

// PostTx posts a transaction to the chain and returns the transaction ID
func (nc *NodeConnector) PostTx(asset ChainPostable) (txID string, code TransactionCode, err error) {
	// //serialize the whole transaction
	serializedTX, err := asset.SerializeSignedAsset()
	if err != nil {
		return "", CodeTypeEncodingError, err
	}
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX)

	// // TODO: use net/rpc

	//broadcast_tx_commit - broadcast and wait until its in a new block
	//broadcast_tx_async  - broadcast and return - no checks
	//broadcast_tx_sync   - broadcast and wait for CheckTx result

	body := strings.NewReader(`{
		"jsonrpc": "2.0",
		"id": "anything",
		"method": "broadcast_tx_commit",
		"params": {
			"tx": "` + base64EncodedTX + `"}
	}`)
	url := "http://" + nc.tmNodeAddr

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", CodeTypeHTTPError, errors.Wrap(err, "post to blockchain node")
	}
	req.Header.Set("Content-Type", "text/plain;")

	resp, err := nc.httpClient.Do(req)
	if err != nil {
		return "", CodeTypeHTTPError, errors.Wrap(err, "post to blockchain node")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var f interface{}
	err2 := json.Unmarshal(b, &f)
	if err2 != nil {
		return
	}

	data := f.(map[string]interface{})
	result := data["result"].(map[string]interface{})
	txID = result["hash"].(string)
	checktx := result["check_tx"].(map[string]interface{})
	codef64 := checktx["code"].(float64)
	code = TransactionCode(codef64)
	return txID, code, err
}

func (nc *NodeConnector) GetTx(txHash string) ([]byte, error) {
	//query := fmt.Sprintf("tag.txid='%s'", txHash)
	query := "tag.myname='chris'"
	print("QUERY:", query, "\n")
	result, err := nc.tmClient.TxSearch(query, true, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(result.Txs) == 0 {
		return nil, errors.Errorf("Document not found: %v", txHash)
	}

	return result.Txs[0].Tx, nil
}
