package widgets

import (
	"fmt"
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
)

type RunChart struct {
	Block
	items []ChartItem
	grid  ChartGrid
	mutex *sync.Mutex
}

type TimePoint struct {
	Value float64
	Time  time.Time
}

type ChartItem struct {
	timePoints []TimePoint
	label      string
	color      Color
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
		items: []ChartItem{},
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
		valueExtremum:   GetValueExtremum(self.items),
	}
}

func (self *RunChart) Draw(buf *Buffer) {

	self.mutex.Lock()
	self.Block.Draw(buf)
	self.grid = self.newChartGrid()
	self.plotAxes(buf)

	drawArea := image.Rect(
		self.Inner.Min.X+yAxisLabelsWidth+1, self.Inner.Min.Y,
		self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	self.renderItems(buf, drawArea)
	self.mutex.Unlock()
}

func (self *RunChart) ConsumeValue(value string, label string) {

	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatalf("Expected float number, but got %v", value) // TODO visual notification
	}

	timePoint := TimePoint{Value: float, Time: time.Now()}
	self.mutex.Lock()
	itemExists := false

	for i, item := range self.items {
		if item.label == label {
			item.timePoints = append(item.timePoints, timePoint)
			self.items[i] = item
			itemExists = true
		}
	}

	if !itemExists {
		item := &ChartItem{
			timePoints: []TimePoint{timePoint},
			label:      label,
			color:      ColorYellow,
		}
		self.items = append(self.items, *item)
	}

	self.trimOutOfRangeValues()
	self.mutex.Unlock()
}

func (self *RunChart) ConsumeError(err error) {
	// TODO visual notification
}

func (self *RunChart) trimOutOfRangeValues() {
	for i, item := range self.items {
		lastOutOfRangeValueIndex := -1

		for j, timePoint := range item.timePoints {
			if !self.isTimePointInRange(timePoint) {
				lastOutOfRangeValueIndex = j
			}
		}

		if lastOutOfRangeValueIndex > 0 {
			item.timePoints = append(item.timePoints[:0], item.timePoints[lastOutOfRangeValueIndex+1:]...)
			self.items[i] = item
		}
	}
}

func (self *RunChart) renderItems(buf *Buffer, drawArea image.Rectangle) {

	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	for _, item := range self.items {

		xToPoint := make(map[int]image.Point)
		pointsOrder := make([]int, 0)

		for _, timePoint := range item.timePoints {

			if !self.isTimePointInRange(timePoint) {
				continue
			}

			timeDeltaWithGridMaxTime := self.grid.timeExtremum.max.Sub(timePoint.Time)
			deltaToPaddingRelation := float64(timeDeltaWithGridMaxTime.Nanoseconds()) / float64(self.grid.paddingDuration.Nanoseconds())
			x := self.grid.maxTimeWidth - (int(float64(self.grid.paddingWidth) * deltaToPaddingRelation))

			valuePerYDot := (self.grid.valueExtremum.max - self.grid.valueExtremum.min) / float64(drawArea.Dy()-1)
			y := int(float64(timePoint.Value-self.grid.valueExtremum.min) / valuePerYDot)

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

			//buf.SetCell(
			//	NewCell(self.DotRune, NewStyle(SelectColor(self.LineColors, 0))),
			//	currentPoint,
			//)

			canvas.Line(
				braillePoint(previousPoint),
				braillePoint(currentPoint),
				item.color,
			)
		}
	}

	canvas.Draw(buf)
}

func (self *RunChart) isTimePointInRange(point TimePoint) bool {
	return point.Time.After(self.grid.timeExtremum.min.Add(self.grid.paddingDuration))
}

func (self *RunChart) plotAxes(buf *Buffer) {
	// draw origin cell
	buf.SetCell(
		NewCell(BOTTOM_LEFT, NewStyle(ColorWhite)),
		image.Pt(self.Inner.Min.X+yAxisLabelsWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)

	// draw x axis line
	for i := yAxisLabelsWidth + 1; i < self.Inner.Dx(); i++ {
		buf.SetCell(
			NewCell(HORIZONTAL_DASH, NewStyle(ColorWhite)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}

	// draw grid
	for y := 0; y < self.Inner.Dy()-xAxisLabelsHeight-1; y = y + 2 {
		for x := 0; x < self.grid.linesCount; x++ {
			buf.SetCell(
				NewCell(VERTICAL_DASH, NewStyle(ColorDarkGrey)),
				image.Pt(self.grid.maxTimeWidth-x*self.grid.paddingWidth, y+self.Inner.Min.Y+1),
			)
		}
	}

	// draw y axis line
	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			NewCell(VERTICAL_DASH, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Min.X+yAxisLabelsWidth, i+self.Inner.Min.Y),
		)
	}

	// draw x axis time labels
	for i := 0; i < self.grid.linesCount; i++ {
		labelTime := self.grid.timeExtremum.max.Add(time.Duration(-i) * time.Second)
		buf.SetString(
			labelTime.Format("15:04:05"),
			NewStyle(ColorWhite),
			image.Pt(self.grid.maxTimeWidth-xAxisLabelsWidth/2-i*(self.grid.paddingWidth), self.Inner.Max.Y-1),
		)
	}

	// draw y axis labels
	verticalScale := self.grid.valueExtremum.max - self.grid.valueExtremum.min/float64(self.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 1; i*(yAxisLabelsGap+1) <= self.Inner.Dy()-1; i++ {
		buf.SetString(
			fmt.Sprintf("%.3f", float64(i)*self.grid.valueExtremum.min*verticalScale*(yAxisLabelsGap+1)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func GetValueExtremum(items []ChartItem) ValueExtremum {

	if len(items) == 0 {
		return ValueExtremum{0, 0}
	}

	var max, min = -math.MaxFloat64, math.MaxFloat64

	for _, item := range items {
		for _, point := range item.timePoints {
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

func GetTimeExtremum(linesCount int, paddingDuration time.Duration) TimeExtremum {
	maxTime := time.Now()
	return TimeExtremum{
		max: maxTime,
		min: maxTime.Add(-time.Duration(paddingDuration.Nanoseconds() * int64(linesCount))),
	}
}
