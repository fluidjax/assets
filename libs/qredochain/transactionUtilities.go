package qredochain

import (
	"fmt"
	"strconv"

	"github.com/dgraph-io/badger"
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

func (app *QredoChain) GetWallet(assetID []byte) (*assets.Wallet, error) {
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

func (app *QredoChain) GetGroup(assetID []byte) (*assets.Group, error) {
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

func (app *QredoChain) BatchGet(key []byte) ([]byte, error) {
	var res []byte
	item, err := app.currentBatch.Get(key)
	if item == nil {
		return nil, nil
	}
	err = item.Value(func(val []byte) error {
		res = append([]byte{}, val...) //this copies the item so we can use it outside the closure
		return nil
	})
	return res, err
}

func (app *QredoChain) Get(key []byte) ([]byte, error) {
	var res []byte
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if item == nil {
			return nil
		}
		err = item.Value(func(val []byte) error {
			res = append([]byte{}, val...) //this copies the item so we can use it outside the closure
			return nil
		})
		return err
	})
	return res, err
}

func (app *QredoChain) BatchSet(key []byte, data []byte) error {
	txn := app.currentBatch
	err := txn.Set(key, data)
	return err
}

//Make index 0 padded 8 char
func IndexFormater(index int64) string {
	s := strconv.FormatInt(index, 10)
	return fmt.Sprintf("%08s", s)
}
