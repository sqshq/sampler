package config

import (
	"github.com/sqshq/sampler/component/asciibox"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
)

const (
	defaultRefreshRateMs = 1000
	defaultScale         = 1
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
		if chart.Scale == nil {
			p := defaultScale
			chart.Scale = &p
		}
		if chart.Legend == nil {
			chart.Legend = &LegendConfig{true, true}
			c.RunCharts[i] = chart
		}
		c.RunCharts[i] = chart
	}

	for i, chart := range c.BarCharts {
		if chart.RefreshRateMs == nil {
			r := defaultRefreshRateMs
			chart.RefreshRateMs = &r
		}
		if chart.Scale == nil {
			p := defaultScale
			chart.Scale = &p
		}
		c.BarCharts[i] = chart
	}

	for i, g := range c.Gauges {
		if g.RefreshRateMs == nil {
			r := defaultRefreshRateMs
			g.RefreshRateMs = &r
		}
		if g.Scale == nil {
			p := defaultScale
			g.Scale = &p
		}
		var items []data.Item
		for label, script := range g.Values {
			l := label
			items = append(items, data.Item{Label: &l, Script: script})
		}
		g.Items = items
		c.Gauges[i] = g
	}

	for i, box := range c.AsciiBoxes {
		if box.RefreshRateMs == nil {
			r := defaultRefreshRateMs
			box.RefreshRateMs = &r
		}
		if box.Label == nil {
			label := string(i)
			box.Label = &label
		}
		if box.Font == nil {
			font := asciibox.AsciiFontFlat
			box.Font = &font
		}
		if box.Color == nil {
			color := console.ColorWhite
			box.Color = &color
		}
		c.AsciiBoxes[i] = box
	}
}

func (c *Config) setDefaultLayout() {
	// TODO auto-arrange components
}

func (c *Config) setDefaultColors() {

	palette := console.GetPalette(*c.Theme)
	colorsCount := len(palette.Colors)

	for _, chart := range c.RunCharts {
		for j, item := range chart.Items {
			if item.Color == nil {
				item.Color = &palette.Colors[j%colorsCount]
				chart.Items[j] = item
			}
		}
	}

	for _, chart := range c.BarCharts {
		for j, item := range chart.Items {
			if item.Color == nil {
				item.Color = &palette.Colors[j%colorsCount]
				chart.Items[j] = item
			}
		}
	}

	for i, g := range c.Gauges {
		if g.Color == nil {
			g.Color = &palette.Colors[i%colorsCount]
			c.Gauges[i] = g
		}
	}
}
