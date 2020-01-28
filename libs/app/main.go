package main

//This is a wrapper app that simply starts the Qredochain server

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/qredo/assets/libs/qredochain"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {

	//	flag.StringVar(&configFile, "config", "/home/ubuntu/node/config/config.toml", "Path to config.toml")
	flag.StringVar(&configFile, "config", "/tmp/example/config/config.toml", "Path to config.toml")
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	app := qredochain.NewQredoChain(db)

	flag.Parse()

	node, err := qredochain.NewTendermint(app, configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	node.Start()
	defer func() {
		node.Stop()
		node.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}
