package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gookit/color"
	qc "github.com/qredo/assets/libs/clitools/lib"
	"github.com/qredo/assets/libs/qredochain"
	"github.com/urfave/cli/v2"
)

var (
	defaultConnector = "127.0.0.1:26657"
)

func main() {
	var connector string
	cliTool := &qc.CLITool{}

	app := &cli.App{
		Name:     "Qredochain",
		Version:  "v0.1.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Chris Morris",
				Email: "chris@qredo.com",
			},
		},
		Copyright: "(c) 2020 Qredo Ltd",
		HelpName:  "qredochain",
		Usage:     "qredochain command params",
		UsageText: "Interact with the Qredochain from the command line",
		//ArgsUsage: "[args]",

		Commands: []*cli.Command{
			&cli.Command{
				Name:    "tendermintquery",
				Aliases: []string{"tq"},
				Usage:   "Search the underlying tendermint database",
				Description: "Query the underlying tendermint database \n" +
					"   examples:\n" +
					"   qc tq \"tx.hash='528579CDD20444140270C5B476AA2971A484719C7BE02CB99539468AEC93B222'\"\n" +
					"   qc tq \"tx.height>0 and tx.height<10\"\n" +
					"   qc tq \"tag.tagkey='tagvalue'\"\n" +
					"   qc tq \"tag.tagkey contains 'tag'\"\n",
				ArgsUsage:       "querystring",
				Flags:           []cli.Flag{},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					query := "empty query"
					if c.NArg() > 0 {
						query = c.Args().Get(0)
					}
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.PPQredoChainSearch(query)
				},
			},
			&cli.Command{
				Name:    "consensusquery",
				Aliases: []string{"cq", "qq"},
				Usage:   "Search the Qredochain Consensus App database for keys ",
				Description: "Query the Qredochain Consensus App Layer database\n" +
					"   examples:\n" +
					"   qc cq  \"nO3lRBxbYjbEclTiK7joo7uBPObh1CZbB36VHriuSoo=\"\n",
				ArgsUsage:       "querystring",
				Flags:           []cli.Flag{},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					query := "empty query"
					suffix := ""
					if c.NArg() > 0 {
						query = c.Args().Get(0)
					}
					if c.NArg() > 1 {
						suffix = c.Args().Get(1)
					}
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.PPConsensusSearch(query, suffix)
				},
			},
			&cli.Command{
				Name:    "createiddoc",
				Aliases: []string{"cid"},
				Usage:   "Create a new Identity Doc (IDDoc) using a random seed",
				Description: "Generate a new IDDoc with optional supplied authentication reference. Failure to supply an authentication reference will result in a random one being generated.\n" +
					"   qc cid  \"testid\"",
				ArgsUsage: "authref",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "broadcast",
						Aliases: []string{"b"},
						Value:   false,
						Usage:   "broadcast to the Qredo Network",
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					broadcast := c.Bool("broadcast")
					authref := ""
					if c.NArg() > 0 {
						authref = c.Args().Get(0)
					}
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.CreateIDDoc(authref, broadcast)
				},
			},
			&cli.Command{
				Name:        "createwallet",
				Aliases:     []string{"cw"},
				Usage:       "Create a new Wallet",
				Description: "Generate a new Wallet with the supplied JSON",
				ArgsUsage:   "",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "specify json parameters",
					},
					&cli.BoolFlag{
						Name:    "broadcast",
						Aliases: []string{"b"},
						Value:   false,
						Usage:   "broadcast to the Qredo Network",
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					broadcast := c.Bool("broadcast")
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.CreateWalletWithJSON(c.String("json"), broadcast)

				},
			},
			&cli.Command{
				Name:        "prepwalletupdate",
				Aliases:     []string{"pwu"},
				Usage:       "Prepare a New Wallet Update Transaction for Signing",
				Description: "Prepare a New Wallet Update Transaction for Signing",
				ArgsUsage:   "",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "specify json parameters",
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.PrepareWalletUpdateWithJSON(c.String("json"))

				},
			},
			&cli.Command{
				Name:        "sendunderlying",
				Aliases:     []string{"su"},
				Usage:       "Create a new Underlying Transaction",
				Description: "Create a new Underlying Transaction",
				ArgsUsage:   "",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "specify json parameters",
					},
					&cli.BoolFlag{
						Name:    "broadcast",
						Aliases: []string{"b"},
						Value:   false,
						Usage:   "broadcast to the Qredo Network",
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					broadcast := c.Bool("broadcast")
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}

					return cliTool.CreateUnderlyingWithJSON(c.String("json"), broadcast)

				},
			},
			&cli.Command{
				Name:        "sendmpc",
				Aliases:     []string{"smpc"},
				Usage:       "Create a new MPC Transaction",
				Description: "Create a new MPC Transaction",
				ArgsUsage:   "",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "specify json parameters",
					},
					&cli.BoolFlag{
						Name:    "broadcast",
						Aliases: []string{"b"},
						Value:   false,
						Usage:   "broadcast to the Qredo Network",
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					if cliTool.NodeConn == nil {
						return errors.New("Error: Fail to connect to Node - " + cliTool.ConnectString)
					}
					broadcast := c.Bool("broadcast")
					return cliTool.CreateMPCWithJSON(c.String("json"), broadcast)
				},
			},
			&cli.Command{
				Name:        "sign",
				Aliases:     []string{"s"},
				Usage:       "Update an existing wallet with a transfer",
				Description: "Update an existing wallet with a transfer\n",
				ArgsUsage:   "",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "specify json parameters",
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}

					return cliTool.Sign(c.String("json"))

				},
			},
			&cli.Command{
				Name:        "sendwallet",
				Aliases:     []string{"sw", ""},
				Usage:       "Aggregate Sign, Check & Broadcast Wallet update",
				Description: "",
				ArgsUsage:   "",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "json",
						Aliases: []string{"j"},
						Usage:   "specify json parameters",
					},
					&cli.BoolFlag{
						Name:    "broadcast",
						Aliases: []string{"b"},
						Value:   false,
						Usage:   "broadcast to the Qredo Network",
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}

					broadcast := c.Bool("broadcast")
					return cliTool.AggregateWalletSign(c.String("json"), broadcast)
				},
			},
			&cli.Command{
				Name:            "verifytx",
				Aliases:         []string{"vtx"},
				Usage:           "Verify a Raw TX",
				Description:     "Decode and verify a raw TX",
				Flags:           []cli.Flag{},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}

					iddoc := ""
					tx := ""
					if c.NArg() > 0 {
						iddoc = c.Args().Get(0)
					} else {
						return nil
					}
					if c.NArg() > 1 {
						tx = c.Args().Get(1)
					} else {
						return nil
					}
					return cliTool.VerifyTX(iddoc, tx)
				},
			},
			&cli.Command{
				Name:            "generateseed",
				Aliases:         []string{"gs"},
				Usage:           "Generate a random seed value",
				Description:     "Generate a random seed value",
				Flags:           []cli.Flag{},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.GenerateSeed()

				},
			},
			&cli.Command{

				Name:            "balance",
				Aliases:         []string{"bal"},
				Usage:           "Display Balance for Supplied AssetID",
				Description:     "Get balance(s) for supplied AssetID",
				ArgsUsage:       "assetid",
				Flags:           []cli.Flag{},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					assetid := ""
					if c.NArg() > 0 {
						assetid = c.Args().Get(0)
					}
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.Balance(assetid)
				},
			},
			&cli.Command{
				Name:            "status",
				Aliases:         []string{"s"},
				Usage:           "Display Qredochain status information",
				Description:     "Show status of QredoChain",
				Flags:           []cli.Flag{},
				SkipFlagParsing: false,
				HideHelp:        true,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.Status()
				},
			},
			&cli.Command{
				Name:            "monitor",
				Aliases:         []string{"m", "mon"},
				Usage:           "Monitor the Qredochain for incoming transactions",
				Description:     "Monitor the Qredochain for incoming transactions",
				Flags:           []cli.Flag{},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					if cliTool.NodeConn == nil {
						return errors.New("Fail to connect to Node: " + cliTool.ConnectString)
					}
					return cliTool.Monitor()
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "connect",
				Aliases:     []string{"c"},
				Destination: &connector,
				Usage:       "Qredochain connection string",
				DefaultText: defaultConnector,
				Required:    false,
			},
		},
		EnableBashCompletion: false,
		HideHelp:             false,
		HideVersion:          false,
		Before: func(c *cli.Context) error {
			if connector == "" {
				connector = defaultConnector
			}
			nc, _ := qredochain.NewNodeConnector(connector, "", nil, nil)
			//we trap any failure later - else we print the help screen everytime
			cliTool.NodeConn = nc
			cliTool.ConnectString = connector
			return nil
		},
		After: func(c *cli.Context) error {
			if cliTool.NodeConn != nil {
				cliTool.NodeConn.Stop()
			}
			return nil
		},
		CommandNotFound: func(c *cli.Context, command string) {
			fmt.Fprintf(c.App.Writer, "Command not found %q. \n", command)
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}
			//fmt.Fprintf(c.App.Writer, "WRONG: %#v\n", err)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		color.Red.Println(err.Error())
		os.Exit(1)
		return
	}
}

func wopAction(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, ":wave: over here, eh\n")
	return nil
}

type hexWriter struct{}

func (w *hexWriter) Write(p []byte) (int, error) {
	for _, b := range p {
		fmt.Printf("%x", b)
	}
	fmt.Printf("\n")

	return len(p), nil
}

type genericType struct {
	s string
}

func (g *genericType) Set(value string) error {
	g.s = value
	return nil
}

func (g *genericType) String() string {
	return g.s
}
