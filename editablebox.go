package widgets

import (
	"image"

	. "github.com/gizak/termui/v3"
)

type EditableBox struct {
	*Block

	Text      string
	TextStyle Style
}

func NewEditableBox() *EditableBox {
	return &EditableBox{
		Block:     NewBlock(),
		TextStyle: NewStyle(Theme.Default.Fg, Theme.Default.Bg),
	}
}

func (self *EditableBox) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	cells := BuildCellWithXArray(TrimCells(ParseStyles(self.Text, self.TextStyle), self.Inner.Dx()))
	for _, cell := range cells {
		buf.SetCell(cell.Cell, image.Pt(cell.X, 0).Add(self.Inner.Min))
	}
}
