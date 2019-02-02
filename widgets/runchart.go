package widgets

import (
	"fmt"
	"github.com/sqshq/vcmd/console"
	"github.com/sqshq/vcmd/data"
	"image"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	. "github.com/sqshq/termui"
)

const (
	xAxisLabelsHeight = 1
	xAxisLabelsWidth  = 8
	xAxisLabelsGap    = 2
	yAxisLabelsHeight = 1
	yAxisLabelsGap    = 1
	xAxisLegendWidth  = 15
)

type RunChart struct {
	Block
	lines     []TimeLine
	grid      ChartGrid
	precision int
	selection time.Time
	mutex     *sync.Mutex
}

type TimePoint struct {
	Value float64
	Time  time.Time
}

type TimeLine struct {
	points []TimePoint
	item   data.Item
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
		precision: 2, // TODO config
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

func (self *RunChart) Draw(buf *Buffer) {

	self.mutex.Lock()
	self.Block.Draw(buf)
	self.grid = self.newChartGrid()
	self.renderAxes(buf)

	drawArea := image.Rect(
		self.Inner.Min.X+self.grid.minTimeWidth+1, self.Inner.Min.Y,
		self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	self.renderItems(buf, drawArea)
	self.renderLegend(buf, drawArea)
	self.mutex.Unlock()
}

func (self *RunChart) ConsumeValue(item data.Item, value string) {

	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Expected float number, but got %v", value) // TODO visual notification
	}

	timePoint := TimePoint{Value: float, Time: time.Now()}
	self.mutex.Lock()
	itemExists := false

	for i, line := range self.lines {
		if line.item.Label == item.Label {
			line.points = append(line.points, timePoint)
			self.lines[i] = line
			itemExists = true
		}
	}

	if !itemExists {
		item := &TimeLine{
			points: []TimePoint{timePoint},
			item:   item,
		}
		self.lines = append(self.lines, *item)
	}

	self.trimOutOfRangeValues()
	self.mutex.Unlock()
}

func (self *RunChart) ConsumeError(item data.Item, err error) {
	// TODO visual notification
}

func (self *RunChart) SelectValue(x int, y int) {
	// TODO instead of that, find actual time for the given X
	// + make sure that Y is within the given chart
	// once ensured, set "selected time" into the chart structure
	// self.selection = image.Point{X: x, Y: y}
}

func (self *RunChart) trimOutOfRangeValues() {

	minRangeTime := self.grid.timeExtremum.min.Add(-self.grid.paddingDuration * 10)

	for i, item := range self.lines {
		lastOutOfRangeValueIndex := -1

		for j, point := range item.points {
			if point.Time.Before(minRangeTime) {
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

		for _, point := range line.points {

			timeDeltaWithGridMaxTime := self.grid.timeExtremum.max.Sub(point.Time).Nanoseconds()
			timeDeltaToPaddingRelation := float64(timeDeltaWithGridMaxTime) / float64(self.grid.paddingDuration.Nanoseconds())
			x := self.grid.maxTimeWidth - (int(float64(self.grid.paddingWidth) * timeDeltaToPaddingRelation))

			var y int
			if self.grid.valueExtremum.max-self.grid.valueExtremum.min == 0 {
				y = (drawArea.Dy() - 2) / 2
			} else {
				valuePerY := (self.grid.valueExtremum.max - self.grid.valueExtremum.min) / float64(drawArea.Dy()-2)
				y = int(float64(point.Value-self.grid.valueExtremum.min) / valuePerY)
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
				line.item.Color,
			)
		}

		//if point, exists := xToPoint[self.selection.X]; exists {
		//	buffer.SetCell(
		//		NewCell(DOT, NewStyle(line.item.Color)),
		//		point,
		//	)
		//	log.Printf("EXIST!")
		//} else {
		//	//log.Printf("DOES NOT EXIST")
		//}
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

func (self *RunChart) renderLegend(buffer *Buffer, rectangle image.Rectangle) {
	for i, line := range self.lines {

		extremum := GetLineValueExtremum(line.points)

		buffer.SetString(
			string(DOT),
			NewStyle(line.item.Color),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth-2, self.Inner.Min.Y+1+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("%s", line.item.Label),
			NewStyle(line.item.Color),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+1+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("cur %s", formatValue(line.points[len(line.points)-1].Value, self.precision)),
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
			l := len(formatValue(point.Value, self.precision))
			if l > maxValueLength {
				maxValueLength = l
			}
		}
	}

	return maxValueLength
}

func formatValue(value float64, precision int) string {
	format := " %." + strconv.Itoa(precision) + "f"
	return fmt.Sprintf(format, value)
}

func GetChartValueExtremum(items []TimeLine) ValueExtremum {

	if len(items) == 0 {
		return ValueExtremum{0, 0}
	}

	var max, min = -math.MaxFloat64, math.MaxFloat64

	for _, item := range items {
		for _, point := range item.points {
			if point.Value > max {
				max = point.Value
			}
			if point.Value < min {
				min = point.Value
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
		if point.Value > max {
			max = point.Value
		}
		if point.Value < min {
			min = point.Value
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
