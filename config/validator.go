package config

import (
	"fmt"
	"github.com/sqshq/sampler/console"
	"strings"
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
		if c.Min != nil && c.Max != nil && *c.Min >= *c.Max {
			console.Exit(fmt.Sprintf("Config validation error: sparkline '%s' has min (%v) greater than or equal to max (%v)", c.Title, *c.Min, *c.Max))
		}
	}
	for _, c := range c.Gauges {
		components = append(components, c.ComponentConfig)
		if c.IsMulti() {
			validateMultiGauge(c)
		} else {
			validateItemScripts(c.Title, c.Min)
			validateItemScripts(c.Title, c.Max)
			validateItemScripts(c.Title, c.Cur)
		}
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

func validateMultiGauge(g GaugeConfig) {

	if g.Min.SampleScript != nil {
		validateItemScripts(g.Title, g.Min)
	}
	if g.Max.SampleScript != nil {
		validateItemScripts(g.Title, g.Max)
	}

	items := make([]Item, 0)
	for _, it := range g.Items {
		if it.Label == nil || *it.Label == "" {
			console.Exit(fmt.Sprintf("Config validation error: each item of gauge '%s' should have a label", g.Title))
		}
		label := *it.Label
		if label == GaugeMinLabel || label == GaugeMaxLabel || strings.Contains(label, GaugeLabelSeparator) {
			console.Exit(fmt.Sprintf("Config validation error: gauge '%s' item label '%s' is not allowed (reserved words min/max and character '%s')", g.Title, label, GaugeLabelSeparator))
		}
		name := g.Title + " / " + label
		if it.Cur != nil && it.SampleScript != nil {
			console.Exit(fmt.Sprintf("Config validation error: gauge '%s' item '%s' has both 'sample' and 'cur', keep only one", g.Title, label))
		}
		validateItemScripts(name, it.CurItem())
		if it.Min != nil {
			validateItemScripts(name+" min", *it.Min)
		}
		if it.Max != nil {
			validateItemScripts(name+" max", *it.Max)
		} else if g.Max.SampleScript == nil {
			console.Exit(fmt.Sprintf("Config validation error: gauge '%s' item '%s' has no max, define it on the item or on the gauge", g.Title, label))
		}
		items = append(items, it.Item)
	}

	validateLabelsUniqueness(g.Title, items)
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
