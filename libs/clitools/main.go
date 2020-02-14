package main

import (
	"fmt"
	"os"
	"time"

	qc "github.com/qredo/assets/libs/clitools/lib"
	"github.com/urfave/cli/v2"
)

var (
	defaultConnector = "127.0.01:26657"
)

func main() {
	var connector string

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
		ArgsUsage: "[args]",

		Commands: []*cli.Command{
			&cli.Command{
				Name:      "createiddoc",
				Aliases:   []string{"cid", "createid"},
				Category:  "transactions",
				Usage:     "createiddoc (seed) (authenticator_reference)",
				UsageText: "Create a New IDDoc Transaction in the Qredochain using the supplied seed & authenicator reference, ",
				//Description: "",
				//ArgsUsage: "[]",
				// Flags: []cli.Flag{
				// 	&cli.BoolFlag{Name: "forever", Aliases: []string{"forevvarr"}},
				// },
				// Subcommands: []*cli.Command{
				// 	&cli.Command{
				// 		Name:   "wop",
				// 		Action: wopAction,
				// 	},
				// },
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",
				BashComplete: func(c *cli.Context) {
					fmt.Fprintf(c.App.Writer, "--better\n")
				},
				Before: func(c *cli.Context) error {

					//fmt.Fprintf(c.App.Writer, "Before\n")
					return nil
				},
				After: func(c *cli.Context) error {
					//fmt.Fprintf(c.App.Writer, "After\n")
					return nil
				},
				Action: func(c *cli.Context) error {
					fmt.Fprintf(c.App.Writer, "Create an IDDoc\n")
					// c.Command.FullName()
					// c.Command.HasName("wop")
					// c.Command.Names()
					// c.Command.VisibleFlags()
					// fmt.Fprintf(c.App.Writer, "dodododododoodododddooooododododooo\n")
					// if c.Bool("forever") {
					// 	c.Command.Run(c)
					// }
					return nil
				},
				OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
					fmt.Fprintf(c.App.Writer, "for shame\n")
					return err
				},
			},
			&cli.Command{
				Name:    "chainquery",
				Aliases: []string{"cq", "cquery"},
				//Category:    "Query",
				Usage:       "qc chainquery 'search query'",
				Description: "Query the Qredochain (underlying Tendermint) database ",
				//UsageText: "Query the Qredochain (underlying Tendermint) database ",
				ArgsUsage: "querystring",
				//ArgsUsage: "[]",
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
				// Subcommands: []*cli.Command{
				// 	&cli.Command{
				// 		Name:   "wop",
				// 		Action: wopAction,
				// 	},
				// },
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "",

				Action: func(c *cli.Context) error {
					if connector == "" {
						connector = defaultConnector
					}

					query := "empty query"
					if c.NArg() > 0 {
						query = c.Args().Get(0)
					}

					qc.QredoChainSearch(connector, query)
					fmt.Fprintf(c.App.Writer, "Run Query\n")
					return nil
				},
				OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
					fmt.Fprintf(c.App.Writer, "for shame\n")
					return err
				},
			},
		},
		Flags: []cli.Flag{
			// &cli.BoolFlag{Name: "fancy"},
			// &cli.BoolFlag{Value: true, Name: "fancier"},
			// &cli.DurationFlag{Name: "howlong", Aliases: []string{"H"}, Value: time.Second * 3},
			// &cli.Float64Flag{Name: "howmuch"},
			// &cli.GenericFlag{Name: "wat", Value: &genericType{}},
			// &cli.Int64Flag{Name: "longdistance"},
			// &cli.Int64SliceFlag{Name: "intervals"},
			// &cli.IntFlag{Name: "distance"},
			// &cli.IntSliceFlag{Name: "times"},

			// 	&cli.StringSliceFlag{Name: "names", Aliases: []string{"N"}},
			// 	&cli.UintFlag{Name: "age"},
			// 	&cli.Uint64Flag{Name: "bigage"},
		},
		EnableBashCompletion: true,
		HideHelp:             false,
		HideVersion:          false,
		BashComplete: func(c *cli.Context) {
			fmt.Fprintf(c.App.Writer, "Bash Complete General\n")
		},
		Before: func(c *cli.Context) error {
			fmt.Fprintf(c.App.Writer, "Connect to Qredochain node\n")
			return nil
		},
		// After: func(c *cli.Context) error {
		// 	fmt.Fprintf(c.App.Writer, "After App\n")
		// 	return nil
		// },
		CommandNotFound: func(c *cli.Context, command string) {
			fmt.Fprintf(c.App.Writer, "Command not found %q. \n", command)
		},
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}

			fmt.Fprintf(c.App.Writer, "WRONG: %#v\n", err)
			return nil
		},
		Action: func(c *cli.Context) error {
			// cli.DefaultAppComplete(c)
			// cli.HandleExitCoder(errors.New("not an exit coder, though"))
			// cli.ShowAppHelp(c)
			// cli.ShowCommandCompletions(c, "nope")
			// cli.ShowCommandHelp(c, "also-nope")
			// cli.ShowCompletions(c)
			// cli.ShowSubcommandHelp(c)
			// cli.ShowVersion(c)

			// fmt.Printf("%#v\n", c.App.Command("doo"))
			// if c.Bool("infinite") {
			// 	c.App.Run([]string{"app", "doo", "wop"})
			// }

			// if c.Bool("forevar") {
			// 	c.App.RunAsSubcommand(c)
			// }
			// c.App.Setup()
			// fmt.Printf("%#v\n", c.App.VisibleCategories())
			// fmt.Printf("%#v\n", c.App.VisibleCommands())
			// fmt.Printf("%#v\n", c.App.VisibleFlags())

			// fmt.Printf("%#v\n", c.Args().First())
			// if c.Args().Len() > 0 {
			// 	fmt.Printf("%#v\n", c.Args().Get(1))
			// }
			// fmt.Printf("%#v\n", c.Args().Present())
			// fmt.Printf("%#v\n", c.Args().Tail())

			// set := flag.NewFlagSet("contrive", 0)
			// nc := cli.NewContext(c.App, set, c)

			// fmt.Printf("%#v\n", nc.Args())
			// fmt.Printf("%#v\n", nc.Bool("nope"))
			// fmt.Printf("%#v\n", !nc.Bool("nerp"))
			// fmt.Printf("%#v\n", nc.Duration("howlong"))
			// fmt.Printf("%#v\n", nc.Float64("hay"))
			// fmt.Printf("%#v\n", nc.Generic("bloop"))
			// fmt.Printf("%#v\n", nc.Int64("bonk"))
			// fmt.Printf("%#v\n", nc.Int64Slice("burnks"))
			// fmt.Printf("%#v\n", nc.Int("bips"))
			// fmt.Printf("%#v\n", nc.IntSlice("blups"))
			// fmt.Printf("%#v\n", nc.String("snurt"))
			// fmt.Printf("%#v\n", nc.StringSlice("snurkles"))
			// fmt.Printf("%#v\n", nc.Uint("flub"))
			// fmt.Printf("%#v\n", nc.Uint64("florb"))

			// fmt.Printf("%#v\n", nc.FlagNames())
			// fmt.Printf("%#v\n", nc.IsSet("wat"))
			// fmt.Printf("%#v\n", nc.Set("wat", "nope"))
			// fmt.Printf("%#v\n", nc.NArg())
			// fmt.Printf("%#v\n", nc.NumFlags())
			// fmt.Printf("%#v\n", nc.Lineage()[1])
			// nc.Set("wat", "also-nope")

			// ec := cli.Exit("ohwell", 86)
			// fmt.Fprintf(c.App.Writer, "%d", ec.ExitCode())
			// fmt.Printf("made it!\n")
			// return ec
			return nil
		},
		Metadata: map[string]interface{}{
			"layers":          "many",
			"explicable":      false,
			"whatever-values": 19.99,
		},
	}

	// if os.Getenv("HEXY") != "" {
	// 	app.Writer = &hexWriter{}
	// 	app.ErrWriter = &hexWriter{}
	// }

	app.Run(os.Args)
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
