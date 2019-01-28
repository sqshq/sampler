package widgets

import (
	"fmt"
	"image"
	"sync"
	"time"

	. "github.com/sqshq/termui"
)

type TimePlot struct {
	Block
	DataLabels        []string
	MaxValueTimePoint TimePoint
	LineColors        []Color
	ShowAxes          bool
	DotRune           rune
	HorizontalScale   int
	Marker            PlotMarker
	timePoints        []TimePoint
	dataMutex         *sync.Mutex
	grid              PlotGrid
}

const (
	xAxisLabelsHeight = 1
	xAxisLabelsWidth  = 8
	xAxisLabelsGap    = 2
	yAxisLabelsWidth  = 5
	yAxisLabelsGap    = 1
)

type TimePoint struct {
	Value float64
	Time  time.Time
}

type PlotMarker uint

const (
	MarkerBraille PlotMarker = iota
	MarkerDot
)

func NewTimePlot() *TimePlot {
	return &TimePlot{
		Block:           *NewBlock(),
		LineColors:      Theme.Plot.Lines,
		DotRune:         DOT,
		HorizontalScale: 1,
		ShowAxes:        true,
		Marker:          MarkerBraille,
		timePoints:      make([]TimePoint, 0),
		dataMutex:       &sync.Mutex{},
	}
}

type PlotGrid struct {
	count           int
	maxTimeX        int
	maxTime         time.Time
	minTime         time.Time
	maxValue        float64
	minValue        float64
	spacingDuration time.Duration
	spacingWidth    int
}

func (self *TimePlot) newPlotGrid() PlotGrid {

	count := (self.Inner.Max.X - self.Inner.Min.X) / (xAxisLabelsGap + xAxisLabelsWidth)
	spacingDuration := time.Duration(time.Second) // TODO support others and/or adjust automatically depending on refresh rate
	maxTime := time.Now()
	minTime := maxTime.Add(-time.Duration(spacingDuration.Nanoseconds() * int64(count)))
	maxPoint, minPoint := GetMaxAndMinValueTimePoints(self.timePoints)

	return PlotGrid{
		count:           count,
		spacingDuration: spacingDuration,
		spacingWidth:    xAxisLabelsGap + xAxisLabelsWidth,
		maxTimeX:        self.Inner.Max.X - xAxisLabelsWidth/2 - xAxisLabelsGap,
		maxTime:         maxTime,
		minTime:         minTime,
		maxValue:        maxPoint.Value,
		minValue:        minPoint.Value,
	}
}

func (self *TimePlot) Draw(buf *Buffer) {

	self.dataMutex.Lock()
	self.Block.Draw(buf)
	self.grid = self.newPlotGrid()

	if self.ShowAxes {
		self.plotAxes(buf)
	}

	drawArea := self.Inner
	if self.ShowAxes {
		drawArea = image.Rect(
			self.Inner.Min.X+yAxisLabelsWidth+1, self.Inner.Min.Y,
			self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}

	self.renderBraille(buf, drawArea)
	self.dataMutex.Unlock()
}

func (self *TimePlot) AddValue(value float64) {
	self.dataMutex.Lock()
	self.timePoints = append(self.timePoints, TimePoint{Value: value, Time: time.Now()})
	self.trimOutOfRangeValues()
	self.dataMutex.Unlock()
}

func (self *TimePlot) trimOutOfRangeValues() {

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

func (self *TimePlot) renderBraille(buf *Buffer, drawArea image.Rectangle) {

	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	pointPerX := make(map[int]image.Point)
	pointsOrder := make([]int, 0)

	for _, timePoint := range self.timePoints {

		if !self.isTimePointInRange(timePoint) {
			continue
		}

		timeDeltaWithGridMaxTime := self.grid.maxTime.Sub(timePoint.Time)
		deltaToSpacingRelation := float64(timeDeltaWithGridMaxTime.Nanoseconds()) / float64(self.grid.spacingDuration.Nanoseconds())
		x := self.grid.maxTimeX - (int(float64(self.grid.spacingWidth) * deltaToSpacingRelation))

		valuePerYDot := (self.grid.maxValue - self.grid.minValue) / float64(drawArea.Dy()-1)
		y := int(float64(timePoint.Value-self.grid.minValue) / valuePerYDot)

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
			braille(previousPoint),
			braille(currentPoint),
			SelectColor(self.LineColors, 0), //i
		)
	}

	canvas.Draw(buf)
}

func (self *TimePlot) isTimePointInRange(point TimePoint) bool {
	return point.Time.After(self.grid.minTime.Add(self.grid.spacingDuration))
}

func (self *TimePlot) plotAxes(buf *Buffer) {
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
		for x := 0; x < self.grid.count; x++ {
			buf.SetCell(
				NewCell(VERTICAL_DASH, NewStyle(ColorDarkGrey)),
				image.Pt(self.grid.maxTimeX-x*self.grid.spacingWidth, y+self.Inner.Min.Y+1),
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
	for i := 0; i < self.grid.count; i++ {
		labelTime := self.grid.maxTime.Add(time.Duration(-i) * time.Second)
		buf.SetString(
			labelTime.Format("15:04:05"),
			NewStyle(ColorWhite),
			image.Pt(self.grid.maxTimeX-xAxisLabelsWidth/2-i*(self.grid.spacingWidth), self.Inner.Max.Y-1),
		)
	}

	// draw y axis labels
	verticalScale := self.grid.maxValue - self.grid.minValue/float64(self.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 1; i*(yAxisLabelsGap+1) <= self.Inner.Dy()-1; i++ {
		buf.SetString(
			fmt.Sprintf("%.3f", float64(i)*self.grid.minValue*verticalScale*(yAxisLabelsGap+1)),
			NewStyle(ColorWhite),
			image.Pt(self.Inner.Min.X, self.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func GetMaxAndMinValueTimePoints(points []TimePoint) (TimePoint, TimePoint) {

	if len(points) == 0 {
		return TimePoint{0, time.Now()}, TimePoint{0, time.Now()}
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

	return max, min
}
