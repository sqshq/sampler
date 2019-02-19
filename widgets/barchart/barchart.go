package barchart

import (
	rw "github.com/mattn/go-runewidth"
	"github.com/sqshq/sampler/data"
	ui "github.com/sqshq/termui"
	"image"
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
		maxValue: 0,
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

	bar := b.bars[index]
	bar.value = float
	b.bars[index] = bar
}

func (b *BarChart) Draw(buf *ui.Buffer) {
	b.Block.Draw(buf)

	maxVal := b.maxValue

	barXCoordinate := b.Inner.Min.X

	for i, data := range b.Data {
		// draw bar
		height := int((data / maxVal) * float64(b.Inner.Dy()-1))
		for x := barXCoordinate; x < ui.MinInt(barXCoordinate+b.BarWidth, b.Inner.Max.X); x++ {
			for y := b.Inner.Max.Y - 2; y > (b.Inner.Max.Y-2)-height; y-- {
				c := ui.NewCell(barSymbol, ui.NewStyle(ui.SelectColor(b.BarColors, i)))
				buf.SetCell(c, image.Pt(x, y))
			}
		}

		// draw label
		if i < len(b.Labels) {
			labelXCoordinate := barXCoordinate +
				int((float64(b.BarWidth) / 2)) -
				int((float64(rw.StringWidth(b.Labels[i])) / 2))
			buf.SetString(
				b.Labels[i],
				ui.SelectStyle(b.LabelStyles, i),
				image.Pt(labelXCoordinate, b.Inner.Max.Y-1),
			)
		}

		// draw value
		numberXCoordinate := barXCoordinate + int((float64(b.BarWidth) / 2))
		if numberXCoordinate <= b.Inner.Max.X {
			buf.SetString(
				b.NumFmt(data),
				ui.NewStyle(
					ui.SelectStyle(b.NumStyles, i+1).Fg,
					ui.SelectColor(b.BarColors, i),
					ui.SelectStyle(b.NumStyles, i+1).Modifier,
				),
				image.Pt(numberXCoordinate, b.Inner.Max.Y-2),
			)
		}

		barXCoordinate += (b.BarWidth + b.BarGap)
	}
}
