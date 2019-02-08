package widgets

import (
	"fmt"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"math"
	"strconv"
	"sync"
	"time"

	ui "github.com/sqshq/termui"
)

// TODO split into runchart, grid, legend files
const (
	xAxisLegendWidth  = 20
	xAxisLabelsHeight = 1
	xAxisLabelsWidth  = 8
	xAxisLabelsGap    = 2
	xAxisGridWidth    = xAxisLabelsGap + xAxisLabelsWidth
	yAxisLabelsHeight = 1
	yAxisLabelsGap    = 1
)

type ScrollMode int

const (
	Auto   ScrollMode = 0
	Manual ScrollMode = 1
)

type RunChart struct {
	ui.Block
	lines      []TimeLine
	grid       ChartGrid
	timescale  time.Duration
	mutex      *sync.Mutex
	scrollMode ScrollMode
	selection  time.Time
	precision  int
}

type ChartGrid struct {
	timeRange    TimeRange
	timePerPoint time.Duration
	valueExtrema ValueExtrema
	linesCount   int
	maxTimeWidth int
	minTimeWidth int
}

type TimePoint struct {
	value      float64
	time       time.Time
	coordinate int
}

type TimeLine struct {
	points    []TimePoint
	color     ui.Color
	label     string
	selection int
}

type TimeRange struct {
	max time.Time
	min time.Time
}

type ValueExtrema struct {
	max float64
	min float64
}

func NewRunChart(title string, precision int, refreshRateMs int) *RunChart {
	block := *ui.NewBlock()
	block.Title = title
	return &RunChart{
		Block:      block,
		lines:      []TimeLine{},
		timescale:  calculateTimescale(refreshRateMs),
		mutex:      &sync.Mutex{},
		precision:  precision,
		scrollMode: Auto,
	}
}

func (self *RunChart) newChartGrid() ChartGrid {

	linesCount := (self.Inner.Max.X - self.Inner.Min.X - self.grid.minTimeWidth) / xAxisGridWidth
	timeRange := self.getTimeRange(linesCount)

	return ChartGrid{
		timeRange:    timeRange,
		timePerPoint: self.timescale / time.Duration(xAxisGridWidth),
		valueExtrema: getValueExtrema(self.lines, timeRange),
		linesCount:   linesCount,
		maxTimeWidth: self.Inner.Max.X,
		minTimeWidth: self.getMaxValueLength(),
	}
}

func (self *RunChart) newTimePoint(value float64) TimePoint {
	now := time.Now()
	return TimePoint{
		value:      value,
		time:       now,
		coordinate: self.calculateTimeCoordinate(now),
	}
}

