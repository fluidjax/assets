package qc

import (
	"encoding/hex"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/protobuffer"
)

func (cliTool *CLITool) VerifyTX(iddocID string, tx string) error {

	key, err := hex.DecodeString(iddocID)
	if err != nil {
		return err
	}

	res, err := cliTool.NodeConn.GetAsset(iddocID)
	iddoc, err := assets.ReBuildIDDoc(res, key)

	msg := &protobuffer.PBSignedAsset{}
	txbin, err := hex.DecodeString(tx)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(txbin, msg)

	signedAsset := &assets.SignedAsset{}
	signedAsset.CurrentAsset = msg
	signedAsset.DataStore = cliTool.NodeConn

	v, err := signedAsset.Verify(iddoc)
	if err != nil {
		return err
	}

	if v != true {
		return errors.New("TX fails to verify")
	}

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
