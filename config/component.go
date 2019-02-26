package config

import (
	"github.com/sqshq/sampler/component/asciibox"
	"github.com/sqshq/sampler/data"
	ui "github.com/sqshq/termui"
)

type ComponentConfig struct {
	Title         string   `yaml:"title"`
	RefreshRateMs *int     `yaml:"refresh-rate-ms,omitempty"`
	Position      Position `yaml:"position"`
	Size          Size     `yaml:"size"`
}

type GaugeConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int              `yaml:"scale,omitempty"`
	Color           *ui.Color         `yaml:"color,omitempty"`
	Values          map[string]string `yaml:"values"`
	Items           []data.Item
}

type BarChartConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int        `yaml:"scale,omitempty"`
	Items           []data.Item `yaml:"items"`
}

type AsciiBoxConfig struct {
	ComponentConfig `yaml:",inline"`
	data.Item       `yaml:",inline"`
	Font            *asciibox.AsciiFont `yaml:"font,omitempty"`
}

type RunChartConfig struct {
	ComponentConfig `yaml:",inline"`
	Legend          *LegendConfig `yaml:"legend,omitempty"`
	Scale           *int          `yaml:"scale,omitempty"`
	Items           []data.Item   `yaml:"items"`
}

type LegendConfig struct {
	Enabled bool `yaml:"enabled"`
	Details bool `yaml:"details"`
}

type Position struct {
	X int `yaml:"w"`
	Y int `yaml:"h"`
}

type Size struct {
	X int `yaml:"w"`
	Y int `yaml:"h"`
}

type ComponentType rune

const (
	TypeRunChart ComponentType = 0
	TypeBarChart ComponentType = 1
	TypeTextBox  ComponentType = 2
	TypeAsciiBox ComponentType = 3
	TypeGauge    ComponentType = 4
)

type ComponentSettings struct {
	Type     ComponentType
	Title    string
	Size     Size
	Position Position
}
