package qc

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jroimartin/gocui"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	nc                *qredochain.NodeConnector
	out               <-chan ctypes.ResultEvent
	done              = make(chan struct{})
	wg                sync.WaitGroup
	mu                sync.Mutex // protects ctr
	ctr               = 0
	displayTopItem    = 0
	displayBottomItem = 0
	highlightRow      = 0
	txhistory         []ctypes.ResultEvent
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

	nc, err = qredochain.NewNodeConnector(connectorString, "", nil, nil)

	out, err = nc.TmClient.Subscribe(context.Background(), "", "tx.height>0", 1000)
	if err != nil {
		print("Failed to subscribe to node")
		os.Exit(1)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	g.Mouse = true
	g.Highlight = true
	g.Cursor = true
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	wg.Add(1)

	go txListener(g)
	//initial draw

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
	return nil
}

// func layout(g *gocui.Gui) error {
// 	if v, err := g.SetView("ctr", 2, 2, 12, 4); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		fmt.Fprintln(v, "0")
// 	}
// 	return nil
//}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("main", 0, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if _, err := g.SetView("info", 0, maxY/2-1, maxX-1, maxY-1); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.MouseLeft, gocui.ModNone, showTX); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func txListener(g *gocui.Gui) {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		case result := <-out:
			txhistory = append(txhistory, result)
			addItemToScreen(g, result)
		}
	}
}

type head struct {
	name string
	pad  int
}

func showMainHeader(main *gocui.View) {
	str := ""
	cm := []head{
		{"Blk", 3},
		{"Type", 16},
		{"AssetID", 40},
		{"Size", 5},
	}
	for _, val := range cm {
		str += fmt.Sprintf("%s%s", val.name, strings.Repeat(" ", val.pad))
	}
	fmt.Fprintln(main, str)

}

func addItemToScreen(g *gocui.Gui, res ctypes.ResultEvent) {
	g.Update(func(g *gocui.Gui) error {
		main, err := g.View("main")
		main.Editable = true
		main.Highlight = true
		main.Wrap = true
		main.SelBgColor = gocui.ColorGreen
		main.SelFgColor = gocui.ColorBlack
		main.Autoscroll = true
		if err != nil {
			return err
		}
		showTXHistoryLine(main, res)
		return nil
	})
}

func showTXHistoryLine(main *gocui.View, res ctypes.ResultEvent) {
	tx := res.Data.(tmtypes.EventDataTx).Tx
	chainData := res.Data.(tmtypes.EventDataTx)
	txsize := len(tx)
	signedAsset := &protobuffer.PBSignedAsset{}
	err := proto.Unmarshal(tx, signedAsset)
	if err != nil {
		panic("Fatal error unmarshalling payload")
	}
	asset := signedAsset.Asset
	txType := asset.Type.String()
	assetID := asset.ID
	t := time.Now()
	blockHeight := PadRight(strconv.FormatInt(chainData.Height, 10), " ", 5)
	assetIDHex := hex.EncodeToString(assetID)
	fmt.Fprintf(main, " %s  %s %s %s %d\n", t.Format(time.Kitchen), blockHeight, txType, assetIDHex, txsize)

}
func showTX(g *gocui.Gui, main *gocui.View) error {
	info, err := g.View("info")
	info.Clear()
	info.Editable = true
	info.Wrap = true
	if err != nil {
		return err
	}
	_, y := main.Cursor()
	_, sizeY := main.Size()
	historyLength := len(txhistory)
	var itemNumber int
	if len(txhistory) < sizeY {
		itemNumber = y
	} else {
		itemNumber = historyLength - (sizeY - y + 1)
	}
	if itemNumber > historyLength-1 {
		return nil
	}
	res := txhistory[itemNumber]
	tx := res.Data.(tmtypes.EventDataTx).Tx
	signedAsset := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(tx, signedAsset)
	if err != nil {
		return nil
	}
	fmt.Fprintf(info, prettyStringFromSignedAsset(signedAsset))
	return nil
}
