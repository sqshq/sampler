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
	Title    string          `yaml:"title"`
	Position [][]int         `yaml:"position,flow"`
	RateMs   *int            `yaml:"rate-ms,omitempty"`
	Triggers []TriggerConfig `yaml:"triggers,omitempty"`
	Type     ComponentType   `yaml:",omitempty"`
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
}

type ActionsConfig struct {
	TerminalBell *bool   `yaml:"terminal-bell,omitempty"`
	Sound        *bool   `yaml:"sound,omitempty"`
	Visual       *bool   `yaml:"visual,omitempty"`
	Script       *string `yaml:"script,omitempty"`
}

type GaugeConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int      `yaml:"scale,omitempty"`
	Color           *ui.Color `yaml:"color,omitempty"`
	PercentOnly     *bool     `yaml:"percent-only,omitempty"`
	Cur             Item      `yaml:"cur"`
	Max             Item      `yaml:"max"`
	Min             Item      `yaml:"min"`
}

type SparkLineConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int        `yaml:"scale,omitempty"`
	Item            Item        `yaml:",inline"`
	Gradient        *[]ui.Color `yaml:",omitempty"`
}

type BarChartConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int   `yaml:"scale,omitempty"`
	Items           []Item `yaml:"items"`
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
	SampleScript        *string   `yaml:"sample"`
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
