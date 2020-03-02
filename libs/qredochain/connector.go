package qredochain

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/datastore"
	"github.com/qredo/assets/libs/logger"
	"github.com/qredo/assets/libs/protobuffer"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcTypes "github.com/tendermint/tendermint/rpc/lib/types"
)

type ResultPOSTTxCommit struct {
	CheckTx   abci.ResponseCheckTx   `json:"check_tx"`
	DeliverTx abci.ResponseDeliverTx `json:"deliver_tx"`
	Hash      bytes.HexBytes         `json:"hash"`
	Height    string                 `json:"height"`
}

type RPCPostResponse struct {
}

const (
	nodeConnectionTimeout = time.Second * 600
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

func (nc *NodeConnector) BatchGet(assetID []byte) ([]byte, error) {
	return nil, nil
}

//Load by AssetID
func (nc *NodeConnector) RawGet(assetID []byte) ([]byte, error) {
	txid, err := nc.SingleRawConsensusSearch(assetID)
	if err != nil {
		return nil, err
	}
	query := "tx.hash='" + hex.EncodeToString(txid) + "'"
	result, err := nc.SingleRawChainSearch(query)
	return result, nil
}

//Save (key)
func (nc *NodeConnector) BatchSet(key []byte, serializedData []byte) error {
	return nil
}

func (nc *NodeConnector) Stop() {
	nc.TmClient.Stop()
}

func (nc *NodeConnector) TxSearch(query string, prove bool, currentPage int, numPerPage int) (resultRaw *ctypes.ResultTxSearch, err error) {
	resultRaw, err = nc.TmClient.TxSearch(query, prove, currentPage, numPerPage, "")
	return resultRaw, err
}

// PostTx posts a transaction to the chain and returns the transaction ID
func (nc *NodeConnector) PostTx(asset ChainPostable) (txID string, assetErr *assets.AssetsError) {
	// //serialize the whole transaction
	serializedTX, err := asset.SerializeSignedAsset()
	if err != nil {
		return "", assets.NewAssetsError(assets.CodeTypeEncodingError, "Failed to serialize Asset")

	}
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX)

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
		return "", assets.NewAssetsError(assets.CodeTypeHTTPError, "Failed Connect to QredoChain node")

	}
	req.Header.Set("Content-Type", "text/plain;")

	resp, err := nc.HttpClient.Do(req)
	if err != nil {
		return "", assets.NewAssetsError(assets.CodeTypeHTTPError, "Failed to post to  QredoChain node")
	}
	defer resp.Body.Close()

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", assets.NewAssetsError(assets.CodeTypeHTTPError, "Failed to read from QredoChain node")
	}

	rpcResp := rpcTypes.RPCResponse{}
	err = rpcResp.UnmarshalJSON(jsonResp)
	if err != nil {
		return "", assets.NewAssetsError(assets.CodeTypeEncodingError, "Failed to decode RPC Reeponse from QredoChain node")
	}

	rbtxc := &ResultPOSTTxCommit{}
	print(string(rpcResp.Result))
	err = json.Unmarshal(rpcResp.Result, rbtxc)
	if err != nil {
		return "", assets.NewAssetsError(assets.CodeTypeEncodingError, "Failed to decode RPC Reeponse from QredoChain node")
	}

	if rbtxc.CheckTx.Code != 0 {
		//There was some actionable error
		err := errors.New(string(rbtxc.CheckTx.Data))
		return "", &assets.AssetsError{
			Code:  assets.TransactionCode(rbtxc.CheckTx.Code),
			Error: err,
		}
	}

	//No error code
	if rbtxc.Hash == nil {
		return "", assets.NewAssetsError(assets.TransactionCode(rbtxc.CheckTx.Code), "No TxID returned")
	}

	txID = rbtxc.Hash.String()
	return txID, nil

	// data := f.(map[string]interface{})

	// if data["result"] == nil {
	// 	if data["error"] != nil {
	// 		errdata := data["error"].(map[string]interface{})
	// 		codef64 := errdata["code"].(float64)
	// 		code = assets.TransactionCode(codef64)
	// 		return "", code, errors.New("Failed to add new TX")
	// 	} else {
	// 		return "", assets.CodeAlreadyExists, errors.New("Failed to decode response")
	// 	}
	// }

	// // pp, _ := prettyjson.Marshal(data)
	// //	fmt.Println(string(pp))

	// result := data["result"].(map[string]interface{})
	// txID = result["hash"].(string)
	// checktx := result["check_tx"].(map[string]interface{})
	// codef64 := checktx["code"].(float64)
	// code = assets.TransactionCode(codef64)
	// return txID, code, err
}

func (nc *NodeConnector) GetTx(txHash string) ([]byte, error) {
	//query := fmt.Sprintf("tag.txid='%s'", txHash)
	query := "tag.myname='chris'"
	print("QUERY:", query, "\n")
	result, err := nc.TmClient.TxSearch(query, true, 1, 1, "")
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

func (nc *NodeConnector) GetAsset(assetID string) (*protobuffer.PBSignedAsset, error) {
	//Get TX for Asset ID
	txid, err := nc.ConsensusSearch(assetID, "")
	if err != nil {
		return nil, err
	}
	query := "tx.hash='" + hex.EncodeToString(txid) + "'"
	result, err := nc.QredoChainSearch(query)
	if len(result) != 1 {
		return nil, errors.New("Incorrect number of responses")
	}
	return result[0], nil
}

func (nc *NodeConnector) ConsensusSearch(query string, suffix string) (data []byte, err error) {

	key, err := hex.DecodeString(query)

	if suffix != "" {
		fullSuffix := []byte(suffix)
		key = append(key[:], fullSuffix[:]...)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to decode Base64 Query %s", query)
	}
	return nc.SingleRawConsensusSearch(key)
}

func (nc *NodeConnector) SingleRawConsensusSearch(key []byte) (data []byte, err error) {
	tmClient := nc.TmClient

	result, err := tmClient.ABCIQuery("V", key)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to run Consensus query %s", hex.EncodeToString(key))
	}
	data = result.Response.GetValue()
	return data, nil
}

func (nc *NodeConnector) SingleRawChainSearch(query string) (result []byte, err error) {
	tmClient := nc.TmClient
	r, err := tmClient.TxSearch(query, false, 0, 0, "")
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to SingleRawChainSearch - query %s ", query)
	}
	chainTx := r.Txs[0]
	tx := chainTx.Tx
	return tx, nil
}

func (nc *NodeConnector) QredoChainSearch(query string) (results []*protobuffer.PBSignedAsset, err error) {
	processedCount := 0
	currentPage := 0
	numPerPage := 30

	tmClient := nc.TmClient

	for {
		result, err := tmClient.TxSearch(query, false, currentPage, numPerPage, "asc")
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to search to query %s %d %d", query, currentPage, numPerPage)
		}
		totalToProcess := result.TotalCount
		//fmt.Println("Records Found:", totalToProcess)
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
				if nc.checkQuit(processedCount, totalToProcess) == true {
					return results, nil
				}
				continue
			}
			results = append(results, signedAsset)

			if nc.checkQuit(processedCount, totalToProcess) == true {
				return results, nil
			}
		}
		currentPage++
	}
}

func (nc *NodeConnector) checkQuit(processedCount int, totalToProcess int) bool {
	return processedCount == totalToProcess
}
