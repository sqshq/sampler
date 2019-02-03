package widgets

import (
	. "github.com/sqshq/termui"
)

type Layout struct {
	Block
	components []Component
}

const (
	columnsCount = 30
	rowsCount    = 30
)

func NewLayout(width, height int) *Layout {

	block := *NewBlock()
	block.SetRect(0, 0, width, height)

	return &Layout{
		Block:      block,
		components: make([]Component, 0),
	}
}

func (self *Layout) AddComponent(drawable Drawable, position Position, size Size, Type ComponentType) {
	self.components = append(self.components, Component{drawable, position, size, Type})
}

func (self *Layout) GetComponents(Type ComponentType) []Drawable {

	var components []Drawable

	for _, component := range self.components {
		if component.Type == Type {
			components = append(components, component.Drawable)
		}
	}

	return components
}

func (self *Layout) ChangeDimensions(width, height int) {
	self.SetRect(0, 0, width, height)
}

func (self *Layout) Draw(buffer *Buffer) {

	columnWidth := float64(self.GetRect().Dx()) / columnsCount
	rowHeight := float64(self.GetRect().Dy()) / rowsCount

	for _, component := range self.components {

		x1 := float64(component.Position.X) * columnWidth
		y1 := float64(component.Position.Y) * rowHeight
		x2 := x1 + float64(component.Size.X)*columnWidth
		y2 := y1 + float64(component.Size.Y)*rowHeight

		component.Drawable.SetRect(int(x1), int(y1), int(x2), int(y2))
		component.Drawable.Draw(buffer)
	}
}
