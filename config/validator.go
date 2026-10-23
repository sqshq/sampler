package config

import (
	"fmt"
	"github.com/sqshq/sampler/console"
)

func (c *Config) Validate() error {

	var components []ComponentConfig

	for _, c := range c.RunCharts {
		components = append(components, c.ComponentConfig)
		if err := checkLabelsUniqueness(c.Title, c.Items); err != nil {
			return err
		}
		if err := checkItemsScripts(c.Title, c.Items); err != nil {
			return err
		}
	}
	for _, c := range c.BarCharts {
		components = append(components, c.ComponentConfig)
		if err := checkLabelsUniqueness(c.Title, c.Items); err != nil {
			return err
		}
		if err := checkItemsScripts(c.Title, c.Items); err != nil {
			return err
		}
	}
	for _, c := range c.SparkLines {
		components = append(components, c.ComponentConfig)
		if err := checkItemScripts(c.Title, c.Item); err != nil {
			return err
		}
	}
	for _, c := range c.Gauges {
		components = append(components, c.ComponentConfig)
		if err := checkItemScripts(c.Title, c.Min); err != nil {
			return err
		}
		if err := checkItemScripts(c.Title, c.Max); err != nil {
			return err
		}
		if err := checkItemScripts(c.Title, c.Cur); err != nil {
			return err
		}
	}
	for _, c := range c.AsciiBoxes {
		components = append(components, c.ComponentConfig)
		if err := checkItemScripts(c.Title, c.Item); err != nil {
			return err
		}
	}
	for _, c := range c.TextBoxes {
		components = append(components, c.ComponentConfig)
		if err := checkItemScripts(c.Title, c.Item); err != nil {
			return err
		}
	}

	return checkTitlesUniqueness(components)
}

func (c *Config) validate() {
	if err := c.Validate(); err != nil {
		console.Exit(err.Error())
	}
}

func checkItemsScripts(title string, items []Item) error {
	for _, i := range items {
		if err := checkItemScripts(title, i); err != nil {
			return err
		}
	}
	return nil
}

func checkItemScripts(title string, i Item) error {
	if i.InitScript != nil && i.MultiStepInitScript != nil {
		return fmt.Errorf("Config validation error: both init and multistep-init scripts are not allowed for '%s'", title)
	}
	if i.SampleScript == nil {
		return fmt.Errorf("Config validation error: sample script should be specified for '%s'", title)
	}
	return nil
}

func checkLabelsUniqueness(title string, items []Item) error {
	labels := make(map[string]bool)
	for _, i := range items {
		if i.Label != nil {
			label := *i.Label
			if _, contains := labels[label]; contains {
				return fmt.Errorf("Config validation error: item labels should be unique. Please rename '%s' for '%s'", label, title)
			}
			labels[label] = true
		}
	}
	return nil
}

func checkTitlesUniqueness(components []ComponentConfig) error {
	titles := make(map[string]bool)
	for _, c := range components {
		if _, contains := titles[c.Title]; contains {
			return fmt.Errorf("Config validation error: component titles should be unique. Please rename '%s'", c.Title)
		}
		titles[c.Title] = true
	}
	return nil
}

func validateItemsScripts(title string, items []Item) {
	if err := checkItemsScripts(title, items); err != nil {
		console.Exit(err.Error())
	}
}

func validateItemScripts(title string, i Item) {
	if err := checkItemScripts(title, i); err != nil {
		console.Exit(err.Error())
	}
}

func validateLabelsUniqueness(title string, items []Item) {
	if err := checkLabelsUniqueness(title, items); err != nil {
		console.Exit(err.Error())
	}
}

func validateTitlesUniqueness(components []ComponentConfig) {
	if err := checkTitlesUniqueness(components); err != nil {
		console.Exit(err.Error())
	}
}
