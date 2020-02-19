package qc

import "encoding/hex"

func (cliTool *CLITool) Balance(assetID string) (err error) {

	assetIDBin, err := hex.DecodeString(assetID)
	fullSuffix := []byte("." + "balance")
	key := append(assetIDBin[:], fullSuffix[:]...)

	result, err := cliTool.NodeConn.SingleRawConsensusSearch(key)

	print("hello", result)
	return nil
}
