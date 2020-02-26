package qc

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/elliotchance/orderedmap"
	"github.com/gogo/protobuf/proto"
	"github.com/gookit/color"
	"github.com/jroimartin/gocui"
	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/prettyjson"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/qredo/assets/libs/qredochain"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type col struct {
	name string
	pad  int
}

var mux sync.Mutex
var coldefs = orderedmap.NewOrderedMap()

type QredoChainTX struct {
	result  *ctypes.ResultEvent
	balance string
}

var connector *qredochain.NodeConnector
var datalines []*QredoChainTX
var latched = false

//Monitor - Monitor the chain in real time
func (cliTool *CLITool) Monitor() (err error) {

	coldefs.Set("num", col{"Num", 8})
	coldefs.Set("time", col{"Time", 6})
	coldefs.Set("block", col{"Blk", 8})
	coldefs.Set("type", col{"Type", 14})
	coldefs.Set("assetid", col{"AssetID", 64})
	coldefs.Set("size", col{"Size", 4})
	coldefs.Set("amount", col{"Amount", 8})

	connector = cliTool.NodeConn
	out, err = cliTool.NodeConn.TmClient.Subscribe(context.Background(), "", "tx.height>0", 1000)
	if err != nil {
		return err
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
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	wg.Wait()
	return nil
}

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

	main, _ := g.View("main")
	info, _ := g.View("info")
	main.Highlight = true
	main.Cursor()

	// _, y := main.Cursor()
	// if y == 0 {
	// 	main.SetCursor(0, 3)
	// }

	main.Title = "Transactions "
	info.Title = "Detail"

	return nil
}
func DetailUp(g *gocui.Gui, v *gocui.View) error {
	info, _ := g.View("info")
	main, _ := g.View("main")
	x, y := info.Origin()
	info.SetOrigin(x, y-1)

	displayDetail(g, main)
	return nil
}

func Home(g *gocui.Gui, v *gocui.View) error {
	displayTopItem = 0
	main, _ := g.View("main")
	main.SetCursor(0, 1)
	latched = false
	showList(g)
	return nil
}

func End(g *gocui.Gui, v *gocui.View) error {
	main, _ := g.View("main")
	_, maxY := main.Size()
	maxY--

	if len(datalines) > maxY {
		//More than a full screen
		displayTopItem = max(len(datalines)-maxY, 0)
		main.SetCursor(0, maxY)
		latched = true

	} else {
		//not yet a full screen
		endItem := listEndItem(g)
		main.SetCursor(0, endItem)

	}

	showList(g)
	return nil
}

func DetailDown(g *gocui.Gui, v *gocui.View) error {
	info, _ := g.View("info")
	main, _ := g.View("main")
	x, y := info.Origin()
	info.SetOrigin(x, y+1)
	displayDetail(g, main)
	return nil
}

func ListUp(g *gocui.Gui, v *gocui.View) error {
	//info, _ := g.View("main")
	latched = false
	main, _ := g.View("main")

	cx, cy := main.Cursor()

	if cy > 1 {
		main.SetCursor(cx, cy-1)
	} else {
		displayTopItem--
		displayTopItem = max(displayTopItem, 0)
	}

	displayDetail(g, main)
	showList(g)
	return nil
}

func ListDown(g *gocui.Gui, v *gocui.View) error {
	main, _ := g.View("main")
	cx, cy := main.Cursor()
	_, sy := main.Size()
	endItem := listEndItem(g)

	// if endItem == len(datalines) && endItem == cy+displayTopItem+1 {
	// 	latched = true
	// }

	if cy == sy-1 {
		if endItem != len(datalines) {
			//Bottom of screen
			displayTopItem++
			main.SetCursor(cx, cy)
		} else {

		}
	} else if endItem == cy {
		//nothing
	} else {
		main.SetCursor(cx, cy+1)
	}
	displayDetail(g, main)
	showList(g)
	return nil
}

func listEndItem(g *gocui.Gui) int {
	main, _ := g.View("main")
	_, mainHeight := main.Size()
	return min(displayTopItem+len(datalines), displayTopItem+mainHeight-1)
}

func keybindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, ListUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, ListDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, DetailUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, DetailDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, Home); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, End); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.MouseLeft, gocui.ModNone, displayDetail); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, displayDetail); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

type head struct {
	name string
	pad  int
}

func addItemToScreen(g *gocui.Gui, res QredoChainTX) {
	g.Update(func(g *gocui.Gui) error {
		showList(g)
		return nil
	})
}

func showList(g *gocui.Gui) error {

	main, err := g.View("main")
	main.Editable = true
	main.Highlight = true
	main.Wrap = true
	if err != nil {
		return err
	}
	main.Clear()

	endItem := listEndItem(g)

	if latched == true {
		main.SelBgColor = gocui.ColorRed
		main.SelFgColor = gocui.ColorWhite
	} else {
		main.SelBgColor = gocui.ColorGreen
		main.SelFgColor = gocui.ColorBlack
	}

	lastDisplayLine := min(endItem, len(datalines))

	headerCol := color.New(color.FgWhite, color.BgBlue, color.OpBold)

	for _, key := range coldefs.Keys() {
		c, _ := coldefs.Get(key)
		col := c.(col)
		lab := PadRight(col.name, " ", col.pad)
		lab = headerCol.Sprintf("%s ", lab)
		fmt.Fprint(main, lab)
	}
	fmt.Fprint(main, "\n")

	for i := displayTopItem; i < lastDisplayLine; i++ {
		showdatalinesLine(main, datalines[i], i+1)
	}
	// if latched == true {
	// 	main.SetCursor(0, lastDisplayLine-1)
	// }

	return nil
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y

}

func showdatalinesLine(main *gocui.View, qc *QredoChainTX, count int) {
	res := qc.result
	tx := res.Data.(tmtypes.EventDataTx).Tx
	chainData := res.Data.(tmtypes.EventDataTx)
	txsize := fmt.Sprintf("%d", len(tx))
	signedAsset := &protobuffer.PBSignedAsset{}
	err := proto.Unmarshal(tx, signedAsset)
	if err != nil {
		panic("Fatal error unmarshalling payload")
	}
	asset := signedAsset.Asset
	txType := asset.Type.String()

	if asset.Index > 1 {
		txType = "Updated" + txType
	}
	assetID := asset.ID
	assetIDHex := hex.EncodeToString(assetID)

	//Determine Balance on first run
	if qc.balance == "" {
		qc.balance = "-"

		switch asset.Type {
		case protobuffer.PBAssetType_Underlying:
			//Underlying - set balance to show how much addede
			underlying, err := assets.ReBuildUnderlying(signedAsset, assetID)
			if err != nil {
				break
			}
			payload, err := underlying.Payload()
			if err != nil {
				break
			}
			qc.balance = color.Green.Sprintf("+%d", payload.Amount)

		case protobuffer.PBAssetType_Wallet:
			data, err := connector.ConsensusSearch(assetIDHex, ".balance")
			if err != nil {
				panic("Fatal error retrieving balance")
			}
			if len(data) == 8 {

				balance := int64(binary.LittleEndian.Uint64(data))
				balanceStr := PadRight(fmt.Sprintf("%d", balance), " ", getCol("amount").pad)
				qc.balance = color.Cyan.Sprintf("%s", balanceStr)
			}
		}
	}

	t := time.Now()
	blockHeight := PadRight(strconv.FormatInt(chainData.Height, 10), " ", 5)
	countStr := fmt.Sprintf("%d", count)

	f1 := color.White.Sprintf("%s", PadRight(countStr, " ", getCol("num").pad))
	f2 := color.White.Sprintf("%s", PadRight(t.Format(time.Kitchen), " ", getCol("time").pad))
	f3 := color.White.Sprintf("%s", PadRight(blockHeight, " ", getCol("block").pad))

	tyTypeColor := color.White
	switch asset.Type {
	case protobuffer.PBAssetType_Group:
		tyTypeColor = color.Red
	case protobuffer.PBAssetType_Iddoc:
		tyTypeColor = color.Yellow
	case protobuffer.PBAssetType_Underlying:
		tyTypeColor = color.Cyan
	case protobuffer.PBAssetType_KVAsset:
		tyTypeColor = color.Magenta
	case protobuffer.PBAssetType_Wallet:
		tyTypeColor = color.Green
	case protobuffer.PBAssetType_MPC:
		tyTypeColor = color.Red
	}
	f4 := tyTypeColor.Sprintf("%s", PadRight(txType, " ", getCol("type").pad))
	f5 := color.White.Sprintf("%s", PadRight(assetIDHex, " ", getCol("assetid").pad))
	f6 := color.White.Sprintf("%s", PadRight(txsize, " ", getCol("size").pad))
	f7 := qc.balance

	fmt.Fprintln(main, f1, f2, f3, f4, f5, f6, f7)

}

