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

func (app *Qredochain) GetIDDoc(assetID []byte) (*assets.IDDoc, error) {
	key, err := app.Get(assetID)
	if err != nil {
		return nil, err
	}

	signedAssetBytes, err := app.Get(key)
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

func (app *Qredochain) GetWallet(assetID []byte) (*assets.Wallet, error) {
	key, err := app.Get(assetID)
	if err != nil {
		return nil, err
	}
	signedAssetBytes, err := app.Get(key)
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

func (app *Qredochain) GetGroup(assetID []byte) (*assets.Group, error) {
	key, err := app.Get(assetID)
	if err != nil {
		return nil, err
	}
	signedAssetBytes, err := app.Get(key)
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

func (app *Qredochain) Get(key []byte) ([]byte, error) {
	return app.state.db.Get(key)
}

func (app *Qredochain) Set(key []byte, data []byte) error {
	return app.state.db.Set(key, data)
}

func (app *Qredochain) exists(key []byte) bool {
	item, _ := app.Get(key)
	return item != nil
}

func KeySuffix(key []byte, suffix string) []byte {
	suffixBytes := []byte(suffix)
	return append(key, suffixBytes...)
}

//Make index 0 padded 8 char
func IndexFormater(index int64) string {
	s := strconv.FormatInt(index, 10)
	return fmt.Sprintf("%08s", s)
}

func decodeTX(data []byte) (*protobuffer.PBSignedAsset, error) {
	signedAsset := &protobuffer.PBSignedAsset{}

	err := proto.Unmarshal(data, signedAsset)
	if err != nil {
		return nil, err
	}
	return signedAsset, nil
}
