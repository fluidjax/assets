package lib

import (
	"encoding/hex"

	"github.com/qredo/assets/libs/assets"
)

func (cliTool *CLITool) CreateIDDoc(authref string, broadcast bool) (err error) {

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

	txid := ""
	if broadcast == true {
		var err error
		txid, err = cliTool.NodeConn.PostTx(iddoc)
		if err != nil {
			return err
		}
	}

	serializedAsset, err := iddoc.SerializeAsset()
	if err != nil {
		return err
	}

	//Keep all values internally as Base64 - only convert to Hex to display them
	addResultTextItem("txid", txid)
	addResultBinaryItem("assetid", iddoc.Key())
	addResultBinaryItem("seed", iddoc.Seed)
	addResultBinaryItem("serialized", serializedAsset)
	addResultSignedAsset("object", iddoc.CurrentAsset)

	ppResult()

	return
}
