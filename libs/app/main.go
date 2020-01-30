package main

//This is a wrapper app that simply starts the Qredochain server

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/qredochain"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

var configFile string

func init() {

	//	flag.StringVar(&configFile, "config", "/home/ubuntu/node/config/config.toml", "Path to config.toml")
	flag.StringVar(&configFile, "config", "/tmp/example/config/config.toml", "Path to config.toml")
}

func main() {

	// dir, err := ioutil.TempDir("/tmp", "qredochain-leveldb") // TODO
	// if err != nil {
	// 	panic("Fail to create database directory")
	// 	os.Exit(1)
	// }

	qredochain := qredochain.NewQredochain("/tmp/qredochain-leveldb")

	flag.Parse()

	node, err := qredochain.NewTendermint(qredochain, configFile)
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
