package tendermint

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/datastore"
	"github.com/qredo/assets/libs/logger"
	tmclient "github.com/tendermint/tendermint/rpc/client"
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

// PostTx posts a transaction to the chain and returns the transaction ID
func (nc *NodeConnector) PostTx(asset ChainPostable) (txID string, err error) {
	// //serialize the whole transaction
	serializedTX, err := asset.SerializeSignedAsset()
	if err != nil {
		return "", err
	}
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX)

	// // TODO: use net/rpc
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
		return "", errors.Wrap(err, "post to blockchain node")
	}
	req.Header.Set("Content-Type", "text/plain;")

	resp, err := nc.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "post to blockchain node")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respErr string
		if b, err := ioutil.ReadAll(resp.Body); err != nil {
			respErr = resp.Status
		} else {
			respErr = string(b)
		}

		return "", errors.Errorf("Post to blockchain node status %v: %v", resp.StatusCode, respErr)
	}
	nc.log.Debug("Post to chain: Asset %v", asset.Key())
	return
}
