package layout

import (
	. "github.com/sqshq/termui"
)

type Layout struct {
	Block
	items []Item
}

const (
	columnsCount = 30
	rowsCount    = 30
)

func NewLayout(width, height int) *Layout {

	block := *NewBlock()
	block.SetRect(0, 0, width, height)

	return &Layout{
		Block: block,
		items: make([]Item, 0),
	}
}

func (self *Layout) AddItem(drawable Drawable, position Position, size Size) {
	self.items = append(self.items, Item{drawable, position, size})
}

func (self *Layout) ChangeDimensions(width, height int) {
	self.SetRect(0, 0, width, height)
}

func (self *Layout) Draw(buf *Buffer) {

	columnWidth := float64(self.GetRect().Dx()) / columnsCount
	rowHeight := float64(self.GetRect().Dy()) / rowsCount

	for _, item := range self.items {

		x1 := float64(item.Position.X) * columnWidth
		y1 := float64(item.Position.Y) * rowHeight
		x2 := x1 + float64(item.Size.X)*columnWidth
		y2 := y1 + float64(item.Size.Y)*rowHeight

		item.Data.SetRect(int(x1), int(y1), int(x2), int(y2))
		item.Data.Draw(buf)
	}
}