func getCol(key string) col {
	c, _ := coldefs.Get(key)
	return c.(col)

}

func dumpStatus(g *gocui.Gui) {
	main, _ := g.View("main")
	info, _ := g.View("info")

	_, y := main.Origin()
	fmt.Fprintf(info, "Main Origin:Y   %d\n", y)
	_, y = main.Cursor()
	fmt.Fprintf(info, "Main Cursor:Y   %d \n", y)
	_, y = main.Size()
	fmt.Fprintf(info, "Main Size:Y   %d \n", y)
	fmt.Fprintf(info, "displayTopItem   %d \n", displayTopItem)
	fmt.Fprintf(info, "endItem   %d \n", listEndItem(g))
	fmt.Fprintf(info, "length History   %d \n", len(datalines))
	if latched == true {
		fmt.Fprintf(info, "Latched\n")
	}

}

func displayDetail(g *gocui.Gui, main *gocui.View) error {

	info, err := g.View("info")
	info.Clear()
	info.Editable = true
	info.Wrap = true

	if err != nil {
		return err
	}

	_, y := main.Cursor()
	itemNumber := y + displayTopItem - 1

	if itemNumber < 0 || itemNumber >= len(datalines) {
		return nil
	}
	res := datalines[itemNumber].result
	tx := res.Data.(tmtypes.EventDataTx).Tx
	signedAsset := &protobuffer.PBSignedAsset{}
	err = proto.Unmarshal(tx, signedAsset)
	if err != nil {
		return nil
	}

	//fmt.Fprintf(info, "height %d\n", sizeY)

	chainData := res.Data.(tmtypes.EventDataTx)
	hash := chainData.Tx.Hash()
	result := make(map[string]string)
	result["TxID"] = hex.EncodeToString(hash)

	original := reflect.ValueOf(result)
	copy := reflect.New(original.Type()).Elem()
	TranslateRecursive(copy, original)
	pp, _ := prettyjson.Marshal(copy.Interface())
	fmt.Fprintf(info, string(pp))

	//Modify specific fields for ease of display
	//asset := signedAsset.Asset
	// switch asset.Type {
	// case protobuffer.PBAssetType_Group:
	// 	group := signedAsset.Asset.GetGroup()

	// 	data := hex.EncodeToString([]byte("AAA"))
	// 	group.Participants["groupfield_key1"] = []byte(data)

	// case protobuffer.PBAssetType_Iddoc:
	// case protobuffer.PBAssetType_Underlying:
	// case protobuffer.PBAssetType_KVAsset:
	// case protobuffer.PBAssetType_Wallet:
	// case protobuffer.PBAssetType_MPC:
	// }

	fmt.Fprintf(info, prettyStringFromSignedAsset(signedAsset))

	return nil
}

func txListener(g *gocui.Gui) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		case incoming := <-out:
			mux.Lock()
			res := QredoChainTX{
				result:  &incoming,
				balance: "",
			}
			datalines = append(datalines, &res)

			g.Update(func(g *gocui.Gui) error {
				main, _ := g.View("main")
				_, maxY := main.Size()

				_ = maxY
				if latched == true && len(datalines) > maxY-1 {
					displayTopItem++
				}
				showList(g)
				addItemToScreen(g, res)
				displayDetail(g, main)
				mux.Unlock()
				return nil
			})
		}
	}
}
