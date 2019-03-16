package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/data"
)

type Component struct {
	ui.Block
	data.Consumer
	*Alerter
	Type          config.ComponentType
	Title         string
	Position      config.Position
	Size          config.Size
	RefreshRateMs int
}

func NewComponent(c config.ComponentConfig, t config.ComponentType) *Component {

	consumer := data.NewConsumer()
	block := *ui.NewBlock()
	block.Title = c.Title

	return &Component{
		Block:         block,
		Consumer:      consumer,
		Alerter:       NewAlerter(consumer.AlertChannel),
		Type:          t,
		Title:         c.Title,
		Position:      c.Position,
		Size:          c.Size,
		RefreshRateMs: *c.RefreshRateMs,
	}
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
	if c.Position.X < 0 {
		c.Position.X = 0
	}
	if c.Position.Y < 0 {
		c.Position.Y = 0
	}
}
