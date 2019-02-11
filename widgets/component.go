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

func (c *Component) Move(x, y int) {
	c.Position.X += x
	c.Position.Y += y
	c.normalize()
}

func (c *Component) Resize(x, y int) {
	c.Size.X += x
	c.Size.Y += y
	c.normalize()
}

func (c *Component) normalize() {
	if c.Size.X < 0 {
		c.Size.X = 0
	}
	if c.Size.Y < 0 {
		c.Size.Y = 0
	}
	if c.Position.X < 0 {
		c.Position.X = 0
	}
	if c.Position.Y < 0 {
		c.Position.Y = 0
	}
}
