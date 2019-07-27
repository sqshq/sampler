package runchart

import (
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/util"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"math"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
)

const (
	xAxisGridWidth     = xAxisLabelsIndent + xAxisLabelsWidth
	xAxisLabelsHeight  = 1
	xAxisLabelsWidth   = 8
	xAxisLabelsIndent  = 2
	yAxisLabelsHeight  = 1
	yAxisLabelsIndent  = 1
	historyReserveMin  = 2
	xBrailleMultiplier = 2
	yBrailleMultiplier = 4
)

type Mode int

const (
	ModeDefault  Mode = 0
	ModePinpoint Mode = 1
)

const (
	CommandDisableSelection = "DISABLE_SELECTION"
	CommandMoveSelection    = "MOVE_SELECTION"
)

// RunChart displays observed data in a time sequence
type RunChart struct {
	*ui.Block
	*data.Consumer
	lines     []TimeLine
	grid      chartGrid
	timescale time.Duration
	mutex     *sync.Mutex
	mode      Mode
	selection time.Time
	scale     int
	legend    legend
	palette   console.Palette
}

type TimePoint struct {
	value      float64
	time       time.Time
	coordinate int
}

type TimeLine struct {
	points              []TimePoint
	extrema             ValueExtrema
	color               ui.Color
	label               string
	selectionCoordinate int
	selectionPoint      TimePoint
}

type TimeRange struct {
	max time.Time
	min time.Time
}

type ValueExtrema struct {
	max float64
	min float64
}

func NewRunChart(c config.RunChartConfig, palette console.Palette) *RunChart {

	chart := RunChart{
		Block:     component.NewBlock(c.Title, true, palette),
		Consumer:  data.NewConsumer(),
		lines:     []TimeLine{},
		timescale: calculateTimescale(*c.RateMs),
		mutex:     &sync.Mutex{},
		scale:     *c.Scale,
		mode:      ModeDefault,
		legend:    legend{Enabled: c.Legend.Enabled, Details: c.Legend.Details},
		palette:   palette,
	}

	for _, i := range c.Items {
		chart.AddLine(*i.Label, *i.Color)
	}

	go func() {
		for {
			select {
			case sample := <-chart.SampleChannel:
				chart.consumeSample(sample)
			case alert := <-chart.AlertChannel:
				chart.Alert = alert
			case command := <-chart.CommandChannel:
				switch command.Type {
				case CommandDisableSelection:
					chart.disableSelection()
				case CommandMoveSelection:
					chart.moveSelection(command.Value.(int))
				}
			}
		}
	}()

	return &chart
}

func (c *RunChart) newTimePoint(value float64) TimePoint {
	now := time.Now()
	return TimePoint{
		value:      value,
		time:       now,
		coordinate: c.calculateTimeCoordinate(now),
	}
}

func (c *RunChart) Draw(buffer *ui.Buffer) {

	c.mutex.Lock()
	c.Block.Draw(buffer)
	c.grid = c.newChartGrid()

	drawArea := image.Rect(
		c.Inner.Min.X+c.grid.minTimeWidth+2, c.Inner.Min.Y,
		c.Inner.Max.X, c.Inner.Max.Y-xAxisLabelsHeight-1,
	)

	c.renderAxes(buffer)
	c.renderLines(buffer, drawArea)
	c.renderLegend(buffer, drawArea)
	component.RenderAlert(c.Alert, c.Rectangle, buffer)
	c.mutex.Unlock()
}

func (c *RunChart) AddLine(Label string, color ui.Color) {
	line := TimeLine{
		points:  []TimePoint{},
		color:   color,
		label:   Label,
		extrema: ValueExtrema{max: -math.MaxFloat64, min: math.MaxFloat64},
	}
	c.lines = append(c.lines, line)
}

func (c *RunChart) consumeSample(sample *data.Sample) {

	float, err := util.ParseFloat(sample.Value)
	if err != nil {
		c.HandleConsumeFailure("Failed to parse a number", err, sample)
		return
	}

	c.HandleConsumeSuccess()

	c.mutex.Lock()

	index := -1
	for i, line := range c.lines {
		if line.label == sample.Label {
			index = i
		}
	}

	line := c.lines[index]

	if float < line.extrema.min {
		line.extrema.min = float
	}
	if float > line.extrema.max {
		line.extrema.max = float
	}

	line.points = append(line.points, c.newTimePoint(float))
	c.lines[index] = line

	if len(line.points)%100 == 0 {
		c.trimOutOfRangeValues()
	}

	c.mutex.Unlock()
}

