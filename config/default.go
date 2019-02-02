package config

import (
	"github.com/sqshq/vcmd/console"
)

const (
	defaultRefreshRateMs = 300
	defaultTimeScaleSec  = 1
	defaultTheme         = console.ThemeDark
)

func (self *Config) setDefaultValues() {

	if len(self.Theme) == 0 {
		self.Theme = defaultTheme
	}

	for i, chart := range self.RunCharts {
		if chart.RefreshRateMs == 0 {
			chart.RefreshRateMs = defaultRefreshRateMs
		}
		if chart.TimeScaleSec == 0 {
			chart.TimeScaleSec = defaultTimeScaleSec
		}
		self.RunCharts[i] = chart
	}
}

func (config *Config) setDefaultLayout() {

}

func (config *Config) setDefaultColors() {

	palette := console.GetPalette(config.Theme)

	for i, chart := range config.RunCharts {
		for j, item := range chart.Items {
			if item.Color == 0 {
				item.Color = palette.Colors[i+j]
				chart.Items[j] = item
			}
		}
	}
}
