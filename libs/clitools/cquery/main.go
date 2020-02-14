package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hokaccha/go-prettyjson"
	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/protobuffer"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/urfave/cli"
)

func main() {
	//hello
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
	app.Copyright = "(c) 2019 Chris Morris"
	app.UsageText = `USAGE:
    cquery configdir query

DESCRIPTION:
    cquery queries the chain specified in the config directory with the query and dumps the decoded output

EXAMPLE:
	Qredochain - Tendermint Query
		cquery 127.0.0.1:26657 Q "tx.hash='528579CDD20444140270C5B476AA2971A484719C7BE02CB99539468AEC93B222'"
		cquery 127.0.0.1:26657 Q "tx.height>0 and tx.height<10"

		// Where tag added using code such as 
		// i.AddTag("tagkey", []byte("tagvalue"))
		cquery 127.0.0.1:26657 Q "tag.tagkey='tagvalue'"
		cquery 127.0.0.1:26657 Q "tag.tagkey contains 'tag'"


	Consensus Query
		cquery 127.0.0.1:26657 C "nO3lRBxbYjbEclTiK7joo7uBPObh1CZbB36VHriuSoo="


query /Users/chris/.milagro "tag.recipient='Au1WipqVeTx9i2PV4UcCxmY6iQvA9RZXy88xJLRzafwc'" 
	cquery /Users/chris/.milagro "tag.reference contains 'Eat'" 

	tx.height		- block height
	tag.txhash		- hash of the transsaction
	tag.txtype      - Document Type (none=0,Order=1,FulfillOrder=2,OrderOutput=3,OrderSecret=4,
									 FulfillOrderSecret=5,OrderSecretOutput=6,Identity=7,TrusteeSecret=8)
								     (see protobuffer/proto.go for any additional types)
	tag.sender      - ID of the sender
	tag.recipient   - ID of the recipeitn
	tag.reference   - Order Reference

`

	app.Usage = `cquery queries the chain specified in the config directory with the query and dumps the decoded output `

	app.Action = func(c *cli.Context) error {

		qredochain := c.Args().Get(0)
		querytype := c.Args().Get(1)
		query := c.Args().Get(2)

		if len(c.Args()) != 3 {
			print(app.UsageText)
			os.Exit(1)
			return nil
		}

		switch querytype {
		case "Q":
			return QredoChainSearch(qredochain, query)
		case "C":
			return ConsensusSearch(qredochain, query)
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

