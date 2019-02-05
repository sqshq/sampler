package widgets

import (
	"fmt"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	. "github.com/sqshq/termui"
)

const (
	xAxisLegendWidth  = 20
	xAxisLabelsHeight = 1
	xAxisLabelsWidth  = 8
	xAxisLabelsGap    = 2
	xAxisGridWidth    = xAxisLabelsGap + xAxisLabelsWidth

	yAxisLabelsHeight = 1
	yAxisLabelsGap    = 1

	chartHistoryReserve = 5
)

type RunChart struct {
	Block
	lines     []TimeLine
	grid      ChartGrid
	precision int
	timescale time.Duration
	mutex     *sync.Mutex
}

type TimePoint struct {
	time  time.Time
	value float64
	line  *TimeLine
}

type TimeLine struct {
	points []TimePoint
	color  Color
	label  string
}

type ChartGrid struct {
	linesCount    int
	paddingWidth  int
	maxTimeWidth  int
	minTimeWidth  int
	valueExtremum ValueExtremum
	timeExtremum  TimeExtremum
}

type TimeExtremum struct {
	max time.Time
	min time.Time
}

type ValueExtremum struct {
	max float64
	min float64
}

func NewRunChart(title string, precision int, refreshRateMs int) *RunChart {
	block := *NewBlock()
	block.Title = title
	return &RunChart{
		Block:     block,
		lines:     []TimeLine{},
		mutex:     &sync.Mutex{},
		precision: precision,
		timescale: calculateTimescale(refreshRateMs),
	}
}

func (self *RunChart) newChartGrid() ChartGrid {

	linesCount := (self.Inner.Max.X - self.Inner.Min.X - self.grid.minTimeWidth) / xAxisGridWidth

	return ChartGrid{
		linesCount:    linesCount,
		paddingWidth:  xAxisGridWidth,
		maxTimeWidth:  self.Inner.Max.X,
		minTimeWidth:  self.getMaxValueLength(),
		timeExtremum:  getTimeExtremum(linesCount, self.timescale),
		valueExtremum: getChartValueExtremum(self.lines),
	}
}

func (self *RunChart) Draw(buffer *Buffer) {

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
		log.Printf("Expected float number, but got %v", sample.Value) // TODO visual notification + check sample.Error
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

	timePoint := TimePoint{value: float, time: time.Now(), line: &line}
	line.points = append(line.points, timePoint)
	self.lines[lineIndex] = line

	self.trimOutOfRangeValues()
	self.mutex.Unlock()
}

func (self *RunChart) trimOutOfRangeValues() {

	historyReserve := self.timescale * time.Duration(self.grid.linesCount) * chartHistoryReserve
	minRangeTime := self.grid.timeExtremum.min.Add(-historyReserve)

	for i, item := range self.lines {
		lastOutOfRangeValueIndex := -1

		for j, point := range item.points {
			if point.time.Before(minRangeTime) {
				lastOutOfRangeValueIndex = j
			}
		}

		if lastOutOfRangeValueIndex > 0 {
			item.points = append(item.points[:0], item.points[lastOutOfRangeValueIndex+1:]...)
			self.lines[i] = item
		}
	}
}

func (self *RunChart) renderLines(buffer *Buffer, drawArea image.Rectangle) {

	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	for _, line := range self.lines {

		xToPoint := make(map[int]image.Point)
		pointsOrder := make([]int, 0)

		for _, timePoint := range line.points {

			timeDeltaWithGridMaxTime := self.grid.timeExtremum.max.Sub(timePoint.time).Nanoseconds()
			timeDeltaToPaddingRelation := float64(timeDeltaWithGridMaxTime) / float64(self.timescale.Nanoseconds())
			x := self.grid.maxTimeWidth - (int(float64(xAxisGridWidth) * timeDeltaToPaddingRelation))

			var y int
			if self.grid.valueExtremum.max-self.grid.valueExtremum.min == 0 {
				y = (drawArea.Dy() - 2) / 2
			} else {
				valuePerY := (self.grid.valueExtremum.max - self.grid.valueExtremum.min) / float64(drawArea.Dy()-2)
				y = int(float64(timePoint.value-self.grid.valueExtremum.min) / valuePerY)
			}

			point := image.Pt(x, drawArea.Max.Y-y-1)

			if _, exists := xToPoint[x]; exists {
				continue
			}

			if !point.In(drawArea) {
				continue
			}

			xToPoint[x] = point
			pointsOrder = append(pointsOrder, x)
		}

		for i, x := range pointsOrder {

			currentPoint := xToPoint[x]
			var previousPoint image.Point

			if i == 0 {
				previousPoint = currentPoint
			} else {
				previousPoint = xToPoint[pointsOrder[i-1]]
			}

			canvas.Line(
				braillePoint(previousPoint),
				braillePoint(currentPoint),
				line.color,
			)
		}
	}

	canvas.Draw(buffer)
}

