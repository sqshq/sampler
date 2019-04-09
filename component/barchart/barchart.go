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
	"strconv"
)

const (
	barIndent int = 1
)

type BarChart struct {
	*ui.Block
	*data.Consumer
	alert    *data.Alert
	bars     []Bar
	scale    int
	maxValue float64
	count    int64
	palette  console.Palette
}

type Bar struct {
	label string
	color ui.Color
	value float64
	delta float64
}

func NewBarChart(c config.BarChartConfig, palette console.Palette) *BarChart {

	chart := BarChart{
		Block:    component.NewBlock(c.Title, true, palette),
		Consumer: data.NewConsumer(),
		bars:     []Bar{},
		scale:    *c.Scale,
		maxValue: -math.MaxFloat64,
		palette:  palette,
	}

	for _, i := range c.Items {
		chart.AddBar(*i.Label, *i.Color)
	}

	go func() {
		for {
			select {
			case sample := <-chart.SampleChannel:
				chart.consumeSample(sample)
			case alert := <-chart.AlertChannel:
				chart.alert = alert
			}
		}
	}()

	return &chart
}

func (b *BarChart) consumeSample(sample *data.Sample) {

	b.count++

	float, err := util.ParseFloat(sample.Value)
	if err != nil {
		b.AlertChannel <- &data.Alert{
			Title: "FAILED TO PARSE A NUMBER",
			Text:  err.Error(),
			Color: sample.Color,
		}
		return
	}

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

func (b *BarChart) AddBar(label string, color ui.Color) {
	b.bars = append(b.bars, Bar{label: label, color: color, value: 0})
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
		value := formatValue(bar.value, b.scale)
		if bar.delta != 0 {
			value = fmt.Sprintf("%s / %s", value, formatValueWithSign(bar.delta, b.scale))
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

	component.RenderAlert(b.alert, b.Rectangle, buffer)
}

// TODO extract to utils
func formatValue(value float64, scale int) string {
	if math.Abs(value) == math.MaxFloat64 {
		return "Inf"
	} else {
		format := "%." + strconv.Itoa(scale) + "f"
		return fmt.Sprintf(format, value)
	}
}

// TODO extract to utils
func formatValueWithSign(value float64, scale int) string {
	if value == 0 {
		return " 0"
	} else if value > 0 {
		return "+" + formatValue(value, scale)
	} else {
		return formatValue(value, scale)
	}
}
