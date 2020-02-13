package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gookit/color"
	"github.com/qredo/assets/libs/protobuffer"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/urfave/cli"
)
func main() {
	app := cli.NewApp()
	app.Name = "tmget"
	app.Version = "0.1.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Chris Morris",
			Email: "chris@qredo.com",
		},
	}
	app.Copyright = "(c) 2020 Chris Morris"
	app.UsageText = `USAGE:
    cmon chain:port

DESCRIPTION:
	cmon continually monitors the chain for new blocks

EXAMPLE:
	cmon 127.0.0.1:26657
`
	app.Usage = `cmon continually monitors the chain for new blocks `

	app.Action = func(c *cli.Context) error {
		qredochain := c.Args().Get(0)

		if len(c.Args()) != 1 {
			print(app.UsageText)
			os.Exit(1)
			return nil
		}

		tmClient, err := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", qredochain), "/websocket")
		if err := tmClient.Start(); err != nil {
			print("Failed to open websocket")
			os.Exit(1)
		}

		out, err := tmClient.Subscribe(context.Background(), "", "tx.height>0", 1000)
		if err != nil {
			print("Failed to subscribe to node")
			os.Exit(1)
		}

		records := 0

		fmt.Printf("Listening on chain %s\n\n", qredochain)
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
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func getEnv(name, defaultValue string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return v
}

//Use - helper to remove warnings
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

//PadRight - right pad a string
func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}
