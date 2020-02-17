package qc

import (
	"encoding/hex"
	"reflect"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
)

func (cliTool *CLITool) CreateIDDoc(authref string) (err error) {

	if authref == "" {
		randAuth, _ := assets.RandomBytes(32)
		authref = hex.EncodeToString(randAuth)
	}

	iddoc, err := assets.NewIDDoc(authref)
	if err != nil {
		return err
	}
	err = iddoc.Sign(iddoc)

	if err != nil {
		return err
	}

	txid, code, err := cliTool.NodeConn.PostTx(iddoc)

	if code != 0 {
		print(err.Error())
		return errors.Wrap(err, "TX Fails verifications")
	}

	if err != nil {
		return err
	}

	//Keep all values internally as Base64 - only convert to Hex to display them
	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", iddoc.Key())
	addResultBinaryItem("seed", iddoc.Seed)

	//Because json encoding merges binary/string data, and we want binary data converted to
	//hex,  data, we need to convert to hex
	original := reflect.ValueOf(iddoc.CurrentAsset)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)
	addResultItem("object", copy.Interface())
	ppResult()

	return
}
