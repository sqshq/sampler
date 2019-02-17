package config

import (
	"github.com/sqshq/sampler/console"
)

const (
	defaultRefreshRateMs = 300
	defaultPrecision     = 1
	defaultTheme         = console.ThemeDark
)

func (c *Config) setDefaults() {
	c.setDefaultValues()
	c.setDefaultColors()
	c.setDefaultLayout()
}

func (c *Config) setDefaultValues() {

	if c.Theme == nil {
		t := defaultTheme
		c.Theme = &t
	}

	for i, chart := range c.RunCharts {
		if chart.RefreshRateMs == nil {
			r := defaultRefreshRateMs
			chart.RefreshRateMs = &r
		}
		if chart.Precision == nil {
			p := defaultPrecision
			chart.Precision = &p
		}
		if chart.Legend == nil {
			chart.Legend = &LegendConfig{true, true}
			c.RunCharts[i] = chart
		}
		c.RunCharts[i] = chart
	}
}

func (c *Config) setDefaultLayout() {

}

func (c *Config) setDefaultColors() {

	palette := console.GetPalette(*c.Theme)

	for i, chart := range c.RunCharts {
		for j, item := range chart.Items {
			if item.Color == nil {
				item.Color = &palette.Colors[i+j] // TODO handle out of range case
				chart.Items[j] = item
			}
		}
	}
}
