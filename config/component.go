package config

import (
	"github.com/sqshq/sampler/console"
	ui "github.com/sqshq/termui"
)

type ComponentType rune

const (
	TypeRunChart ComponentType = 0
	TypeBarChart ComponentType = 1
	TypeTextBox  ComponentType = 2
	TypeAsciiBox ComponentType = 3
	TypeGauge    ComponentType = 4
)

type ComponentConfig struct {
	Title         string          `yaml:"title"`
	RefreshRateMs *int            `yaml:"refresh-rate-ms,omitempty"`
	Position      Position        `yaml:"position"`
	Size          Size            `yaml:"size"`
	Triggers      []TriggerConfig `yaml:"triggers,omitempty"`
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
	Scale           *int              `yaml:"scale,omitempty"`
	Color           *ui.Color         `yaml:"color,omitempty"`
	Values          map[string]string `yaml:"values"`
	Items           []Item            `yaml:",omitempty"`
}

type BarChartConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int   `yaml:"scale,omitempty"`
	Items           []Item `yaml:"items"`
}

type AsciiBoxConfig struct {
	ComponentConfig `yaml:",inline"`
	Item            `yaml:",inline"`
	Font            *console.AsciiFont `yaml:"font,omitempty"`
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

type Position struct {
	X int `yaml:"w"`
	Y int `yaml:"h"`
}

type Size struct {
	X int `yaml:"w"`
	Y int `yaml:"h"`
}

type Item struct {
	Label  *string   `yaml:"label,omitempty"`
	Script string    `yaml:"value"`
	Color  *ui.Color `yaml:"color,omitempty"`
}

type ComponentSettings struct {
	Type     ComponentType
	Title    string
	Size     Size
	Position Position
}
