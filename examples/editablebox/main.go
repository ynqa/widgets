package main

import (
	"log"

	"github.com/gizak/termui/v3"

	"github.com/ynqa/widgets/pkg/widgets"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatal(err)
	}
	defer termui.Close()

	box := widgets.NewEditableBox()

	event := termui.PollEvents()
	setRect := func() {
		width, _ := termui.TerminalDimensions()
		box.SetRect(0, 1, width, 4)
	}
	setRect()
	termui.Render(box)

	for e := range event {
		switch {
		case len(e.ID) == 1:
			box.Text += e.ID
		case e.ID == "<C-c>":
			return
		case e.ID == "<Enter>":
			box.Text = ""
		case e.ID == "<Resize>":
			setRect()
		}
		termui.Render(box)
	}
}
