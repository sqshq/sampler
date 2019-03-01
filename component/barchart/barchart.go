package barchart

import (
	"fmt"
	rw "github.com/mattn/go-runewidth"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	ui "github.com/sqshq/termui"
	"image"
	"math"
	"strconv"
)

const (
	barIndent int = 1
)

type BarChart struct {
	ui.Block
	bars     []Bar
	scale    int
	maxValue float64
	count    int64
}

type Bar struct {
	label string
	color ui.Color
	value float64
	delta float64
}

func NewBarChart(title string, scale int) *BarChart {
	block := *ui.NewBlock()
	block.Title = title
	return &BarChart{
		Block:    block,
		bars:     []Bar{},
		scale:    scale,
		maxValue: -math.MaxFloat64,
	}
}

func (b *BarChart) AddBar(label string, color ui.Color) {
	b.bars = append(b.bars, Bar{label: label, color: color, value: 0})
}

func (b *BarChart) ConsumeSample(sample data.Sample) {

	b.count++

	float, err := strconv.ParseFloat(sample.Value, 64)
	if err != nil {
		// TODO visual notification + check sample.Error
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

func (b *BarChart) reselectMaxValue() {
	maxValue := -math.MaxFloat64
	for _, bar := range b.bars {
		if bar.value > maxValue {
			maxValue = bar.value
		}
	}
	b.maxValue = maxValue
}

func (b *BarChart) Draw(buf *ui.Buffer) {
	b.Block.Draw(buf)

	barWidth := int(math.Ceil(float64(b.Inner.Dx()-2*barIndent-len(b.bars)*barIndent) / float64(len(b.bars))))
	barXCoordinate := b.Inner.Min.X + barIndent

	labelStyle := ui.NewStyle(console.ColorWhite)

	for _, bar := range b.bars {

		// draw bar
		height := int((bar.value / b.maxValue) * float64(b.Inner.Dy()-1))
		if height <= 1 {
			height = 2
		}

		maxYCoordinate := b.Inner.Max.Y - height
		for x := barXCoordinate; x < ui.MinInt(barXCoordinate+barWidth, b.Inner.Max.X-barIndent); x++ {
			for y := b.Inner.Max.Y - 2; y >= maxYCoordinate; y-- {
				c := ui.NewCell(console.SymbolShade, ui.NewStyle(bar.color))
				buf.SetCell(c, image.Pt(x, y))
			}
		}

		// draw label
		labelXCoordinate := barXCoordinate +
			int(float64(barWidth)/2) -
			int(float64(rw.StringWidth(bar.label))/2)
		buf.SetString(
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
		buf.SetString(
			value,
			labelStyle,
			image.Pt(valueXCoordinate, maxYCoordinate-1))

		barXCoordinate += barWidth + barIndent
	}
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
