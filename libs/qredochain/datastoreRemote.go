package qredochain

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
	rpcTypes "github.com/tendermint/tendermint/rpc/lib/types"
)

func (nc *NodeConnector) BatchGet(assetID []byte) ([]byte, error) {
	return nil, nil
}

func (nc *NodeConnector) GetSignedAsset(key []byte) (*protobuffer.PBSignedAsset, error) {
	val, err := nc.RawGet(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, errors.New("Key not found")
	}
	msg := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(val, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (nc *NodeConnector) GetAssetbyID(assetID []byte) ([]byte, error) { //Get Asset using Asset ID
	txid, err := nc.SingleRawConsensusSearch(assetID)
	if err != nil {
		return nil, err
	}
	query := "tx.hash='" + hex.EncodeToString(txid) + "'"
	result, err := nc.SingleRawChainSearch(query)
	return result, nil
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
func (nc *NodeConnector) Set(key []byte, serializedData []byte) (txID string, err error) {

	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedData)

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
		return "", assets.NewAssetsError(assets.CodeTypeEncodingError, "Failed to decode RPC Response from QredoChain node")
	}

	if rpcResp.Error != nil {
		return "", assets.NewAssetsError(assets.CodeTypeTendermintInternalError, rpcResp.Error.Error())
	}

	rbtxc := &ResultPOSTTxCommit{}
	err = json.Unmarshal(rpcResp.Result, rbtxc)
	if err != nil {
		return "", assets.NewAssetsError(assets.CodeTypeEncodingError, "Failed to decode RPC Response Result from QredoChain node")
	}

	if rbtxc.CheckTx.Code != 0 {
		//There was some actionable error
		err := errors.New(string(rbtxc.CheckTx.Data))
		return "", assets.NewAssetsError(assets.TransactionCode(rbtxc.CheckTx.Code), err.Error())
	}

	//No error code
	if rbtxc.Hash == nil {
		return "", assets.NewAssetsError(assets.TransactionCode(rbtxc.CheckTx.Code), "No TxID returned")
	}

	txID = rbtxc.Hash.String()
	return txID, nil
}
