package widgets

import (
	. "github.com/sqshq/termui"
)

type Item struct {
	drawable    Drawable
	coordinates ItemCoordinates
}

type ItemCoordinates struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

type LayoutDimensions struct {
	width  int
	height int
}

type Layout struct {
	Block
	dimensions LayoutDimensions
	items      []Item
}

const (
	columnsCount = 12
	rowsCount    = 12
)

func NewLayout(width, height int) *Layout {

	b := *NewBlock()
	b.SetRect(0, 0, width, height)

	return &Layout{
		Block:      b,
		dimensions: LayoutDimensions{width, height},
		items:      make([]Item, 0),
	}
}

func (self *Layout) AddItem(drawable interface{}, x1, y1, x2, y2 int) {
	self.items = append(self.items, Item{
		drawable:    drawable.(Drawable),
		coordinates: ItemCoordinates{x1, y1, x2, y2},
	})
}

func (self *Layout) MoveItem(x, y int) {
	self.items[0].coordinates.x1 += x
	self.items[0].coordinates.y1 += y
	self.items[0].coordinates.x2 += x
	self.items[0].coordinates.y2 += y
}

func (self *Layout) ChangeDimensions(width, height int) {
	self.dimensions = LayoutDimensions{width, height}
	self.SetRect(0, 0, width, height)
}

func (self *Layout) Draw(buf *Buffer) {

	columnWidth := float64(self.dimensions.width) / columnsCount
	rowHeight := float64(self.dimensions.height) / rowsCount

	for _, item := range self.items {

		x1 := float64(item.coordinates.x1) * columnWidth
		y1 := float64(item.coordinates.y1) * rowHeight
		x2 := float64(item.coordinates.x2) * columnWidth
		y2 := float64(item.coordinates.y2) * rowHeight

		item.drawable.SetRect(int(x1), int(y1), int(x2), int(y2))
		item.drawable.Draw(buf)
	}
}
