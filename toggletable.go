package widgets

import (
	"image"
	"strings"
	"unicode/utf8"

	. "github.com/gizak/termui/v3"

	"github.com/ynqa/widgets/node"
)

const (
	widthFromLeftBorder = 2
)

type ToggleTable struct {
	*Block

	Headers []string
	Widths  []int

	Node         *node.Node
	HeaderStyle  Style
	CursorStyle  Style
	RowStyle     Style
	FoldSymbol   rune
	UnfoldSymbol rune

	SelectedRow    int
	drawInitialRow int
}

func NewToggleTable() *ToggleTable {
	return &ToggleTable{
		Block:        NewBlock(),
		HeaderStyle:  NewStyle(Theme.Default.Fg, Theme.Default.Bg, ModifierBold),
		CursorStyle:  NewStyle(ColorBlack, ColorYellow),
		RowStyle:     NewStyle(Theme.Default.Fg),
		FoldSymbol:   '▸',
		UnfoldSymbol: '▾',
	}
}

func (self *ToggleTable) rowPrefix(n *node.Node) string {
	arrow := string(self.FoldSymbol)
	if n.ChildVisible() {
		arrow = string(self.UnfoldSymbol)
	}
	if n.IsLeaf() {
		arrow = strings.Repeat(" ", utf8.RuneCountInString(arrow))
	}
	return strings.Repeat(" ", n.Depth()) + arrow + " "
}

func (self *ToggleTable) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	// coordinates: space from border + header + initial row = 3
	if self.Inner.Dy() >= 3 {
		// store start positions for each column
		var (
			colPos []int
			cur    int = widthFromLeftBorder
		)
		for _, w := range self.Widths {
			colPos = append(colPos, cur)
			cur += w
		}

		// draw headers
		for i, h := range self.Headers {
			// replace to '…' if the field is over
			h := TrimString(h, self.Widths[i]-widthFromLeftBorder)
			buf.SetString(
				h,
				self.HeaderStyle,
				image.Pt(
					self.Inner.Min.X+colPos[i],
					// coordinates: space from border = 1
					self.Inner.Min.Y+1),
			)
		}

		if self.SelectedRow < self.drawInitialRow {
			self.drawInitialRow = self.SelectedRow
		} else if self.SelectedRow >= self.drawInitialRow+self.Inner.Dy()-2 {
			// coordinates: space from border + header = 2
			self.drawInitialRow += self.Inner.Dy() - 2
		}

		nodes := self.Node.Flatten()

		// draw rows
		for idx := self.drawInitialRow; idx >= 0 && idx < len(nodes) && idx < self.drawInitialRow+self.Inner.Dy()-2; idx++ {
			node := nodes[idx]
			// coordinates: space from border + header = 2
			y := self.Inner.Min.Y + idx - self.drawInitialRow + 2
			if idx == self.SelectedRow {
				buf.SetString(
					strings.Repeat(" ", self.Inner.Dx()),
					self.CursorStyle,
					image.Pt(self.Inner.Min.X, y),
				)
				self.setselected(idx)
			}
			style := self.RowStyle
			if idx == self.SelectedRow {
				style = self.CursorStyle
			}
			for i, width := range self.Widths {
				row := node.Row()[i]
				if i == 0 {
					row = self.rowPrefix(node) + node.Row()[i]
				}
				r := TrimString(row, width-widthFromLeftBorder)
				buf.SetString(
					r,
					style,
					image.Pt(self.Inner.Min.X+colPos[i], y),
				)
			}
		}
	}
}

func (self *ToggleTable) setselected(idx int) {
	rows := self.Node.Flatten()
	self.SelectedRow = idx
	max := len(rows) - 1
	if max >= 0 && self.SelectedRow < 0 {
		self.SelectedRow = max
	} else if max >= 0 && self.SelectedRow > max {
		self.SelectedRow = 0
	}
}

func (self *ToggleTable) ScrollUp() {
	self.setselected(self.SelectedRow - 1)
}

func (self *ToggleTable) ScrollDown() {
	self.setselected(self.SelectedRow + 1)
}
