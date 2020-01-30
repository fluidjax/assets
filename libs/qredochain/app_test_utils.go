package qredochain

//Utilties and globals used to Test the App

import (
	"flag"
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/tendermint/tendermint/node"
)

var done chan bool
var ready chan bool
var app *Qredochain
var tnode *node.Node

func ShutDown() {
	if tnode != nil {
		tnode.Stop()
		tnode.Wait()
	}
}

func InitiateChain() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	app = NewQredochain(db)

	flag.Parse()

	tnode, err := NewTendermint(app, "/tmp/example/config/config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	tnode.Start()
	defer func() {
		tnode.Stop()
		tnode.Wait()
	}()

	done = make(chan bool, 1)
	ready <- true //notify the server is up
	<-done        //wait
}

func StartTestChain() {
	go InitiateChain()
	ready = make(chan bool, 1)
	<-ready //wait for server to come up
}
