package config

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/console"
	"image"
)

func (c *Config) setDefaultArrangement() {

	components := getComponents(c)

	if allHaveNoPosition(components) {
		setSingleComponentPosition(components[0])
	}

	for _, component := range components {
		if component.Position == nil || len(component.Position) == 0 {

			lc := getLargestComponent(components)
			le := getLargestEmptySpaceRectangle(components)

			if getSquare(lc) > le.Dx()*le.Dy()*2 {
				arrangeIntoLargestComponent(component, lc)
			} else {
				arrangeIntoLargestEmptySpace(component, le)
			}
		}
	}
}

func arrangeIntoLargestComponent(component *ComponentConfig, lc *ComponentConfig) {

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

func arrangeIntoLargestEmptySpace(component *ComponentConfig, le image.Rectangle) {
	component.Position = [][]int{
		{le.Min.X, le.Min.Y},
		{le.Dx(), le.Dy()},
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
	w := int(console.ColumnsCount)
	h := int(console.RowsCount * 0.6)
	c.Position = [][]int{
		{(console.ColumnsCount - w) / 2, (console.RowsCount - h) / 2},
		{w, h},
	}
}

func getLargestEmptySpaceRectangle(components []*ComponentConfig) image.Rectangle {

	grid := [console.RowsCount][console.ColumnsCount]int{}

	for _, component := range components {
		rect := component.GetRectangle()
		for r := ui.MinInt(console.RowsCount, rect.Min.Y); r < ui.MinInt(console.RowsCount, rect.Max.Y); r++ {
			for c := ui.MinInt(console.ColumnsCount, rect.Min.X); c < ui.MinInt(console.ColumnsCount, rect.Max.X); c++ {
				grid[r][c] = 1
			}
		}
	}

	mr := image.ZR

	for row := 0; row < console.RowsCount; row++ {
		histogram := createHistogram(grid, row)
		r := calcMaxRectangle(histogram, row)
		if r.Dx()*r.Dy() > mr.Dx()*mr.Dy() {
			mr = r
		}
	}

	return mr
}

func calcMaxRectangle(histogram [console.ColumnsCount]int, row int) image.Rectangle {

	maxRectangle := image.ZR
	maxArea := 0

	for i := 0; i < len(histogram); i++ {

		height := histogram[i]
		if height > maxArea {
			maxArea = height
			maxRectangle = image.Rect(i, row, i, row+height)
		}

		for j := i - 1; j >= 0; j-- {
			width := i - j + 1
			height = ui.MinInt(height, histogram[j])
			if width*height > maxArea {
				maxArea = width * height
				maxRectangle = image.Rect(j, row, j+width, row+height)
			}
		}
	}

	return maxRectangle
}

func createHistogram(grid [console.RowsCount][console.ColumnsCount]int, row int) [console.ColumnsCount]int {
	histogram := [console.ColumnsCount]int{}
	for column := 0; column < console.ColumnsCount; column++ {
		histogram[column] = countEmptyCellsBelow(grid, row, column)
	}
	return histogram
}

func countEmptyCellsBelow(grid [console.RowsCount][console.ColumnsCount]int, row int, column int) int {
	count := 0
	for r := row; r < console.RowsCount; r++ {
		if grid[r][column] == 1 {
			return count
		}
		count++
	}
	return count
}
