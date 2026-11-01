package config

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/console"
	"image"
)

type ComponentType rune

const (
	TypeRunChart  ComponentType = 0
	TypeBarChart  ComponentType = 1
	TypeSparkLine ComponentType = 2
	TypeTextBox   ComponentType = 3
	TypeAsciiBox  ComponentType = 4
	TypeGauge     ComponentType = 5
)

type ComponentConfig struct {
	Title     string          `yaml:"title"`
	Position  [][]int         `yaml:"position,flow"`
	RateMs    *int            `yaml:"rate-ms,omitempty"`
	Triggers  []TriggerConfig `yaml:"triggers,omitempty"`
	Type      ComponentType   `yaml:",omitempty"`
	ShowTitle *bool           `yaml:"show-title,omitempty"`
}

func (c *ComponentConfig) GetLocation() Location {
	return Location{X: c.Position[0][0], Y: c.Position[0][1]}
}

func (c *ComponentConfig) GetSize() Size {
	return Size{X: c.Position[1][0], Y: c.Position[1][1]}
}

func (c *ComponentConfig) GetRectangle() image.Rectangle {
	if c.Position == nil || len(c.Position) == 0 {
		return image.ZR
	}
	return image.Rect(
		c.Position[0][0],
		c.Position[0][1],
		c.Position[0][0]+c.Position[1][0],
		c.Position[0][1]+c.Position[1][1])
}

type TriggerConfig struct {
	Title     string         `yaml:"title"`
	Condition string         `yaml:"condition"`
	Actions   *ActionsConfig `yaml:"actions,omitempty"`
	Label     string         `yaml:"-"`
}

type ActionsConfig struct {
	TerminalBell *bool     `yaml:"terminal-bell,omitempty"`
	Sound        *bool     `yaml:"sound,omitempty"`
	Visual       *bool     `yaml:"visual,omitempty"`
	Script       *string   `yaml:"script,omitempty"`
	Color        *ui.Color `yaml:"color,omitempty"`
}

type GaugeConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int              `yaml:"scale,omitempty"`
	Color           *ui.Color         `yaml:"color,omitempty"`
	PercentOnly     *bool             `yaml:"percent-only,omitempty"`
	Cur             Item              `yaml:"cur,omitempty"`
	Max             Item              `yaml:"max,omitempty"`
	Min             Item              `yaml:"min,omitempty"`
	Items           []GaugeItemConfig `yaml:"items,omitempty"`
	Border          *bool             `yaml:"border,omitempty"`
}

type GaugeItemConfig struct {
	Item        `yaml:",inline"`
	Scale       *int            `yaml:"scale,omitempty"`
	PercentOnly *bool           `yaml:"percent-only,omitempty"`
	Cur         *Item           `yaml:"cur,omitempty"`
	Min         *Item           `yaml:"min,omitempty"`
	Max         *Item           `yaml:"max,omitempty"`
	Triggers    []TriggerConfig `yaml:"triggers,omitempty"`
}

func (it *GaugeItemConfig) CurItem() Item {
	if it.Cur == nil {
		return it.Item
	}
	cur := *it.Cur
	cur.Label = it.Label
	cur.Color = it.Color
	return cur
}

const (
	GaugeCurLabel       = "cur"
	GaugeMinLabel       = "min"
	GaugeMaxLabel       = "max"
	GaugeLabelSeparator = ":"
)

func GaugeItemMinLabel(label string) string {
	return label + GaugeLabelSeparator + GaugeMinLabel
}

func GaugeItemMaxLabel(label string) string {
	return label + GaugeLabelSeparator + GaugeMaxLabel
}

func (g *GaugeConfig) IsMulti() bool {
	return len(g.Items) > 0
}

func (g *GaugeConfig) DataItems() []Item {
	if !g.IsMulti() {
		return []Item{g.Cur, g.Min, g.Max}
	}
	items := make([]Item, 0)
	if g.Min.SampleScript != nil {
		items = append(items, withLabel(g.Min, GaugeMinLabel))
	}
	if g.Max.SampleScript != nil {
		items = append(items, withLabel(g.Max, GaugeMaxLabel))
	}
	for _, it := range g.Items {
		items = append(items, withLabel(it.CurItem(), *it.Label))
		if it.Min != nil {
			items = append(items, withLabel(*it.Min, GaugeItemMinLabel(*it.Label)))
		}
		if it.Max != nil {
			items = append(items, withLabel(*it.Max, GaugeItemMaxLabel(*it.Label)))
		}
	}
	return items
}

func (g *GaugeConfig) AllTriggers() []TriggerConfig {
	triggers := append([]TriggerConfig{}, g.Triggers...)
	for _, it := range g.Items {
		for _, t := range it.Triggers {
			t.Label = *it.Label
			triggers = append(triggers, t)
		}
	}
	return triggers
}

func withLabel(i Item, label string) Item {
	i.Label = &label
	if i.Pty == nil {
		pty := false
		i.Pty = &pty
	}
	return i
}

type SparkLineConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int        `yaml:"scale,omitempty"`
	Item            Item        `yaml:",inline"`
	Gradient        *[]ui.Color `yaml:",omitempty"`
	Min             *float64    `yaml:"min,omitempty"`
	Max             *float64    `yaml:"max,omitempty"`
	Border          *bool       `yaml:"border,omitempty"`
}

type BarChartConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int   `yaml:"scale,omitempty"`
	Items           []Item `yaml:"items"`
	Border          *bool  `yaml:"border,omitempty"`
}

type AsciiBoxConfig struct {
	ComponentConfig `yaml:",inline"`
	Item            `yaml:",inline"`
	Border          *bool              `yaml:"border,omitempty"`
	Font            *console.AsciiFont `yaml:"font,omitempty"`
}

type TextBoxConfig struct {
	ComponentConfig `yaml:",inline"`
	Item            `yaml:",inline"`
	Border          *bool `yaml:"border,omitempty"`
}

type RunChartConfig struct {
	ComponentConfig `yaml:",inline"`
	Legend          *LegendConfig `yaml:"legend,omitempty"`
	Scale           *int          `yaml:"scale,omitempty"`
	Items           []Item        `yaml:"items"`
	Border          *bool         `yaml:"border,omitempty"`
}

type LegendConfig struct {
	Enabled bool `yaml:"enabled"`
	Details bool `yaml:"details"`
}

type Item struct {
	Label               *string   `yaml:"label,omitempty"`
	Color               *ui.Color `yaml:"color,omitempty"`
	Pty                 *bool     `yaml:"pty,omitempty"`
	InitScript          *string   `yaml:"init,omitempty"`
	MultiStepInitScript *[]string `yaml:"multistep-init,omitempty"`
	SampleScript        *string   `yaml:"sample,omitempty"`
	TransformScript     *string   `yaml:"transform,omitempty"`
}

type Location struct {
	X int
	Y int
}

type Size struct {
	X int
	Y int
}

type ComponentSettings struct {
	Type     ComponentType
	Title    string
	Size     Size
	Location Location
}

func getPosition(location Location, size Size) [][]int {
	return [][]int{
		{location.X, location.Y},
		{size.X, size.Y},
	}
}

func (c *ComponentConfig) IsTitleVisible() bool {
	return c.ShowTitle == nil || *c.ShowTitle
}
