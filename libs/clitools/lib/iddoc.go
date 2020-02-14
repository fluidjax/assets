package qc

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/qredochain"
)

func CreateIDDoc(connectorString string, authref string) (err error) {

	nc, err := qredochain.NewNodeConnector(connectorString, "", nil, nil)
	defer nc.Stop()

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

	txid, code, err := nc.PostTx(iddoc)

	if code != 0 {
		print(err.Error())
		return errors.Wrap(err, "TX Fails verifications")
	}

	if err != nil {
		return err
	}

	//Keep all values internally as Base64 - only convert to Hex to display them
	addResultItem("txid", hex2base64(txid))
	addResultItem("object", iddoc.CurrentAsset)
	addResultItem("assetid", iddoc.Key())
	addResultItem("seed", iddoc.Seed)

	ppResult()
	return
}
