package widgets

import (
	. "github.com/sqshq/termui"
)

type Component struct {
	Drawable Drawable
	Title    string
	Position Position
	Size     Size
	Type     ComponentType
}

type ComponentType rune

const (
	TypeRunChart ComponentType = 0
	TypeBarChart ComponentType = 1
)

type Position struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type Size struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

func (self *Component) Move(x, y int) {
	self.Position.X += x
	self.Position.Y += y
}

func (self *Component) Resize(x, y int) {
	self.Size.X += x
	self.Size.Y += y
}
