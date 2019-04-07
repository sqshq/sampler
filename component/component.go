package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/data"
)

type Component struct {
	ui.Drawable
	*data.Consumer
	Type     config.ComponentType
	Title    string
	Location config.Location
	Size     config.Size
	RateMs   int
}

func NewComponent(dbl ui.Drawable, cmr *data.Consumer, cfg config.ComponentConfig) *Component {
	return &Component{
		Drawable: dbl,
		Consumer: cmr,
		Type:     cfg.Type,
		Title:    cfg.Title,
		Location: cfg.GetLocation(),
		Size:     cfg.GetSize(),
		RateMs:   *cfg.RateMs,
	}
}

func (c *Component) Move(x, y int) {
	c.Location.X += x
	c.Location.Y += y
	c.normalize()
}

func (c *Component) Resize(x, y int) {
	c.Size.X += x
	c.Size.Y += y
	c.normalize()
}

func (c *Component) normalize() {
	if c.Location.X < 0 {
		c.Location.X = 0
	}
	if c.Location.Y < 0 {
		c.Location.Y = 0
	}
}
