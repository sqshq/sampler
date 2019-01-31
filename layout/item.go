package layout

import (
	. "github.com/sqshq/termui"
)

type Item struct {
	Data     Drawable
	Position Position
	Size     Size
}

type Position struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type Size struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

func (self *Item) MoveItem(x, y int) {
	self.Position.X += x
	self.Position.Y += y
}

func (self *Item) ResizeItem(x, y int) {
	self.Size.X += x
	self.Size.Y += y
}
