package qredochain

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func PostTx(base64EncodedTX string, host string) (txID string, err error) {

	// TODO: use net/rpc
	body := strings.NewReader(`{
		"jsonrpc": "2.0",
		"id": "anything",
		"method": "broadcast_tx_commit",
		"params": {
			"tx": "` + base64EncodedTX + `"}
	}`)
	url := "http://" + host

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", errors.Wrap(err, "post to blockchain node")
	}
	req.Header.Set("Content-Type", "text/plain;")

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Do(req)
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
		return "", errors.New("Failed to decode response")
	}
	result := data["result"].(map[string]interface{})
	txID = result["hash"].(string)

	return
}
