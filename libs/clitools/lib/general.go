package qc

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gookit/color"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
	tmtypes "github.com/tendermint/tendermint/types"
)

func Status(connectorString string) (err error) {
	nc, err := qredochain.NewNodeConnector(connectorString, "", nil, nil)
	status, err := nc.TmClient.Status()
	if err != nil {
		return err
	}
	consensusState, err := nc.TmClient.ConsensusState()
	if err != nil {
		return err
	}
	health, err := nc.TmClient.Health()
	if err != nil {
		return err
	}
	netInfo, err := nc.TmClient.NetInfo()
	if err != nil {
		return err
	}
	addResultItem("status", status)
	addResultItem("ConsensusState", consensusState)
	addResultItem("Health", health)
	addResultItem("NetInfo", netInfo)
	ppResult()
	return
}

func Monitor(connectorString string) (err error) {
	nc, err := qredochain.NewNodeConnector(connectorString, "", nil, nil)

	out, err := nc.TmClient.Subscribe(context.Background(), "", "tx.height>0", 1000)
	if err != nil {
		print("Failed to subscribe to node")
		os.Exit(1)
	}

	records := 0

	fmt.Printf("Listening on chain %s\n\n", connectorString)
	fmt.Printf("         Blk   Type                 AssetID                                                          Size\n")
	fmt.Printf("--------------------------------------------------------------------------------------------------------\n")
	for {
		select {
		case result := <-out:

			tx := result.Data.(tmtypes.EventDataTx).Tx
			chainData := result.Data.(tmtypes.EventDataTx)
			txsize := len(tx)

			signedAsset := &protobuffer.PBSignedAsset{}
			err := proto.Unmarshal(tx, signedAsset)

			if err != nil {
				panic("Fatal error unmarshalling payload")
			}

			asset := signedAsset.Asset
			typeCol := color.FgWhite.Render
			txType := asset.Type.String()
			assetID := asset.ID

			switch asset.Type {
			case 0: //Wallet
				typeCol = color.FgLightBlue.Render
			case 1: //Group
				typeCol = color.FgCyan.Render
			case 2: //IDDoc
				typeCol = color.FgMagenta.Render
			case 3: //Underlying
				typeCol = color.FgYellow.Render
			case 4: //KVAsset
				typeCol = color.FgBlue.Render
			}
			t := time.Now()
			records++
			blockHeight := PadRight(strconv.FormatInt(chainData.Height, 10), " ", 5)

			assetIDHex := hex.EncodeToString(assetID)
			txType = PadRight(txType, " ", 20)

			fmt.Printf("%s  %s %s %s %d\n", t.Format(time.Kitchen), blockHeight, typeCol(txType), assetIDHex, txsize)
		}
	}
}
