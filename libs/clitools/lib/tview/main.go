package main

import (
	"fmt"

	"github.com/rivo/tview"
)

var box *tview.TextView

func main() {
	//app := tview.NewApplication()

	list := tview.NewList()
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Item %d", i)
		list.AddItem(msg, "", 0, doit)
	}
	list.ShowSecondaryText(false)

	box = tview.NewTextView()
	box.SetText("hello")

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	//menu := newPrimitive("Menu")
	//main := newPrimitive("Main content")
	//main := tview.NewTextView()
	//sideBar := newPrimitive("Side Bar")

	grid := tview.NewGrid()
	grid.SetRows(3, 0)
	grid.SetColumns(30, 0)
	grid.SetBorders(true)
	grid.AddItem(newPrimitive("Header"), 0, 0, 1, 2, 0, 0, false)
	//grid.AddItem(newPrimitive("Footer"), 2, 0, 1, 2, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(list, 1, 0, 1, 1, 0, 100, false)
	//grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false)
	grid.AddItem(box, 1, 1, 1, 1, 0, 100, false)

	//grid.AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	if err := tview.NewApplication().SetRoot(grid, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}

}

func doit() {
	msg := box.GetText(false)
	box.SetText(msg + " hello")
}