func (self *RunChart) Draw(buffer *ui.Buffer) {

	self.mutex.Lock()
	self.Block.Draw(buffer)
	self.grid = self.newChartGrid()

	drawArea := image.Rect(
		self.Inner.Min.X+self.grid.minTimeWidth+1, self.Inner.Min.Y,
		self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	self.renderAxes(buffer)
	self.renderLines(buffer, drawArea)
	self.renderLegend(buffer, drawArea)
	self.mutex.Unlock()
}

func (self *RunChart) ConsumeSample(sample data.Sample) {

	float, err := strconv.ParseFloat(sample.Value, 64)

	if err != nil {
		// TODO visual notification + check sample.Error
	}

	self.mutex.Lock()

	lineIndex := -1

	for i, line := range self.lines {
		if line.label == sample.Label {
			lineIndex = i
		}
	}

	if lineIndex == -1 {
		line := &TimeLine{
			points: []TimePoint{},
			color:  sample.Color,
			label:  sample.Label,
		}
		self.lines = append(self.lines, *line)
		lineIndex = len(self.lines) - 1
	}

	line := self.lines[lineIndex]
	timePoint := self.newTimePoint(float)
	line.points = append(line.points, timePoint)
	self.lines[lineIndex] = line

	self.trimOutOfRangeValues()
	self.mutex.Unlock()
}

func (self *RunChart) renderLines(buffer *ui.Buffer, drawArea image.Rectangle) {

	canvas := ui.NewCanvas()
	canvas.Rectangle = drawArea

	if len(self.lines) == 0 || len(self.lines[0].points) == 0 {
		return
	}

	selectionCoordinate := self.calculateTimeCoordinate(self.selection)
	selectionPoints := make(map[int]image.Point)

	probe := self.lines[0].points[0]
	delta := ui.AbsInt(self.calculateTimeCoordinate(probe.time) - probe.coordinate)

	for i, line := range self.lines { // TODO start from right side, break on out of range

		xPoint := make(map[int]image.Point)
		xOrder := make([]int, 0)

		if line.selection != 0 {
			line.selection -= delta
			self.lines[i].selection = line.selection
		}

		for j, timePoint := range line.points {

			timePoint.coordinate -= delta
			line.points[j] = timePoint

			var y int
			if self.grid.valueExtrema.max == self.grid.valueExtrema.min {
				y = (drawArea.Dy() - 2) / 2
			} else {
				valuePerY := (self.grid.valueExtrema.max - self.grid.valueExtrema.min) / float64(drawArea.Dy()-2)
				y = int(float64(timePoint.value-self.grid.valueExtrema.min) / valuePerY)
			}

			point := image.Pt(timePoint.coordinate, drawArea.Max.Y-y-1)

			if _, exists := xPoint[point.X]; exists {
				continue
			}

			if !point.In(drawArea) {
				continue
			}

			if line.selection == 0 {
				if len(line.points) > j+1 && ui.AbsInt(timePoint.coordinate-selectionCoordinate) > ui.AbsInt(line.points[j+1].coordinate-selectionCoordinate) {
					selectionPoints[i] = point
				}
			} else {
				if timePoint.coordinate == line.selection {
					selectionPoints[i] = point
				}
			}

			xPoint[point.X] = point
			xOrder = append(xOrder, point.X)
		}

		for i, x := range xOrder {

			currentPoint := xPoint[x]
			var previousPoint image.Point

			if i == 0 {
				previousPoint = currentPoint
			} else {
				previousPoint = xPoint[xOrder[i-1]]
			}

			canvas.Line(
				braillePoint(previousPoint),
				braillePoint(currentPoint),
				line.color,
			)
		}
	}

	canvas.Draw(buffer)

	if self.scrollMode == Manual {
		for lineIndex, point := range selectionPoints {
			buffer.SetCell(ui.NewCell(console.SymbolSelection, ui.NewStyle(self.lines[lineIndex].color)), point)
			if self.lines[lineIndex].selection == 0 {
				self.lines[lineIndex].selection = point.X
			}
		}
	}
}

func (self *RunChart) renderAxes(buffer *ui.Buffer) {
	// draw origin cell
	buffer.SetCell(
		ui.NewCell(ui.BOTTOM_LEFT, ui.NewStyle(ui.ColorWhite)),
		image.Pt(self.Inner.Min.X+self.grid.minTimeWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)

	// draw x axis line
	for i := self.grid.minTimeWidth + 1; i < self.Inner.Dx(); i++ {
		buffer.SetCell(
			ui.NewCell(ui.HORIZONTAL_DASH, ui.NewStyle(ui.ColorWhite)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}

	// draw grid lines
	for y := 0; y < self.Inner.Dy()-xAxisLabelsHeight-2; y = y + 2 {
		for x := 1; x <= self.grid.linesCount; x++ {
			buffer.SetCell(
				ui.NewCell(ui.VERTICAL_DASH, ui.NewStyle(console.ColorDarkGrey)),
				image.Pt(self.grid.maxTimeWidth-x*xAxisGridWidth, y+self.Inner.Min.Y+1),
			)
		}
	}

	// draw y axis line
	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buffer.SetCell(
			ui.NewCell(ui.VERTICAL_DASH, ui.NewStyle(ui.ColorWhite)),
			image.Pt(self.Inner.Min.X+self.grid.minTimeWidth, i+self.Inner.Min.Y),
		)
	}

	// draw x axis time labels
	for i := 1; i <= self.grid.linesCount; i++ {
		labelTime := self.grid.timeRange.max.Add(time.Duration(-i) * self.timescale)
		buffer.SetString(
			labelTime.Format("15:04:05"),
			ui.NewStyle(ui.ColorWhite),
			image.Pt(self.grid.maxTimeWidth-xAxisLabelsWidth/2-i*(xAxisGridWidth), self.Inner.Max.Y-1),
		)
	}

	// draw y axis labels
	if self.grid.valueExtrema.max != self.grid.valueExtrema.min {
		labelsCount := (self.Inner.Dy() - xAxisLabelsHeight - 1) / (yAxisLabelsGap + yAxisLabelsHeight)
		valuePerY := (self.grid.valueExtrema.max - self.grid.valueExtrema.min) / float64(self.Inner.Dy()-xAxisLabelsHeight-3)
		for i := 0; i < int(labelsCount); i++ {
			value := self.grid.valueExtrema.max - (valuePerY * float64(i) * (yAxisLabelsGap + yAxisLabelsHeight))
			buffer.SetString(
				formatValue(value, self.precision),
				ui.NewStyle(ui.ColorWhite),
				image.Pt(self.Inner.Min.X, 1+self.Inner.Min.Y+i*(yAxisLabelsGap+yAxisLabelsHeight)),
			)
		}
	} else {
		buffer.SetString(
			formatValue(self.grid.valueExtrema.max, self.precision),
			ui.NewStyle(ui.ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Min.Y+self.Inner.Dy()/2))
	}
}

func (self *RunChart) renderLegend(buffer *ui.Buffer, rectangle image.Rectangle) {

	for i, line := range self.lines {

		extremum := getLineValueExtremum(line.points)

		buffer.SetString(
			string(ui.DOT),
			ui.NewStyle(line.color),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth-2, self.Inner.Min.Y+1+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("%s", line.label),
			ui.NewStyle(line.color),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+1+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("cur %s", formatValue(line.points[len(line.points)-1].value, self.precision)),
			ui.NewStyle(ui.ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+2+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("max %s", formatValue(extremum.max, self.precision)),
			ui.NewStyle(ui.ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+3+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("min %s", formatValue(extremum.min, self.precision)),
			ui.NewStyle(ui.ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+4+i*5),
		)
	}
}

func (self *RunChart) trimOutOfRangeValues() {
	// TODO use hard limit
}

func (self *RunChart) calculateTimeCoordinate(t time.Time) int {
	timeDeltaWithGridMaxTime := self.grid.timeRange.max.Sub(t).Nanoseconds()
	timeDeltaToPaddingRelation := float64(timeDeltaWithGridMaxTime) / float64(self.timescale.Nanoseconds())
	return self.grid.maxTimeWidth - int(math.Ceil(float64(xAxisGridWidth)*timeDeltaToPaddingRelation))
}

func (self *RunChart) getMaxValueLength() int {

	maxValueLength := 0

	for _, line := range self.lines {
		for _, point := range line.points {
			l := len(formatValue(point.value, self.precision))
			if l > maxValueLength {
				maxValueLength = l
			}
		}
	}

	return maxValueLength
}

func (self *RunChart) MoveSelection(shift int) {

	if self.scrollMode == Auto {
		self.scrollMode = Manual
		self.selection = getMidRangeTime(self.grid.timeRange)
		return
	} else {
		self.selection = self.selection.Add(self.grid.timePerPoint * time.Duration(shift))
		if self.selection.After(self.grid.timeRange.max) {
			self.selection = self.grid.timeRange.max
		} else if self.selection.Before(self.grid.timeRange.min) {
			self.selection = self.grid.timeRange.min
		}
	}

	for i := range self.lines {
		self.lines[i].selection = 0
	}
}

func (self *RunChart) DisableSelection() {
	if self.scrollMode == Manual {
		self.scrollMode = Auto
		return
	}
}

func getMidRangeTime(r TimeRange) time.Time {
	delta := r.max.Sub(r.min)
	return r.max.Add(-delta / 2)
}

func formatValue(value float64, precision int) string {
	format := "%." + strconv.Itoa(precision) + "f"
	return fmt.Sprintf(format, value)
}

func getValueExtrema(items []TimeLine, timeRange TimeRange) ValueExtrema {

	if len(items) == 0 {
		return ValueExtrema{0, 0}
	}

	var max, min = -math.MaxFloat64, math.MaxFloat64

	for _, item := range items {
		for _, point := range item.points {
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

func getLineValueExtremum(points []TimePoint) ValueExtrema {

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

func (self *RunChart) getTimeRange(linesCount int) TimeRange {

	if self.scrollMode == Manual {
		return self.grid.timeRange
	}

	width := time.Duration(self.timescale.Nanoseconds() * int64(linesCount))
	max := time.Now()

	return TimeRange{
		max: max,
		min: max.Add(-width),
	}
}

// time duration between grid lines
func calculateTimescale(refreshRateMs int) time.Duration {

	multiplier := refreshRateMs * xAxisGridWidth / 2
	timescale := time.Duration(time.Millisecond * time.Duration(multiplier)).Round(time.Second)

	if timescale.Seconds() == 0 {
		return time.Second
	} else {
		return timescale
	}
}
