package main

import (
	"log"

	"github.com/gizak/termui/v3"

	"github.com/ynqa/widgets"
	"github.com/ynqa/widgets/node"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatal(err)
	}
	defer termui.Close()

	table := widgets.NewToggleTable()
	table.Headers = []string{"index", "item"}
	table.Widths = []int{10, 10}
	table.Title = "example"

	root := node.Root()
	root.AddChildren(
		node.New("0", []string{"0", "aaa"}).AddChildren(
			node.New("1", []string{"1", "bbb"}).AddChildren(
				node.Leaf("2", []string{"2", "ccc"}),
			),
			node.Leaf("3", []string{"3", "ddd"}),
		),
	)
	table.Node = root

	event := termui.PollEvents()
	setRect := func() {
		width, height := termui.TerminalDimensions()
		table.SetRect(0, 1, width, height-1)
	}
	setRect()
	termui.Render(table)

	for e := range event {
		switch e.ID {
		case "<Enter>":
			root.Toggle(table.SelectedRow)
		case "<Down>":
			table.ScrollDown()
		case "<Up>":
			table.ScrollUp()
		case "q", "<C-c>":
			return
		case "<Resize>":
			setRect()
		}
		termui.Render(table)
	}
}