func (c *RunChart) renderLines(buffer *ui.Buffer, drawArea image.Rectangle) {

	canvas := ui.NewCanvas()
	canvas.Rectangle = drawArea

	if len(c.lines) == 0 || len(c.lines[0].points) == 0 {
		return
	}

	selectionCoordinate := c.calculateTimeCoordinate(c.selection)
	selectionPoints := make(map[int]image.Point)

	probe := c.lines[0].points[0]
	probeCalculatedCoordinate := c.calculateTimeCoordinate(probe.time)
	delta := probe.coordinate - probeCalculatedCoordinate

	for i, line := range c.lines {

		xPoint := make(map[int]image.Point)
		xOrder := make([]int, 0)

		// move selection on a delta, if it was instantiated after cursor move
		if line.selectionCoordinate != 0 {
			line.selectionCoordinate -= delta
			c.lines[i].selectionCoordinate = line.selectionCoordinate
		}

		for j, timePoint := range line.points {

			timePoint.coordinate -= delta
			line.points[j] = timePoint

			var y int
			if c.grid.valueExtrema.max == c.grid.valueExtrema.min {
				y = (drawArea.Dy() - 2) / 2
			} else {
				valuePerY := (c.grid.valueExtrema.max - c.grid.valueExtrema.min) / float64(drawArea.Dy()-2)
				y = int(float64(timePoint.value-c.grid.valueExtrema.min) / valuePerY)
			}

			point := image.Pt(timePoint.coordinate, drawArea.Max.Y-y-1)

			if _, exists := xPoint[point.X]; exists {
				continue
			}

			if !point.In(drawArea) {
				continue
			}

			if line.selectionCoordinate == 0 {
				// instantiate selection coordinate as the closest point to the cursor time
				if len(line.points) > j+1 && ui.AbsInt(timePoint.coordinate-selectionCoordinate) > ui.AbsInt(line.points[j+1].coordinate-selectionCoordinate) {
					selectionPoints[i] = point
					c.lines[i].selectionPoint = timePoint
				}
			} else if timePoint.coordinate == line.selectionCoordinate {
				selectionPoints[i] = point
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

			canvas.SetLine(
				braillePoint(previousPoint),
				braillePoint(currentPoint),
				line.color,
			)
		}
	}

	canvas.Draw(buffer)

	if c.mode == ModePinpoint {
		for lineIndex, point := range selectionPoints {
			buffer.SetCell(ui.NewCell(console.SymbolSelection, ui.NewStyle(c.lines[lineIndex].color)), point)
			if c.lines[lineIndex].selectionCoordinate == 0 {
				c.lines[lineIndex].selectionCoordinate = point.X
			}
		}
	}
}

func (c *RunChart) trimOutOfRangeValues() {

	minRangeTime := c.grid.timeRange.min.Add(-time.Minute * time.Duration(historyReserveMin))

	for i, item := range c.lines {
		lastOutOfRangeValueIndex := -1

		for j, point := range item.points {
			if point.time.Before(minRangeTime) {
				lastOutOfRangeValueIndex = j
			}
		}

		if lastOutOfRangeValueIndex > 0 {
			item.points = append(item.points[:0], item.points[lastOutOfRangeValueIndex+1:]...)
			c.lines[i] = item
		}
	}
}

func (c *RunChart) calculateTimeCoordinate(t time.Time) int {
	timeDeltaWithGridMaxTime := c.grid.timeRange.max.Sub(t).Nanoseconds()
	timeDeltaToPaddingRelation := float64(timeDeltaWithGridMaxTime) / float64(c.timescale.Nanoseconds())
	return c.grid.maxTimeWidth - int(math.Ceil(float64(xAxisGridWidth)*timeDeltaToPaddingRelation))
}

func (c *RunChart) moveSelection(shift int) {

	if c.mode == ModeDefault {
		c.mode = ModePinpoint
		c.selection = getMidRangeTime(c.grid.timeRange)
		return
	}

	c.selection = c.selection.Add(c.grid.timePerPoint * time.Duration(shift))
	if c.selection.After(c.grid.timeRange.max) {
		c.selection = c.grid.timeRange.max
	} else if c.selection.Before(c.grid.timeRange.min) {
		c.selection = c.grid.timeRange.min
	}

	for i := range c.lines {
		c.lines[i].selectionCoordinate = 0
	}
}

func (c *RunChart) disableSelection() {
	if c.mode == ModePinpoint {
		c.mode = ModeDefault
		return
	}
}

func getMidRangeTime(r TimeRange) time.Time {
	delta := r.max.Sub(r.min)
	return r.max.Add(-delta / 2)
}

// time duration between grid lines
func calculateTimescale(rateMs int) time.Duration {

	multiplier := rateMs * xAxisGridWidth / 2
	timescale := time.Duration(time.Millisecond * time.Duration(multiplier)).Round(time.Second)

	if timescale.Seconds() == 0 {
		return time.Second
	}

	return timescale
}

func braillePoint(point image.Point) image.Point {
	return image.Point{X: point.X * xBrailleMultiplier, Y: point.Y * yBrailleMultiplier}
}
