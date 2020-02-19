module github.com/qredo/assets

go 1.13

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.1
	github.com/dgraph-io/badger v1.6.0
	github.com/fatih/color v1.9.0
	github.com/go-kit/kit v0.9.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.3
	github.com/gookit/color v1.2.2
	github.com/hokaccha/go-prettyjson v0.0.0-20190818114111-108c894c2c0e
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jroimartin/gocui v0.4.0
	github.com/mattn/go-runewidth v0.0.8 // indirect
	github.com/nsf/termbox-go v0.0.0-20200204031403-4d2b513ad8be // indirect
	github.com/pkg/errors v0.9.1
	github.com/rivo/tview v0.0.0-20200204110323-ae3d8cac5e4b
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/abci v0.12.0
	github.com/tendermint/tendermint v0.33.0
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/urfave/cli v1.22.2
	github.com/urfave/cli/v2 v2.1.1
	go.etcd.io/bbolt v1.3.3
)

replace github.com/btcsuite/btcd => github.com/qredo/btcd v0.21.1
