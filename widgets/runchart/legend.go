package runchart

import (
	"fmt"
	ui "github.com/sqshq/termui"
	"image"
	"math"
)

const (
	xAxisLegendIndent = 10
	yAxisLegendIndent = 1
	heightOnDefault   = 2
	heightOnPinpoint  = 4
	heightOnDetails   = 6
)

type Legend struct {
	Enabled bool
	Details bool
}

func (c *RunChart) renderLegend(buffer *ui.Buffer, rectangle image.Rectangle) {

	if !c.legend.Enabled {
		return
	}

	height := heightOnDefault

	if c.mode == ModePinpoint {
		height = heightOnPinpoint
	} else if c.legend.Details {
		height = heightOnDetails
	}

	rowCount := (c.Dx() - yAxisLegendIndent) / (height + yAxisLegendIndent)
	columnCount := int(math.Ceil(float64(len(c.lines)) / float64(rowCount)))
	columnWidth := getColumnWidth(c.lines, c.precision)

	for col := 0; col < columnCount; col++ {
		for row := 0; row < rowCount; row++ {

			lineIndex := row + rowCount*col
			if len(c.lines) <= lineIndex {
				break
			}

			line := c.lines[row+rowCount*col]
			extrema := getLineValueExtrema(line.points)

			x := c.Inner.Max.X - (columnWidth+xAxisLegendIndent)*(col+1)
			y := c.Inner.Min.Y + yAxisLegendIndent + row*height

			titleStyle := ui.NewStyle(line.color)
			detailsStyle := ui.NewStyle(ui.ColorWhite)

			buffer.SetString(string(ui.DOT), titleStyle, image.Pt(x-2, y))
			buffer.SetString(line.label, titleStyle, image.Pt(x, y))

			if c.mode == ModePinpoint {
				continue
			}

			if !c.legend.Details {
				continue
			}

			details := [4]string{
				fmt.Sprintf("cur %s", formatValue(line.points[len(line.points)-1].value, c.precision)),
				fmt.Sprintf("max %s", formatValue(extrema.max, c.precision)),
				fmt.Sprintf("min %s", formatValue(extrema.min, c.precision)),
				fmt.Sprintf("dif %s", formatValue(1, c.precision)),
			}

			for i, detail := range details {
				buffer.SetString(detail, detailsStyle, image.Pt(x, y+i+yAxisLegendIndent))
			}
		}
	}
}

func getColumnWidth(lines []TimeLine, precision int) int {
	width := len(formatValue(0, precision))
	for _, line := range lines {
		if len(line.label) > width {
			width = len(line.label)
		}
	}
	return width
}

// TODO remove and use the one from line
func getLineValueExtrema(points []TimePoint) ValueExtrema {

	if len(points) == 0 {
		return ValueExtrema{0, 0}
	}

	var max, min = -math.MaxFloat64, math.MaxFloat64

	for _, point := range points {
		if point.value > max {
			max = point.value
		}
		if point.value < min {
			min = point.value
		}
	}

	return ValueExtrema{max: max, min: min}
}
