package qc

import (
	"context"

	"github.com/rivo/tview"
)

var box *tview.TextView

func listener() {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		case result := <-out:
			_ = result
			//			txhistory = append(txhistory, result)
			doit("got one")
		}
	}
}

func (cliTool *CLITool) Monitor2() (err error) {
	//app := tview.NewApplication()

	connector = cliTool.NodeConn
	out, err = cliTool.NodeConn.TmClient.Subscribe(context.Background(), "", "tx.height>0", 1000)
	if err != nil {
		return err
	}
	go listener()

	list := tview.NewList()
	// for i := 0; i < 100; i++ {
	// 	msg := fmt.Sprintf("Item %d", i)
	// 	list.AddItem(msg, "", 0, doit)
	// }
	list.ShowSecondaryText(false)

	box = tview.NewTextView()
	box.SetText("hello")

	// newPrimitive := func(text string) tview.Primitive {
	// 	return tview.NewTextView().
	// 		SetTextAlign(tview.AlignCenter).
	// 		SetText(text)
	//}

	grid := tview.NewGrid()
	grid.SetRows(0, 0)
	grid.SetColumns(0)
	grid.SetBorders(true)
	//grid.AddItem(newPrimitive("Header"), 0, 0, 1, 2, 0, 0, false)
	//(Primitive, row int, column int, rowSpan int, colSpan int, minGridHeight int, minGridWidth int, focus bool
	grid.AddItem(list, 0, 0, 1, 1, 0, 100, false)
	grid.AddItem(box, 1, 0, 1, 1, 0, 100, false)

	if err := tview.NewApplication().SetRoot(grid, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
	return nil

}

func doit(str string) {
	msg := box.GetText(false)
	box.SetText(msg + str)
}
