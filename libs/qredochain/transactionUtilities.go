package qredochain

import (
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/gookit/color"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/protobuffer"
)

func dumpMessage(t int, msg string) {
	switch t {
	case 1:
		color.Cyan.Printf("%s\n", msg)
	case 2:
		color.Red.Printf("%s\n", msg)
	case 3:
		color.Yellow.Printf("%s\n", msg)
	case 4:
		color.Green.Printf("%s\n", msg)
	case 5:
		color.Magenta.Printf("%s\n", msg)
	}
}

func (app *QredoChain) GetIDDoc(assetID []byte) (*assets.IDDoc, error) {
	key, err := app.RawGet(assetID)
	if err != nil {
		return nil, err
	}

	signedAssetBytes, err := app.RawGet(key)
	if err != nil {
		return nil, err
	}
	msg := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(signedAssetBytes, msg)
	if err != nil {
		return nil, err
	}
	return assets.ReBuildIDDoc(msg, key)
}

func (app *QredoChain) GetWallet(assetID []byte) (*assets.Wallet, error) {
	key, err := app.RawGet(assetID)
	if err != nil {
		return nil, err
	}
	signedAssetBytes, err := app.RawGet(key)
	if err != nil {
		return nil, err
	}
	msg := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(signedAssetBytes, msg)
	if err != nil {
		return nil, err
	}
	return assets.ReBuildWallet(msg, key)
}

func (app *QredoChain) GetGroup(assetID []byte) (*assets.Group, error) {
	key, err := app.RawGet(assetID)
	if err != nil {
		return nil, err
	}
	signedAssetBytes, err := app.RawGet(key)
	if err != nil {
		return nil, err
	}
	msg := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(signedAssetBytes, msg)
	if err != nil {
		return nil, err
	}
	return assets.ReBuildGroup(msg, key)
}

//Make index 0 padded 8 char
func IndexFormater(index int64) string {
	s := strconv.FormatInt(index, 10)
	return fmt.Sprintf("%08s", s)
}
