package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
					cliTool.PPQredoChainSearch(query)
					return nil
				},
			},
			&cli.Command{
				Name:    "consensusquery",
				Aliases: []string{"cq"},
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
					if c.NArg() > 0 {
						query = c.Args().Get(0)
					}
					cliTool.PPConsensusSearch(query)
					return nil
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
					return cliTool.CreateIDDoc(authref, broadcast)
				},
			},
			&cli.Command{
				Name:    "createwallet",
				Aliases: []string{"cw"},
				Usage:   "Create a new Wallet with supplied Seed for the already created IDDoc",
				Description: "Generate a new Wallet with the supplied Seed (IDDoc)\n" +
					"   qc cw  dedd7dfb323a7013d7528b3dc753aa5f992c3803f5b183e7df20a5972861dfe7",
				ArgsUsage: "seed",
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
					if c.String("json") != "" {
						return cliTool.CreateWalletWithJSON(c.String("json"), broadcast)

					}
					seed := ""
					if c.NArg() > 0 {
						seed = c.Args().Get(0)
					} else {
						return nil
					}
					return cliTool.CreateWallet(seed, broadcast)
				},
			},
			&cli.Command{
				Name:        "updatewallet",
				Aliases:     []string{"uw"},
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
					if c.String("json") != "" {
						return cliTool.UpdateWalletWithJSON(c.String("json"))
					}
					return nil

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
					if c.String("json") != "" {
						return cliTool.Sign(c.String("json"))
					}
					return nil
				},
			},
			&cli.Command{
				Name:        "aggregatesign",
				Aliases:     []string{"as", "agsign"},
				Usage:       "",
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
					broadcast := c.Bool("broadcast")
					if c.String("json") != "" {
						return cliTool.AggregateSign(c.String("json"), broadcast)
					}
					return nil
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
					err := cliTool.VerifyTX(iddoc, tx)
					if err != nil {
						return cli.Exit(fmt.Sprintf("Verify Fails: %s ", err.Error()), 100)
					}
					return nil
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
					cliTool.GenerateSeed()
					return nil
				},
			},
			&cli.Command{
				Name:            "status",
				Aliases:         []string{"s"},
				Usage:           "Display Qredochain status information",
				Description:     "Show status of QredoChain",
				Flags:           []cli.Flag{},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				Action: func(c *cli.Context) error {
					cliTool.Status()
					return nil
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
					cliTool.Monitor()
					return nil
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
			nc, err := qredochain.NewNodeConnector(connector, "", nil, nil)
			if err != nil {
				return err
			}
			cliTool.NodeConn = nc
			return nil
		},
		After: func(c *cli.Context) error {
			cliTool.NodeConn.Stop()
			return nil
		},
		CommandNotFound: func(c *cli.Context, command string) {
			fmt.Fprintf(c.App.Writer, "Command not found %q. \n", command)
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			print(err)
			if isSubcommand {
				return err
			}

			fmt.Fprintf(c.App.Writer, "WRONG: %#v\n", err)
			return nil
		},
		Action: func(c *cli.Context) error {
			return nil
		},
		Metadata: map[string]interface{}{
			"layers":          "many",
			"explicable":      false,
			"whatever-values": 19.99,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
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
