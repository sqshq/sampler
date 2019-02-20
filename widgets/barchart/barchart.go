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
	barSymbol rune = '⠿'
	barIndent int  = 1
)

type BarChart struct {
	ui.Block
	bars     []Bar
	scale    int
	maxValue float64
}

type Bar struct {
	label string
	color ui.Color
	value float64
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

	if float > b.maxValue {
		b.maxValue = float
	}

	bar := b.bars[index]
	bar.value = float
	b.bars[index] = bar
}

func (b *BarChart) Draw(buf *ui.Buffer) {
	b.Block.Draw(buf)

	barWidth := b.Inner.Dx() / len(b.bars)
	barXCoordinate := b.Inner.Min.X

	labelStyle := ui.NewStyle(console.ColorWhite)

	for _, bar := range b.bars {
		// draw bar
		height := int((bar.value / b.maxValue) * float64(b.Inner.Dy()-1))
		for x := barXCoordinate; x < ui.MinInt(barXCoordinate+barWidth, b.Inner.Max.X); x++ {
			for y := b.Inner.Max.Y - 2; y > (b.Inner.Max.Y-2)-height; y-- {
				c := ui.NewCell(barSymbol, ui.NewStyle(bar.color))
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
			image.Pt(labelXCoordinate, b.Inner.Max.Y-1),
		)

		// draw value
		numberXCoordinate := barXCoordinate + int(float64(barWidth)/2)
		if numberXCoordinate <= b.Inner.Max.X {
			buf.SetString(
				formatValue(bar.value, b.scale),
				labelStyle,
				image.Pt(numberXCoordinate, b.Inner.Max.Y-2),
			)
		}

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
