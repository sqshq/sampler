package config

import (
	"fmt"
	"github.com/sqshq/sampler/console"
)

func (c *Config) validate() {

	var components []ComponentConfig

	for _, c := range c.RunCharts {
		components = append(components, c.ComponentConfig)
		validateLabelsUniqueness(c.Title, c.Items)
	}
	for _, c := range c.BarCharts {
		components = append(components, c.ComponentConfig)
		validateLabelsUniqueness(c.Title, c.Items)
	}
	for _, c := range c.SparkLines {
		components = append(components, c.ComponentConfig)
	}
	for _, c := range c.Gauges {
		components = append(components, c.ComponentConfig)
	}
	for _, c := range c.AsciiBoxes {
		components = append(components, c.ComponentConfig)
	}
	for _, c := range c.TextBoxes {
		components = append(components, c.ComponentConfig)
	}

	validateTitlesUniqueness(components)
}

func validateLabelsUniqueness(title string, items []Item) {
	labels := make(map[string]bool)
	for _, c := range items {
		label := *c.Label
		if _, contains := labels[label]; contains {
			console.Exit(fmt.Sprintf("Config validation error: item labels should be unique. Please rename '%s' in '%s'", label, title))
		}
		labels[label] = true
	}
}

func validateTitlesUniqueness(components []ComponentConfig) {
	titles := make(map[string]bool)
	for _, c := range components {
		if _, contains := titles[c.Title]; contains {
			console.Exit(fmt.Sprintf("Config validation error: component titles should be unique. Please rename '%s'", c.Title))
		}
		titles[c.Title] = true
	}
}
