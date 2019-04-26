package config

import (
	"github.com/sqshq/sampler/console"
)

func (c *Config) setDefaultArrangement() {

	components := getComponents(c)

	if allHaveNoPosition(components) {
		setSingleComponentPosition(components[0])
	}

	for _, component := range components {
		if component.Position == nil {

			lc := getLargestComponent(components)
			lr := lc.GetRectangle()

			var w, h, ws, hs int

			if lr.Dx()/2 > lr.Dy() {
				w = lr.Dx() / 2
				h = lr.Dy()
				ws = w
				hs = 0
			} else {
				w = lr.Dx()
				h = lr.Dy() / 2
				ws = 0
				hs = h
			}

			component.Position = [][]int{
				{lr.Min.X + ws, lr.Min.Y + hs},
				{w, h},
			}

			lc.Position = [][]int{
				{lr.Min.X, lr.Min.Y},
				{w, h},
			}
		}
	}
}

func allHaveNoPosition(components []*ComponentConfig) bool {
	for _, component := range components {
		if component.Position != nil {
			return false
		}
	}
	return true
}

func getLargestComponent(components []*ComponentConfig) *ComponentConfig {

	largestComponent := components[0]

	for _, component := range components {
		if getSquare(component) > getSquare(largestComponent) {
			largestComponent = component
		}
	}

	return largestComponent
}

func getSquare(c *ComponentConfig) int {
	r := c.GetRectangle()
	return r.Dx() * r.Dy()
}

func getComponents(c *Config) []*ComponentConfig {

	var components []*ComponentConfig

	for i := range c.RunCharts {
		components = append(components, &c.RunCharts[i].ComponentConfig)
	}
	for i := range c.BarCharts {
		components = append(components, &c.BarCharts[i].ComponentConfig)
	}
	for i := range c.Gauges {
		components = append(components, &c.Gauges[i].ComponentConfig)
	}
	for i := range c.SparkLines {
		components = append(components, &c.SparkLines[i].ComponentConfig)
	}
	for i := range c.AsciiBoxes {
		components = append(components, &c.AsciiBoxes[i].ComponentConfig)
	}
	for i := range c.TextBoxes {
		components = append(components, &c.TextBoxes[i].ComponentConfig)
	}

	return components
}

func setSingleComponentPosition(c *ComponentConfig) {
	w := int(console.ColumnsCount * 0.8)
	h := int(console.RowsCount * 0.7)
	c.Position = [][]int{
		{(console.ColumnsCount - w) / 2, (console.RowsCount - h) / 2},
		{w, h},
	}
}
