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

type TransactionCode uint32

const (
	CodeTypeOK                  TransactionCode = 0
	CodeTypeEncodingError                       = 1
	CodeTypeBadNonce                            = 2
	CodeTypeUnauthorized                        = 3
	CodeAlreadyExists                           = 4
	CodeDatabaseFail                            = 5
	CodeFailVerfication                         = 6
	CodeTypeHTTPError                           = 7
	CodeTendermintInternalError                 = -32603
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
	NodeID     string
	TmNodeAddr string
	HttpClient *http.Client
	TmClient   *tmclient.HTTP
	Log        *logger.Logger
	//Store      *datastore.Store
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
		TmNodeAddr: tmNodeAddr,
		NodeID:     nodeID,
		Log:        log,
		//	Store:      store,
		HttpClient: &http.Client{
			Timeout: nodeConnectionTimeout,
		},
		TmClient: tmClient,
	}, nil
}

func (nc *NodeConnector) TxSearch(query string, prove bool, currentPage int, numPerPage int) (resultRaw *ctypes.ResultTxSearch, err error) {
	resultRaw, err = nc.TmClient.TxSearch(query, prove, currentPage, numPerPage)
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
	url := "http://" + nc.TmNodeAddr

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", CodeTypeHTTPError, errors.Wrap(err, "post to blockchain node")
	}
	req.Header.Set("Content-Type", "text/plain;")

	resp, err := nc.HttpClient.Do(req)
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

	if data["result"] == nil {
		if data["error"] != nil {
			errdata := data["error"].(map[string]interface{})
			codef64 := errdata["code"].(float64)
			code = TransactionCode(codef64)
			return "", code, errors.New("Failed to add new TX")
		} else {
			return "", CodeAlreadyExists, errors.New("Failed to decode response")
		}
	}

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
	result, err := nc.TmClient.TxSearch(query, true, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(result.Txs) == 0 {
		return nil, errors.Errorf("Document not found: %v", txHash)
	}

	return result.Txs[0].Tx, nil
}

// func (nc *NodeConnector) subscribeAndQueue(ctx context.Context, txQueue chan *protobuffer.TX) error {
// 	query := "tag.recipient='" + nc.nodeID + "'"

// 	out, err := nc.tmClient.Subscribe(context.Background(), "", query, 1000)
// 	if err != nil {
// 		return errors.Wrapf(err, "Failed to subscribe to query %s", query)
// 	}

// 	go func() {
// 		for {
// 			select {
// 			case result := <-out:
// 				tx := result.Data.(tmtypes.EventDataTx).Tx
// 				incomingTX := &protobuffer.TX{}
// 				err := proto.Unmarshal(tx, incomingTX)
// 				incomingTX.Height = result.Data.(tmtypes.EventDataTx).Height
// 				incomingTX.Index = result.Data.(tmtypes.EventDataTx).Index

// 				if err != nil {
// 					nc.log.Debug("IGNORED TX - Invalid!")
// 					break
// 				}

// 				//check if this node is in receipient list
// 				if incomingTX.RecipientCID != nc.nodeID {
// 					nc.log.Debug("IGNORED TX! Recipient not match the query! (%v != %v)", incomingTX.RecipientCID, nc.nodeID)
// 					break
// 				}

// 				//Add into the waitingQueue for later processing
// 				txQueue <- incomingTX
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}()

// 	return nil
// }