func (self *RunChart) renderAxes(buffer *Buffer) {
	// draw origin cell
	buffer.SetCell(
		NewCell(BOTTOM_LEFT, NewStyle(ColorWhite)),
		image.Pt(self.Inner.Min.X+self.grid.minTimeWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)

	// draw x axis line
	for i := self.grid.minTimeWidth + 1; i < self.Inner.Dx(); i++ {
		buffer.SetCell(
			NewCell(HORIZONTAL_DASH, NewStyle(ColorWhite)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}

	// draw grid lines
	for y := 0; y < self.Inner.Dy()-xAxisLabelsHeight-2; y = y + 2 {
		for x := 1; x <= self.grid.linesCount; x++ {
			buffer.SetCell(
				NewCell(VERTICAL_DASH, NewStyle(console.ColorDarkGrey)),
				image.Pt(self.grid.maxTimeWidth-x*xAxisGridWidth, y+self.Inner.Min.Y+1),
			)
		}
	}

	// draw y axis line
	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buffer.SetCell(
			NewCell(VERTICAL_DASH, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Min.X+self.grid.minTimeWidth, i+self.Inner.Min.Y),
		)
	}

	// draw x axis time labels
	for i := 1; i <= self.grid.linesCount; i++ {
		labelTime := self.grid.timeExtremum.max.Add(time.Duration(-i) * self.timescale)
		buffer.SetString(
			labelTime.Format("15:04:05"),
			NewStyle(ColorWhite),
			image.Pt(self.grid.maxTimeWidth-xAxisLabelsWidth/2-i*(xAxisGridWidth), self.Inner.Max.Y-1),
		)
	}

	// draw y axis labels
	if self.grid.valueExtremum.max != self.grid.valueExtremum.min {
		labelsCount := (self.Inner.Dy() - xAxisLabelsHeight - 1) / (yAxisLabelsGap + yAxisLabelsHeight)
		valuePerY := (self.grid.valueExtremum.max - self.grid.valueExtremum.min) / float64(self.Inner.Dy()-xAxisLabelsHeight-3)
		for i := 0; i < int(labelsCount); i++ {
			value := self.grid.valueExtremum.max - (valuePerY * float64(i) * (yAxisLabelsGap + yAxisLabelsHeight))
			buffer.SetString(
				formatValue(value, self.precision),
				NewStyle(ColorWhite),
				image.Pt(self.Inner.Min.X, 1+self.Inner.Min.Y+i*(yAxisLabelsGap+yAxisLabelsHeight)),
			)
		}
	} else {
		buffer.SetString(
			formatValue(self.grid.valueExtremum.max, self.precision),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Dy()/2))
	}
}

func (self *RunChart) renderLegend(buffer *Buffer, rectangle image.Rectangle) {

	for i, line := range self.lines {

		extremum := getLineValueExtremum(line.points)

		buffer.SetString(
			string(DOT),
			NewStyle(line.color),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth-2, self.Inner.Min.Y+1+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("%s", line.label),
			NewStyle(line.color),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+1+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("cur %s", formatValue(line.points[len(line.points)-1].value, self.precision)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+2+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("max %s", formatValue(extremum.max, self.precision)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+3+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("min %s", formatValue(extremum.min, self.precision)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+4+i*5),
		)
	}
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

func formatValue(value float64, precision int) string {
	format := "%." + strconv.Itoa(precision) + "f"
	return fmt.Sprintf(format, value)
}

func getChartValueExtremum(items []TimeLine) ValueExtremum {

	if len(items) == 0 {
		return ValueExtremum{0, 0}
	}

	var max, min = -math.MaxFloat64, math.MaxFloat64

	for _, item := range items {
		for _, point := range item.points {
			if point.value > max {
				max = point.value
			}
			if point.value < min {
				min = point.value
			}
		}
	}

	return ValueExtremum{max: max, min: min}
}

func getLineValueExtremum(points []TimePoint) ValueExtremum {

	if len(points) == 0 {
		return ValueExtremum{0, 0}
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

	return ValueExtremum{max: max, min: min}
}

func getTimeExtremum(linesCount int, scale time.Duration) TimeExtremum {
	maxTime := time.Now()
	return TimeExtremum{
		max: maxTime,
		min: maxTime.Add(-time.Duration(scale.Nanoseconds() * int64(linesCount))),
	}
}

func calculateTimescale(refreshRateMs int) time.Duration {

	multiplier := refreshRateMs * xAxisGridWidth / 2
	timescale := time.Duration(time.Millisecond * time.Duration(multiplier)).Round(time.Second)

	if timescale.Seconds() == 0 {
		return time.Second
	} else {
		return timescale
	}
}
