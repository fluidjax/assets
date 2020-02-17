package qc

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
)

func (cliTool *CLITool) CreateWallet(seedHex string) (err error) {

	iddoc, err := cliTool.GetIDDocForSeed(seedHex)
	if err != nil {
		return err
	}

	wallet, err := assets.NewWallet(iddoc, "BTC")
	if err != nil {
		return err
	}
	wallet.Sign(iddoc)

	txid, code, err := cliTool.NodeConn.PostTx(wallet)

	if code != 0 {
		print(err.Error())
		return errors.Wrap(err, "TX Fails verifications")
	}

	if err != nil {
		return err
	}

	addResultTextItem("txid", txid)
	original := reflect.ValueOf(wallet.CurrentAsset)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)
	addResultItem("object", copy.Interface())
	ppResult()

	return
}
