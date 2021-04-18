package table

import (
	"image"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/gizak/termui/v3"
	"github.com/ynqa/widgets/pkg/table/node"
)

const (
	widthFromLeftBorder = 2
)

type Table struct {
	opts option

	headers []Header
	node    *node.Node

	selectedRow    int
	drawInitialRow int

	sync.Mutex
}

type Header struct {
	Header string
	Width  int
}

type option struct {
	block            *termui.Block
	headerStyle      termui.Style
	cursoredRowStyle termui.Style
	defaultRowStyle  termui.Style
	foldedSymbol     rune
	unfoldedSymbol   rune
}

type Option func(*option)

func Block(block *termui.Block) Option {
	return Option(func(o *option) {
		o.block = block
	})
}

func HeaderStyle(style termui.Style) Option {
	return Option(func(o *option) {
		o.headerStyle = style
	})
}

func CursoredRowStyle(style termui.Style) Option {
	return Option(func(o *option) {
		o.cursoredRowStyle = style
	})
}

func DefaultRowStyle(style termui.Style) Option {
	return Option(func(o *option) {
		o.defaultRowStyle = style
	})
}

func FoldedSymbol(symbol rune) Option {
	return Option(func(o *option) {
		o.foldedSymbol = symbol
	})
}

func UnfoldedSymbol(symbol rune) Option {
	return Option(func(o *option) {
		o.unfoldedSymbol = symbol
	})
}

func New(headers []Header, opts ...Option) *Table {
	option := option{
		block:            termui.NewBlock(),
		headerStyle:      termui.NewStyle(termui.Theme.Default.Fg, termui.Theme.Default.Bg, termui.ModifierBold),
		cursoredRowStyle: termui.NewStyle(termui.ColorBlack, termui.ColorYellow),
		defaultRowStyle:  termui.NewStyle(termui.Theme.Default.Fg),
		foldedSymbol:     '▸',
		unfoldedSymbol:   '▾',
	}
	for _, fn := range opts {
		fn(&option)
	}

	return &Table{
		opts:    option,
		headers: headers,
		node:    node.Root(),
	}
}

func (self *Table) GetRect() image.Rectangle {
	return self.opts.block.Rectangle
}

func (self *Table) SetRect(x1, y1, x2, y2 int) {
	self.opts.block.SetRect(x1, y1, x2, y2)
}

func (self *Table) SelectedRow() int {
	return self.selectedRow
}

func (self *Table) SetNode(node *node.Node) {
	self.node = node
}

func (self *Table) rowPrefix(n *node.Node) string {
	arrow := string(self.opts.foldedSymbol)
	if n.ChildVisible() {
		arrow = string(self.opts.unfoldedSymbol)
	}
	if n.IsLeaf() {
		arrow = strings.Repeat(" ", utf8.RuneCountInString(arrow))
	}
	return strings.Repeat(" ", n.Depth()) + arrow + " "
}

func (self *Table) Draw(buf *termui.Buffer) {
	self.opts.block.Draw(buf)

	if self.opts.block.Inner.Dy() >= 3 {
		var (
			colPos []int
			cur    int = widthFromLeftBorder
		)
		for _, h := range self.headers {
			colPos = append(colPos, cur)
			cur += h.Width
		}

		for i, h := range self.headers {
			buf.SetString(
				termui.TrimString(h.Header, h.Width-widthFromLeftBorder),
				self.opts.headerStyle,
				image.Pt(
					self.opts.block.Inner.Min.X+colPos[i],
					self.opts.block.Inner.Min.Y+1,
				),
			)
		}

		if self.selectedRow < self.drawInitialRow {
			self.drawInitialRow = self.selectedRow
		} else if self.selectedRow >= self.drawInitialRow+self.opts.block.Inner.Dy()-2 {
			self.drawInitialRow += self.opts.block.Inner.Dy() - 2
		}

		nodes := self.node.Flatten()

		for idx := self.drawInitialRow; idx >= 0 && idx < len(nodes) && idx < self.drawInitialRow+self.opts.block.Inner.Dy()-2; idx++ {
			node := nodes[idx]
			y := self.opts.block.Inner.Min.Y + idx - self.drawInitialRow + 2
			if idx == self.selectedRow {
				buf.SetString(
					strings.Repeat(" ", self.opts.block.Inner.Dx()),
					self.opts.cursoredRowStyle,
					image.Pt(self.opts.block.Inner.Min.X, y),
				)
				self.setselected(idx)
			}
			style := self.opts.defaultRowStyle
			if idx == self.selectedRow {
				style = self.opts.cursoredRowStyle
			}
			for i, h := range self.headers {
				row := node.Row()[i]
				if i == 0 {
					row = self.rowPrefix(node) + node.Row()[i]
				}
				buf.SetString(
					termui.TrimString(row, h.Width-widthFromLeftBorder),
					style,
					image.Pt(self.opts.block.Inner.Min.X+colPos[i], y),
				)
			}
		}
	}
}

func (self *Table) setselected(idx int) {
	rows := self.node.Flatten()
	self.selectedRow = idx
	max := len(rows) - 1
	if max >= 0 && self.selectedRow < 0 {
		self.selectedRow = max
	} else if max >= 0 && self.selectedRow > max {
		self.selectedRow = 0
	}
}

func (self *Table) ScrollUp() {
	self.setselected(self.selectedRow - 1)
}

func (self *Table) ScrollDown() {
	self.setselected(self.selectedRow + 1)
}
