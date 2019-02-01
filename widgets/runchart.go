package widgets

import (
	"fmt"
	"github.com/sqshq/vcmd/data"
	"github.com/sqshq/vcmd/settings"
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
	yAxisLabelsWidth  = 5
	yAxisLabelsGap    = 1
	xAxisLegendWidth  = 15
)

type RunChart struct {
	Block
	lines []TimeLine
	grid  ChartGrid
	mutex *sync.Mutex
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
		Block: block,
		lines: []TimeLine{},
		mutex: &sync.Mutex{},
	}
}

func (self *RunChart) newChartGrid() ChartGrid {

	linesCount := (self.Inner.Max.X - self.Inner.Min.X) / (xAxisLabelsGap + xAxisLabelsWidth)
	paddingDuration := time.Duration(time.Second) // TODO support others and/or adjust automatically depending on refresh rate

	return ChartGrid{
		linesCount:      linesCount,
		paddingDuration: paddingDuration,
		paddingWidth:    xAxisLabelsGap + xAxisLabelsWidth,
		maxTimeWidth:    self.Inner.Max.X - xAxisLabelsWidth/2 - xAxisLabelsGap,
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
		self.Inner.Min.X+yAxisLabelsWidth+1, self.Inner.Min.Y,
		self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	self.renderItems(buf, drawArea)
	self.renderLegend(buf, drawArea)
	self.mutex.Unlock()
}

func (self *RunChart) ConsumeValue(item data.Item, value string) {

	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatalf("Expected float number, but got %v", value) // TODO visual notification
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

func (self *RunChart) trimOutOfRangeValues() {
	for i, item := range self.lines {
		lastOutOfRangeValueIndex := -1

		for j, timePoint := range item.points {
			if !self.isTimePointInRange(timePoint) {
				lastOutOfRangeValueIndex = j
			}
		}

		if lastOutOfRangeValueIndex > 0 {
			item.points = append(item.points[:0], item.points[lastOutOfRangeValueIndex+1:]...)
			self.lines[i] = item
		}
	}
}

func (self *RunChart) renderAxes(buffer *Buffer) {
	// draw origin cell
	buffer.SetCell(
		NewCell(BOTTOM_LEFT, NewStyle(ColorWhite)),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)

	// draw x axis line
	for i := yAxisLabelsWidth + 1; i < self.Inner.Dx(); i++ {
		buffer.SetCell(
			NewCell(HORIZONTAL_DASH, NewStyle(ColorWhite)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}

	// draw grid
	for y := 0; y < self.Inner.Dy()-xAxisLabelsHeight-1; y = y + 2 {
		for x := 0; x < self.grid.linesCount; x++ {
			buffer.SetCell(
				NewCell(VERTICAL_DASH, NewStyle(settings.ColorDarkGrey)),
				image.Pt(self.grid.maxTimeWidth-x*self.grid.paddingWidth, y+self.Inner.Min.Y+1),
			)
		}
	}

	// draw y axis line
	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buffer.SetCell(
			NewCell(VERTICAL_DASH, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Min.X+yAxisLabelsWidth, i+self.Inner.Min.Y),
		)
	}

	// draw x axis time labels
	for i := 0; i < self.grid.linesCount; i++ {
		labelTime := self.grid.timeExtremum.max.Add(time.Duration(-i) * time.Second)
		buffer.SetString(
			labelTime.Format("15:04:05"),
			NewStyle(ColorWhite),
			image.Pt(self.grid.maxTimeWidth-xAxisLabelsWidth/2-i*(self.grid.paddingWidth), self.Inner.Max.Y-1),
		)
	}

	// draw y axis labels
	verticalScale := self.grid.valueExtremum.max - self.grid.valueExtremum.min/float64(self.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 1; i*(yAxisLabelsGap+1) <= self.Inner.Dy()-1; i++ {
		buffer.SetString(
			fmt.Sprintf("%.3f", float64(i)*self.grid.valueExtremum.min*verticalScale*(yAxisLabelsGap+1)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func (self *RunChart) renderItems(buffer *Buffer, drawArea image.Rectangle) {

	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	for _, line := range self.lines {

		xToPoint := make(map[int]image.Point)
		pointsOrder := make([]int, 0)

		for _, point := range line.points {

			if !self.isTimePointInRange(point) {
				continue
			}

			timeDeltaWithGridMaxTime := self.grid.timeExtremum.max.Sub(point.Time)
			deltaToPaddingRelation := float64(timeDeltaWithGridMaxTime.Nanoseconds()) / float64(self.grid.paddingDuration.Nanoseconds())
			x := self.grid.maxTimeWidth - (int(float64(self.grid.paddingWidth) * deltaToPaddingRelation))

			valuePerYDot := (self.grid.valueExtremum.max - self.grid.valueExtremum.min) / float64(drawArea.Dy()-1)
			y := int(float64(point.Value-self.grid.valueExtremum.min) / valuePerYDot)

			if _, exists := xToPoint[x]; exists {
				continue
			}

			xToPoint[x] = image.Pt(x, drawArea.Max.Y-y-1)
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

			//buffer.SetCell(
			//	NewCell(self.DotRune, NewStyle(SelectColor(self.LineColors, 0))),
			//	currentPoint,
			//)

			canvas.Line(
				braillePoint(previousPoint),
				braillePoint(currentPoint),
				line.item.Color,
			)
		}
	}

	canvas.Draw(buffer)
}

func (self *RunChart) renderLegend(buffer *Buffer, rectangle image.Rectangle) {
	for i, line := range self.lines {

		extremum := GetLineValueExtremum(line.points)

		buffer.SetString(
			fmt.Sprintf("•"),
			NewStyle(line.item.Color),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth-2, self.Inner.Min.Y+1+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("%s", line.item.Label),
			NewStyle(line.item.Color),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+1+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("max %.3f", extremum.max),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+2+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("min %.3f", extremum.min),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+3+i*5),
		)
		buffer.SetString(
			fmt.Sprintf("cur %.3f", line.points[len(line.points)-1].Value),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Max.X-xAxisLegendWidth, self.Inner.Min.Y+4+i*5),
		)
	}
}

func (self *RunChart) isTimePointInRange(point TimePoint) bool {
	return point.Time.After(self.grid.timeExtremum.min.Add(self.grid.paddingDuration))
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
