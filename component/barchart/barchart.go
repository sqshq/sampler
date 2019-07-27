package barchart

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	rw "github.com/mattn/go-runewidth"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/util"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"math"
)

const (
	barIndent int = 1
)

// BarChart presents categorical data with rectangular bars
type BarChart struct {
	*ui.Block
	*data.Consumer
	bars     []bar
	scale    int
	maxValue float64
	count    int64
	palette  console.Palette
}

type bar struct {
	label string
	color ui.Color
	value float64
	delta float64
}

func NewBarChart(c config.BarChartConfig, palette console.Palette) *BarChart {

	chart := BarChart{
		Block:    component.NewBlock(c.Title, true, palette),
		Consumer: data.NewConsumer(),
		bars:     []bar{},
		scale:    *c.Scale,
		maxValue: -math.MaxFloat64,
		palette:  palette,
	}

	for _, i := range c.Items {
		chart.addBar(*i.Label, *i.Color)
	}

	go func() {
		for {
			select {
			case sample := <-chart.SampleChannel:
				chart.consumeSample(sample)
			case alert := <-chart.AlertChannel:
				chart.Alert = alert
			}
		}
	}()

	return &chart
}

func (b *BarChart) consumeSample(sample *data.Sample) {

	b.count++

	float, err := util.ParseFloat(sample.Value)

	if err != nil {
		b.HandleConsumeFailure("Failed to parse a number", err, sample)
		return
	}

	b.HandleConsumeSuccess()

	index := -1
	for i, bar := range b.bars {
		if bar.label == sample.Label {
			index = i
		}
	}

	bar := b.bars[index]
	bar.delta = float - bar.value
	bar.value = float
	b.bars[index] = bar

	if float > b.maxValue {
		b.maxValue = float
	}

	// normalize bars height once in a while
	if b.count%500 == 0 {
		b.reselectMaxValue()
	}
}

func (b *BarChart) addBar(label string, color ui.Color) {
	b.bars = append(b.bars, bar{label: label, color: color, value: 0})
}

func (b *BarChart) reselectMaxValue() {
	maxValue := -math.MaxFloat64
	for _, bar := range b.bars {
		if bar.value > maxValue {
			maxValue = bar.value
		}
	}
	b.maxValue = maxValue
}

// Draw renders the barchart
func (b *BarChart) Draw(buffer *ui.Buffer) {
	b.Block.Draw(buffer)

	barWidth := int(math.Floor(float64(b.Inner.Dx()-len(b.bars)*barIndent) / float64(len(b.bars))))
	barXCoordinate := b.Inner.Min.X + barIndent

	labelStyle := ui.NewStyle(b.palette.BaseColor)

	for _, bar := range b.bars {

		// draw bar
		height := int((bar.value / b.maxValue) * float64(b.Inner.Dy()-1))
		if height <= 1 {
			height = 2
		}

		maxYCoordinate := b.Inner.Max.Y - height
		for x := barXCoordinate; x < ui.MinInt(barXCoordinate+barWidth, b.Inner.Max.X-barIndent); x++ {
			for y := b.Inner.Max.Y - 2; y >= maxYCoordinate; y-- {
				c := ui.NewCell(console.SymbolHorizontalBar, ui.NewStyle(bar.color))
				buffer.SetCell(c, image.Pt(x, y))
			}
		}

		// draw label
		labelXCoordinate := barXCoordinate +
			int(float64(barWidth)/2) -
			int(float64(rw.StringWidth(bar.label))/2)
		buffer.SetString(
			bar.label,
			labelStyle,
			image.Pt(labelXCoordinate, b.Inner.Max.Y-1))

		// draw value & delta
		value := util.FormatValue(bar.value, b.scale)
		if bar.delta != 0 {
			value = fmt.Sprintf("%s/%s", value, util.FormatDelta(bar.delta, b.scale))
		}
		valueXCoordinate := barXCoordinate +
			int(float64(barWidth)/2) -
			int(float64(rw.StringWidth(value))/2)
		buffer.SetString(
			value,
			labelStyle,
			image.Pt(valueXCoordinate, maxYCoordinate-1))

		barXCoordinate += barWidth + barIndent
	}

	component.RenderAlert(b.Alert, b.Rectangle, buffer)
}
