package main

import (
	"log"

	"github.com/gizak/termui/v3"

	"github.com/ynqa/widgets/pkg/table"
	"github.com/ynqa/widgets/pkg/table/node"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatal(err)
	}
	defer termui.Close()

	headers := []table.Header{
		{
			Header: "index",
			Width:  10,
		},
		{
			Header: "item",
			Width:  10,
		},
	}
	block := termui.NewBlock()
	block.Title = "example"
	t := table.New(headers, table.Block(block))

	root := node.Root()
	root.Append(
		node.New("0", []string{"0", "aaa"}).Append(
			node.New("1", []string{"1", "bbb"}).Append(
				node.New("2", []string{"2", "ccc"}),
			),
			node.New("3", []string{"3", "ddd"}),
		),
	)
	t.SetNode(root)

	event := termui.PollEvents()
	setRect := func() {
		width, height := termui.TerminalDimensions()
		t.SetRect(0, 1, width, height-1)
	}
	setRect()
	termui.Render(t)

	for e := range event {
		switch e.ID {
		case "<Enter>":
			root.Toggle(t.SelectedRow())
		case "<Down>":
			t.ScrollDown()
		case "<Up>":
			t.ScrollUp()
		case "q", "<C-c>":
			return
		case "<Resize>":
			setRect()
		}
		termui.Render(t)
	}
}
