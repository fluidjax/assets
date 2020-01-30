/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package assets

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//RandomBytes - generate n random bytes
func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

func contains(s [][]byte, e []byte) bool {
	for _, a := range s {
		res := bytes.Compare(a, e)
		if res == 0 {
			return true
		}
	}
	return false
}

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
	result := data["result"].(map[string]interface{})
	txID = result["hash"].(string)

	return
}
