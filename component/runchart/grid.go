package runchart

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component/util"
	"image"
	"math"
	"time"
)

const defaultValueLength = 4

type chartGrid struct {
	timeRange    TimeRange
	timePerPoint time.Duration
	valueExtrema ValueExtrema
	linesCount   int
	maxTimeWidth int
	minTimeWidth int
}

func (c *RunChart) newChartGrid() chartGrid {

	linesCount := (c.Inner.Max.X - c.Inner.Min.X - c.grid.minTimeWidth) / xAxisGridWidth
	timeRange := c.getTimeRange(linesCount)

	return chartGrid{
		timeRange:    timeRange,
		timePerPoint: c.timescale / time.Duration(xAxisGridWidth),
		valueExtrema: getLocalExtrema(c.lines, timeRange),
		linesCount:   linesCount,
		maxTimeWidth: c.Inner.Max.X,
		minTimeWidth: defaultValueLength,
	}
}

func (c *RunChart) renderAxes(buffer *ui.Buffer) {

	// draw y axis labels
	if c.grid.valueExtrema.max != c.grid.valueExtrema.min {
		labelsCount := (c.Inner.Dy() - xAxisLabelsHeight - 1) / (yAxisLabelsIndent + yAxisLabelsHeight)
		valuePerY := (c.grid.valueExtrema.max - c.grid.valueExtrema.min) / float64(c.Inner.Dy()-xAxisLabelsHeight-3)
		for i := 0; i < int(labelsCount); i++ {
			val := c.grid.valueExtrema.max - (valuePerY * float64(i) * (yAxisLabelsIndent + yAxisLabelsHeight))
			fmt := util.FormatValue(val, c.scale)
			if len(fmt) > c.grid.minTimeWidth {
				c.grid.minTimeWidth = len(fmt)
			}
			buffer.SetString(
				fmt,
				ui.NewStyle(c.palette.BaseColor),
				image.Pt(c.Inner.Min.X, 1+c.Inner.Min.Y+i*(yAxisLabelsIndent+yAxisLabelsHeight)))
		}
	} else {
		fmt := util.FormatValue(c.grid.valueExtrema.max, c.scale)
		c.grid.minTimeWidth = len(fmt)
		buffer.SetString(
			fmt,
			ui.NewStyle(c.palette.BaseColor),
			image.Pt(c.Inner.Min.X, c.Inner.Min.Y+c.Inner.Dy()/2))
	}

	// draw origin cell
	buffer.SetCell(
		ui.NewCell(ui.BOTTOM_LEFT, ui.NewStyle(c.palette.BaseColor)),
		image.Pt(c.Inner.Min.X+c.grid.minTimeWidth, c.Inner.Max.Y-xAxisLabelsHeight-1))

	// draw x axis line
	for i := c.grid.minTimeWidth + 1; i < c.Inner.Dx(); i++ {
		buffer.SetCell(
			ui.NewCell(ui.HORIZONTAL_DASH, ui.NewStyle(c.palette.BaseColor)),
			image.Pt(i+c.Inner.Min.X, c.Inner.Max.Y-xAxisLabelsHeight-1))
	}

	// draw grid lines
	for y := 0; y < c.Inner.Dy()-xAxisLabelsHeight-2; y = y + 2 {
		for x := 1; x <= c.grid.linesCount; x++ {
			buffer.SetCell(
				ui.NewCell(ui.VERTICAL_DASH, ui.NewStyle(c.palette.MediumColor)),
				image.Pt(c.grid.maxTimeWidth-x*xAxisGridWidth, y+c.Inner.Min.Y+1))
		}
	}

	// draw y axis line
	for i := 0; i < c.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buffer.SetCell(
			ui.NewCell(ui.VERTICAL_DASH, ui.NewStyle(c.palette.BaseColor)),
			image.Pt(c.Inner.Min.X+c.grid.minTimeWidth, i+c.Inner.Min.Y))
	}

	// draw x axis time labels
	for i := 1; i <= c.grid.linesCount; i++ {
		labelTime := c.grid.timeRange.max.Add(time.Duration(-i) * c.timescale)
		buffer.SetString(
			labelTime.Format("15:04:05"),
			ui.NewStyle(c.palette.BaseColor),
			image.Pt(c.grid.maxTimeWidth-xAxisLabelsWidth/2-i*(xAxisGridWidth), c.Inner.Max.Y-1))
	}
}

func (c *RunChart) getTimeRange(linesCount int) TimeRange {

	if c.mode == ModePinpoint {
		return c.grid.timeRange
	}

	width := time.Duration(c.timescale.Nanoseconds() * int64(linesCount))
	max := time.Now()

	return TimeRange{
		max: max,
		min: max.Add(-width),
	}
}

func getLocalExtrema(items []TimeLine, timeRange TimeRange) ValueExtrema {

	if len(items) == 0 {
		return ValueExtrema{0, 0}
	}

	var max, min = -math.MaxFloat64, math.MaxFloat64

	for _, item := range items {
		started := false
		for i := len(item.points) - 1; i > 0; i-- {
			point := item.points[i]
			if timeRange.isInRange(point.time) {
				started = true
			} else if started == true && !timeRange.isInRange(point.time) {
				break
			}
			if point.value > max && timeRange.isInRange(point.time) {
				max = point.value
			}
			if point.value < min && timeRange.isInRange(point.time) {
				min = point.value
			}
		}
	}

	return ValueExtrema{max: max, min: min}
}

func (r *TimeRange) isInRange(time time.Time) bool {
	return time.After(r.min) && time.Before(r.max)
}
