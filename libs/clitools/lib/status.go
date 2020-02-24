package qc

import (
	prettyjson "github.com/hokaccha/go-prettyjson"
)

func (cliTool *CLITool) Status(silent bool) (err error) {

	stat, err := cliTool.NodeConn.TmClient.Status()
	if err != nil {
		return nil
	}
	consensusState, err := cliTool.NodeConn.TmClient.ConsensusState()
	if err != nil {
		return err
	}
	health, err := cliTool.NodeConn.TmClient.Health()
	if err != nil {
		return err
	}
	netInfo, err := cliTool.NodeConn.TmClient.NetInfo()
	if err != nil {
		return err
	}

	addResultItem("status", stat)
	addResultItem("ConsensusState", consensusState)
	addResultItem("Health", health)
	addResultItem("NetInfo", netInfo)

	pp, _ := prettyjson.Marshal(res)

	if silent == false {
		print(string(pp))
	}
	return nil
}
