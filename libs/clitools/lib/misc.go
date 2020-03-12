package lib

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/prettyjson"
)

// VerifyTX -
func (cliTool *CLITool) VerifyTX(iddocID string, tx string) error {

	// key, err := hex.DecodeString(iddocID)
	// if err != nil {
	// 	return err
	// }

	// res, err := cliTool.NodeConn.GetAsset(iddocID)
	// iddoc, err := assets.ReBuildIDDoc(res, key)

	// msg := &protobuffer.PBSignedAsset{}
	// txbin, err := hex.DecodeString(tx)
	// if err != nil {
	// 	return err
	// }
	// err = proto.Unmarshal(txbin, msg)

	// signedAsset := &assets.SignedAsset{}
	// signedAsset.CurrentAsset = msg
	// signedAsset.DataStore = cliTool.NodeConn

	// err = signedAsset.Verify(iddoc)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (cliTool *CLITool) GenerateSeed() error {
	seed, err := cryptowallet.RandomBytes(48)
	if err != nil {
		return err
	}
	seedHex := hex.EncodeToString(seed)
	fmt.Printf("{\"seed\":\"%s\"}", seedHex)
	return nil
}

func (cliTool *CLITool) Sign(jsonParams string) (err error) {
	//Decode the JSON
	signJSON := &SignJSON{}
	err = json.Unmarshal([]byte(jsonParams), signJSON)
	if err != nil {
		return err
	}

	seed, err := hex.DecodeString(signJSON.Seed)
	if err != nil {
		return err
	}

	key := assets.KeyFromSeed(seed)
	iddoc, err := assets.LoadIDDoc(cliTool.NodeConn, key)
	iddoc.Seed = seed

	msgToSign, err := hex.DecodeString(signJSON.Msg)
	if err != nil {
		return err
	}

	signature, err := assets.Sign(msgToSign, iddoc)

	addResultBinaryItem("signature", signature)
	ppResult()
	return
}

func (cliTool *CLITool) Balance(assetID string) (err error) {
	data, err := cliTool.NodeConn.ConsensusSearch(assetID, ".balance")
	if err != nil {
		return err
	}
	if data == nil {
		return errors.New("Balance doesn't exist")
	}
	currentBalance := int64(binary.LittleEndian.Uint64(data))
	addResultItem("amount", currentBalance)

	original := reflect.ValueOf(res)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)

	pp, _ := prettyjson.Marshal(copy.Interface())
	fmt.Println(string(pp))
	return nil
}

func (cliTool *CLITool) Status(silent bool) (err error) {

	stat, err := cliTool.NodeConn.TmClient.Status()
	if err != nil {
		return nil
	}
	consensusState, err := cliTool.NodeConn.TmClient.ConsensusState()
	if err != nil {
		return err
	}
	health, err := cliTool.NodeConn.TmClient.Health()
	if err != nil {
		return err
	}
	netInfo, err := cliTool.NodeConn.TmClient.NetInfo()
	if err != nil {
		return err
	}

	addResultItem("status", stat)
	addResultItem("ConsensusState", consensusState)
	addResultItem("Health", health)
	addResultItem("NetInfo", netInfo)

	pp, _ := prettyjson.Marshal(res)

	if silent == false {
		print(string(pp))
	}
	return nil
}
