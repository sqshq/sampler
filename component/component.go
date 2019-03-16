package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/data"
)

type Component struct {
	ui.Drawable
	*data.Consumer
	Type          config.ComponentType
	Title         string
	Position      config.Position
	Size          config.Size
	RefreshRateMs int
}

func NewComponent(dbl ui.Drawable, cmr *data.Consumer, cfg config.ComponentConfig, ct config.ComponentType) *Component {
	return &Component{
		Drawable:      dbl,
		Consumer:      cmr,
		Type:          ct,
		Title:         cfg.Title,
		Position:      cfg.Position,
		Size:          cfg.Size,
		RefreshRateMs: *cfg.RefreshRateMs,
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
