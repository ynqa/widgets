package main

import (
	"log"

	"github.com/gizak/termui/v3"

	"github.com/ynqa/widgets/pkg/node"
	"github.com/ynqa/widgets/pkg/widgets"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatal(err)
	}
	defer termui.Close()

	toggle := widgets.NewToggle()
	toggle.Headers = []string{"index", "item"}
	toggle.Widths = []int{10, 10}
	toggle.Title = "example"

	root := node.Root()
	root.Append(
		node.New("0", []string{"0", "aaa"}).Append(
			node.New("1", []string{"1", "bbb"}).Append(
				node.New("2", []string{"2", "ccc"}),
			),
			node.New("3", []string{"3", "ddd"}),
		),
	)
	toggle.Node = root

	event := termui.PollEvents()
	setRect := func() {
		width, height := termui.TerminalDimensions()
		toggle.SetRect(0, 1, width, height-1)
	}
	setRect()
	termui.Render(toggle)

	for e := range event {
		switch e.ID {
		case "<Enter>":
			root.Toggle(toggle.SelectedRow)
		case "<Down>":
			toggle.ScrollDown()
		case "<Up>":
			toggle.ScrollUp()
		case "q", "<C-c>":
			return
		case "<Resize>":
			setRect()
		}
		termui.Render(toggle)
	}
}
