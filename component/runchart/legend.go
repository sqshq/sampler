package runchart

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component/util"
	"image"
	"math"
)

const (
	xAxisLegendIndent = 10
	yAxisLegendIndent = 1
	heightOnDefault   = 2
	heightOnPinpoint  = 4
	heightOnDetails   = 6

	timeFormat = "15:04:05.000"
)

type legend struct {
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

	rowCount := (c.Dy() - 2*yAxisLegendIndent) / (height + yAxisLegendIndent)
	columnCount := int(math.Ceil(float64(len(c.lines)) / float64(rowCount)))
	columnWidth := getColumnWidth(c.mode, c.lines, c.scale)

	for col := 0; col < columnCount; col++ {
		for row := 0; row < rowCount; row++ {

			lineIndex := row + rowCount*col
			if len(c.lines) <= lineIndex {
				break
			}

			line := c.lines[row+rowCount*col]

			x := c.Inner.Max.X - (columnWidth+xAxisLegendIndent)*(col+1)
			y := c.Inner.Min.Y + yAxisLegendIndent + row*height

			titleStyle := ui.NewStyle(line.color)
			detailsStyle := ui.NewStyle(c.palette.BaseColor)

			buffer.SetString(string(ui.DOT), titleStyle, image.Pt(x-2, y))
			buffer.SetString(line.label, titleStyle, image.Pt(x, y))

			if c.mode == ModePinpoint {
				buffer.SetString(fmt.Sprintf("time  %s", line.selectionPoint.time.Format("15:04:05.000")), detailsStyle, image.Pt(x, y+1))
				buffer.SetString(fmt.Sprintf("value %s", util.FormatValue(line.selectionPoint.value, c.scale)), detailsStyle, image.Pt(x, y+2))
				continue
			}

			if !c.legend.Details {
				continue
			}

			details := [4]string{
				fmt.Sprintf("cur  %s", util.FormatValue(getCurrentValue(line), c.scale)),
				fmt.Sprintf("dlt %s", util.FormatDelta(getDiffWithPreviousValue(line), c.scale)),
				fmt.Sprintf("max  %s", util.FormatValue(line.extrema.max, c.scale)),
				fmt.Sprintf("min  %s", util.FormatValue(line.extrema.min, c.scale)),
			}

			for i, detail := range details {
				buffer.SetString(detail, detailsStyle, image.Pt(x, y+i+yAxisLegendIndent))
			}
		}
	}
}

func getColumnWidth(mode Mode, lines []TimeLine, scale int) int {

	if mode == ModePinpoint {
		return len(timeFormat)
	}

	width := len(util.FormatValue(0, scale))
	for _, line := range lines {
		if len(line.label) > width {
			width = len(line.label)
		}
	}
	return width
}

func getDiffWithPreviousValue(line TimeLine) float64 {
	if len(line.points) < 2 {
		return 0
	}
	return line.points[len(line.points)-1].value - line.points[len(line.points)-2].value
}

func getCurrentValue(line TimeLine) float64 {
	if len(line.points) == 0 {
		return 0
	}
	return line.points[len(line.points)-1].value
}
