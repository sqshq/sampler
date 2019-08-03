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
		validateItemsScripts(c.Title, c.Items)
	}
	for _, c := range c.BarCharts {
		components = append(components, c.ComponentConfig)
		validateLabelsUniqueness(c.Title, c.Items)
		validateItemsScripts(c.Title, c.Items)
	}
	for _, c := range c.SparkLines {
		components = append(components, c.ComponentConfig)
		validateItemScripts(c.Title, c.Item)
	}
	for _, c := range c.Gauges {
		components = append(components, c.ComponentConfig)
		validateItemScripts(c.Title, c.Min)
		validateItemScripts(c.Title, c.Max)
		validateItemScripts(c.Title, c.Cur)
	}
	for _, c := range c.AsciiBoxes {
		components = append(components, c.ComponentConfig)
		validateItemScripts(c.Title, c.Item)
	}
	for _, c := range c.TextBoxes {
		components = append(components, c.ComponentConfig)
		validateItemScripts(c.Title, c.Item)
	}

	validateTitlesUniqueness(components)
}

func validateItemsScripts(title string, items []Item) {
	for _, i := range items {
		validateItemScripts(title, i)
	}
}

func validateItemScripts(title string, i Item) {
	if i.InitScript != nil && i.MultiStepInitScript != nil {
		console.Exit(fmt.Sprintf("Config validation error: both init and multistep-init scripts are not allowed for '%s'", title))
	}
	if i.SampleScript == nil {
		console.Exit(fmt.Sprintf("Config validation error: sample script should be specified for '%s'", title))
	}
}

func validateLabelsUniqueness(title string, items []Item) {
	labels := make(map[string]bool)
	for _, i := range items {
		label := *i.Label
		if _, contains := labels[label]; contains {
			console.Exit(fmt.Sprintf("Config validation error: item labels should be unique. Please rename '%s' for '%s'", label, title))
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
