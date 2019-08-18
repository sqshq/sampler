package config

import (
	"github.com/sqshq/sampler/console"
)

const (
	defaultRateMs = 1000
	defaultScale  = 1
	defaultTheme  = console.ThemeDark
)

func (c *Config) setDefaults() {
	c.setDefaultValues()
	c.setDefaultItemSettings()
	c.setDefaultArrangement()
}

func (c *Config) setDefaultValues() {

	if c.Theme == nil {
		t := defaultTheme
		c.Theme = &t
	}

	for i, chart := range c.RunCharts {

		setDefaultTriggersValues(chart.Triggers)
		chart.ComponentConfig.Type = TypeRunChart

		if chart.RateMs == nil {
			r := defaultRateMs
			chart.RateMs = &r
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

	for i, line := range c.SparkLines {

		setDefaultTriggersValues(line.Triggers)
		line.ComponentConfig.Type = TypeSparkLine
		line.Item.Label = &line.Title

		if line.RateMs == nil {
			r := defaultRateMs
			line.RateMs = &r
		}
		if line.Scale == nil {
			p := defaultScale
			line.Scale = &p
		}
		c.SparkLines[i] = line
	}

	for i, chart := range c.BarCharts {

		setDefaultTriggersValues(chart.Triggers)
		chart.ComponentConfig.Type = TypeBarChart

		if chart.RateMs == nil {
			r := defaultRateMs
			chart.RateMs = &r
		}
		if chart.Scale == nil {
			p := defaultScale
			chart.Scale = &p
		}
		c.BarCharts[i] = chart
	}

	for i, g := range c.Gauges {

		setDefaultTriggersValues(g.Triggers)
		g.ComponentConfig.Type = TypeGauge

		if g.RateMs == nil {
			r := defaultRateMs
			g.RateMs = &r
		}
		if g.Scale == nil {
			p := defaultScale
			g.Scale = &p
		}

		cur := "cur"
		max := "max"
		min := "min"

		g.Cur.Label = &cur
		g.Max.Label = &max
		g.Min.Label = &min

		c.Gauges[i] = g
	}

	for i, box := range c.AsciiBoxes {

		setDefaultTriggersValues(box.Triggers)
		box.ComponentConfig.Type = TypeAsciiBox

		if box.RateMs == nil {
			r := defaultRateMs
			box.RateMs = &r
		}
		if box.Label == nil {
			label := box.Title
			box.Label = &label
		}
		if box.Font == nil {
			font := console.AsciiFont2D
			box.Font = &font
		}
		if box.Border == nil {
			border := true
			box.Border = &border
		}
		c.AsciiBoxes[i] = box
	}

	for i, box := range c.TextBoxes {

		setDefaultTriggersValues(box.Triggers)
		box.ComponentConfig.Type = TypeTextBox

		if box.RateMs == nil {
			r := defaultRateMs
			box.RateMs = &r
		}
		if box.Label == nil {
			label := box.Title
			box.Label = &label
		}
		if box.Border == nil {
			border := true
			box.Border = &border
		}

		c.TextBoxes[i] = box
	}
}

func setDefaultTriggersValues(triggers []TriggerConfig) {

	defaultTerminalBell := false
	defaultSound := false
	defaultVisual := false

	for i, trigger := range triggers {

		if trigger.Actions == nil {
			trigger.Actions = &ActionsConfig{TerminalBell: &defaultTerminalBell, Sound: &defaultSound, Visual: &defaultVisual, Script: nil}
		} else {
			if trigger.Actions.TerminalBell == nil {
				trigger.Actions.TerminalBell = &defaultTerminalBell
			}
			if trigger.Actions.Sound == nil {
				trigger.Actions.Sound = &defaultSound
			}
			if trigger.Actions.Visual == nil {
				trigger.Actions.Visual = &defaultVisual
			}
		}

		triggers[i] = trigger
	}
}

func (c *Config) setDefaultItemSettings() {

	palette := console.GetPalette(*c.Theme)
	colorsCount := len(palette.ContentColors)
	defaultPty := false
	defaultPercentOnly := false

	for _, ch := range c.RunCharts {
		for j, item := range ch.Items {
			if item.Color == nil {
				item.Color = &palette.ContentColors[j%colorsCount]
			}
			if item.Pty == nil {
				item.Pty = &defaultPty
			}
			ch.Items[j] = item
		}
	}

	for _, b := range c.BarCharts {
		for j, item := range b.Items {
			if item.Color == nil {
				item.Color = &palette.ContentColors[j%colorsCount]
			}
			if item.Pty == nil {
				item.Pty = &defaultPty
			}
			b.Items[j] = item
		}
	}

	for i, s := range c.SparkLines {
		s.Gradient = &palette.GradientColors[i%(len(palette.GradientColors))]
		if s.Item.Pty == nil {
			s.Item.Pty = &defaultPty
		}
		c.SparkLines[i] = s
	}

	for i, g := range c.Gauges {
		if g.Min.Pty == nil {
			g.Min.Pty = &defaultPty
		}
		if g.Max.Pty == nil {
			g.Max.Pty = &defaultPty
		}
		if g.Cur.Pty == nil {
			g.Cur.Pty = &defaultPty
		}
		if g.Color == nil {
			g.Color = &palette.ContentColors[i%colorsCount]
		}
		if g.PercentOnly == nil {
			g.PercentOnly = &defaultPercentOnly
		}
		c.Gauges[i] = g
	}

	for i, a := range c.AsciiBoxes {
		if a.Item.Pty == nil {
			a.Item.Pty = &defaultPty
		}
		c.AsciiBoxes[i] = a
	}

	for i, t := range c.TextBoxes {
		if t.Item.Pty == nil {
			t.Item.Pty = &defaultPty
		}
		c.TextBoxes[i] = t
	}
}
