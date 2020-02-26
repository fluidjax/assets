package qc

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/prettyjson"
)

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
