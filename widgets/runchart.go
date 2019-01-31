package widgets

import (
	"fmt"
	"image"
	"log"
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
	DataLabels []string
	LineColors []Color
	timePoints []TimePoint
	dataMutex  *sync.Mutex
	grid       ChartGrid
}

type Item struct {
	timePoints []TimePoint
	color      Color
	label      string
}

type TimePoint struct {
	Value float64
	Time  time.Time
}

func NewRunChart(title string) *RunChart {
	block := *NewBlock()
	block.Title = title
	//self.LineColors[0] = ui.ColorYellow
	return &RunChart{
		Block:      block,
		LineColors: Theme.Plot.Lines,
		timePoints: make([]TimePoint, 0),
		dataMutex:  &sync.Mutex{},
	}
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

func (self *RunChart) newChartGrid() ChartGrid {

	linesCount := (self.Inner.Max.X - self.Inner.Min.X) / (xAxisLabelsGap + xAxisLabelsWidth)
	paddingDuration := time.Duration(time.Second) // TODO support others and/or adjust automatically depending on refresh rate

	return ChartGrid{
		linesCount:      linesCount,
		paddingDuration: paddingDuration,
		paddingWidth:    xAxisLabelsGap + xAxisLabelsWidth,
		maxTimeWidth:    self.Inner.Max.X - xAxisLabelsWidth/2 - xAxisLabelsGap,
		timeExtremum:    GetTimeExtremum(linesCount, paddingDuration),
		valueExtremum:   GetValueExtremum(self.timePoints),
	}
}

func (self *RunChart) Draw(buf *Buffer) {

	self.dataMutex.Lock()
	self.Block.Draw(buf)
	self.grid = self.newChartGrid()
	self.plotAxes(buf)

	drawArea := image.Rect(
		self.Inner.Min.X+yAxisLabelsWidth+1, self.Inner.Min.Y,
		self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	self.renderBraille(buf, drawArea)
	self.dataMutex.Unlock()
}

func (self *RunChart) ConsumeValue(value string, label string) {

	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatalf("Expected float number, but got %v", value) // TODO visual notification
	}

	self.dataMutex.Lock()
	self.timePoints = append(self.timePoints, TimePoint{Value: float, Time: time.Now()})
	self.trimOutOfRangeValues()
	self.dataMutex.Unlock()
}

func (self *RunChart) ConsumeError(err error) {
	// TODO visual notification
}

func (self *RunChart) trimOutOfRangeValues() {

	lastOutOfRangeValueIndex := -1

	for i, timePoint := range self.timePoints {
		if !self.isTimePointInRange(timePoint) {
			lastOutOfRangeValueIndex = i
		}
	}

	if lastOutOfRangeValueIndex > 0 {
		self.timePoints = append(self.timePoints[:0], self.timePoints[lastOutOfRangeValueIndex+1:]...)
	}
}

func (self *RunChart) renderBraille(buf *Buffer, drawArea image.Rectangle) {

	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	pointPerX := make(map[int]image.Point)
	pointsOrder := make([]int, 0)

	for _, timePoint := range self.timePoints {

		if !self.isTimePointInRange(timePoint) {
			continue
		}

		timeDeltaWithGridMaxTime := self.grid.timeExtremum.max.Sub(timePoint.Time)
		deltaToPaddingRelation := float64(timeDeltaWithGridMaxTime.Nanoseconds()) / float64(self.grid.paddingDuration.Nanoseconds())
		x := self.grid.maxTimeWidth - (int(float64(self.grid.paddingWidth) * deltaToPaddingRelation))

		valuePerYDot := (self.grid.valueExtremum.max - self.grid.valueExtremum.min) / float64(drawArea.Dy()-1)
		y := int(float64(timePoint.Value-self.grid.valueExtremum.min) / valuePerYDot)

		if _, exists := pointPerX[x]; exists {
			continue
		}

		pointPerX[x] = image.Pt(x, drawArea.Max.Y-y-1)
		pointsOrder = append(pointsOrder, x)
	}

	for i, x := range pointsOrder {

		currentPoint := pointPerX[x]
		var previousPoint image.Point

		if i == 0 {
			previousPoint = currentPoint
		} else {
			previousPoint = pointPerX[pointsOrder[i-1]]
		}

		//buf.SetCell(
		//	NewCell(self.DotRune, NewStyle(SelectColor(self.LineColors, 0))),
		//	currentPoint,
		//)

		canvas.Line(
			braillePoint(previousPoint),
			braillePoint(currentPoint),
			SelectColor(self.LineColors, 0), //i
		)
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

func GetValueExtremum(points []TimePoint) ValueExtremum {

	if len(points) == 0 {
		return ValueExtremum{0, 0}
	}

	var max, min = points[0], points[0]

	for _, point := range points {
		if point.Value > max.Value {
			max = point
		}
		if point.Value < min.Value {
			min = point
		}
	}

	return ValueExtremum{max: max.Value, min: min.Value}
}

func GetTimeExtremum(linesCount int, paddingDuration time.Duration) TimeExtremum {
	maxTime := time.Now()
	return TimeExtremum{
		max: maxTime,
		min: maxTime.Add(-time.Duration(paddingDuration.Nanoseconds() * int64(linesCount))),
	}
}
