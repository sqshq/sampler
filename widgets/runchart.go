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
	chartHistoryReserve = 10
	xAxisLabelsHeight   = 1
	xAxisLabelsWidth    = 8
	xAxisLabelsGap      = 2
	yAxisLabelsHeight   = 1
	yAxisLabelsGap      = 1
	xAxisLegendWidth    = 20
)

type RunChart struct {
	Block
	lines     []TimeLine
	grid      ChartGrid
	precision int
	selection *time.Time
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
	linesCount      int
	paddingDuration time.Duration
	paddingWidth    int
	maxTimeWidth    int
	minTimeWidth    int
	valueExtremum   ValueExtremum
	timeExtremum    TimeExtremum
}

type TimeExtremum struct {
	max time.Time
	min time.Time
}

type ValueExtremum struct {
	max float64
	min float64
}

func NewRunChart(title string) *RunChart {
	block := *NewBlock()
	block.Title = title
	return &RunChart{
		Block:     block,
		lines:     []TimeLine{},
		mutex:     &sync.Mutex{},
		precision: 2, // TODO move to config
	}
}

func (self *RunChart) newChartGrid() ChartGrid {

	linesCount := (self.Inner.Max.X - self.Inner.Min.X - self.grid.minTimeWidth) / (xAxisLabelsGap + xAxisLabelsWidth)
	paddingDuration := time.Duration(time.Second) // TODO support others and/or adjust automatically depending on refresh rate

	return ChartGrid{
		linesCount:      linesCount,
		paddingDuration: paddingDuration,
		paddingWidth:    xAxisLabelsGap + xAxisLabelsWidth,
		maxTimeWidth:    self.Inner.Max.X,
		minTimeWidth:    self.getMaxValueLength(),
		timeExtremum:    GetTimeExtremum(linesCount, paddingDuration),
		valueExtremum:   GetChartValueExtremum(self.lines),
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

	selectedPoints := self.getSelectedTimePoints()

	self.renderAxes(buffer)
	self.renderItems(buffer, drawArea)
	self.renderSelection(buffer, drawArea, selectedPoints)
	self.renderLegend(buffer, drawArea, selectedPoints)
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

func (self *RunChart) SelectPoint(x int, y int) {

	point := image.Point{X: x, Y: y}

	if !point.In(self.Rectangle) {
		self.selection = nil
		return
	}

	timeDeltaToPaddingRelation := (self.grid.maxTimeWidth - x) / self.grid.paddingWidth
	timeDeltaWithGridMaxTime := timeDeltaToPaddingRelation * int(self.grid.paddingDuration.Nanoseconds())
	selection := self.grid.timeExtremum.max.Add(-time.Duration(timeDeltaWithGridMaxTime) * time.Nanosecond)

	self.selection = &selection
}

func (self *RunChart) getSelectedTimePoints() []TimePoint {

	selected := []TimePoint{}

	if self.selection == nil {
		return selected
	}

	for _, line := range self.lines {

		if len(line.points) == 0 {
			continue
		}

		closest := line.points[0]

		for _, point := range line.points {

			diffWithClosest := math.Abs(float64(self.selection.UnixNano() - closest.time.UnixNano()))
			diffWithCurrent := math.Abs(float64(self.selection.UnixNano() - point.time.UnixNano()))

			if diffWithClosest > diffWithCurrent {
				closest = point
			}
		}

		selected = append(selected, closest)
	}

	return selected
}

func (self *RunChart) trimOutOfRangeValues() {

	historyReserve := self.grid.paddingDuration * time.Duration(self.grid.linesCount) * chartHistoryReserve
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

func (self *RunChart) renderItems(buffer *Buffer, drawArea image.Rectangle) {

	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	for _, line := range self.lines {

		xToPoint := make(map[int]image.Point)
		pointsOrder := make([]int, 0)

		for _, timePoint := range line.points {

			timeDeltaWithGridMaxTime := self.grid.timeExtremum.max.Sub(timePoint.time).Nanoseconds()
			timeDeltaToPaddingRelation := float64(timeDeltaWithGridMaxTime) / float64(self.grid.paddingDuration.Nanoseconds())
			x := self.grid.maxTimeWidth - (int(float64(self.grid.paddingWidth) * timeDeltaToPaddingRelation))

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
				image.Pt(self.grid.maxTimeWidth-x*self.grid.paddingWidth, y+self.Inner.Min.Y+1),
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
		labelTime := self.grid.timeExtremum.max.Add(time.Duration(-i) * self.grid.paddingDuration)
		buffer.SetString(
			labelTime.Format("15:04:05"),
			NewStyle(ColorWhite),
			image.Pt(self.grid.maxTimeWidth-xAxisLabelsWidth/2-i*(self.grid.paddingWidth), self.Inner.Max.Y-1),
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

func (self *RunChart) renderLegend(buffer *Buffer, rectangle image.Rectangle, selectedPoints []TimePoint) {

	for i, line := range self.lines {

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

		if len(selectedPoints) > 0 {

			index := -1

			for i, p := range selectedPoints {
				if p.line.label == line.label {
					index = i
				}
			}

			if index != -1 {
				buffer.SetString(
					fmt.Sprintf("time:  %v", selectedPoints[index].time.Format("15:04:05.000")),
					NewStyle(ColorWhite),
					image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+2+i*5),
				)
				buffer.SetString(
					fmt.Sprintf("value: %s", formatValue(selectedPoints[index].value, self.precision)),
					NewStyle(ColorWhite),
					image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+3+i*5),
				)
			}
		} else {
			extremum := GetLineValueExtremum(line.points)

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
}

func (self *RunChart) renderSelection(buffer *Buffer, drawArea image.Rectangle, selectedPoints []TimePoint) {

	for _, timePoint := range selectedPoints {

		timeDeltaWithGridMaxTime := self.grid.timeExtremum.max.Sub(timePoint.time).Nanoseconds()
		timeDeltaToPaddingRelation := float64(timeDeltaWithGridMaxTime) / float64(self.grid.paddingDuration.Nanoseconds())
		x := self.grid.maxTimeWidth - (int(float64(self.grid.paddingWidth) * timeDeltaToPaddingRelation))

		var y int
		if self.grid.valueExtremum.max-self.grid.valueExtremum.min == 0 {
			y = (drawArea.Dy() - 2) / 2
		} else {
			valuePerY := (self.grid.valueExtremum.max - self.grid.valueExtremum.min) / float64(drawArea.Dy()-2)
			y = int(float64(timePoint.value-self.grid.valueExtremum.min) / valuePerY)
		}

		point := image.Pt(x, drawArea.Max.Y-y-1)

		if point.In(drawArea) {
			buffer.SetCell(NewCell('▲', NewStyle(timePoint.line.color)), point)
		}
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

func GetChartValueExtremum(items []TimeLine) ValueExtremum {

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

func GetLineValueExtremum(points []TimePoint) ValueExtremum {

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

func GetTimeExtremum(linesCount int, paddingDuration time.Duration) TimeExtremum {
	maxTime := time.Now()
	return TimeExtremum{
		max: maxTime,
		min: maxTime.Add(-time.Duration(paddingDuration.Nanoseconds() * int64(linesCount))),
	}
}
