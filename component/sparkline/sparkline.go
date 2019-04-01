package sparkline

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/util"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"strconv"
)

type SparkLine struct {
	*ui.Block
	*data.Consumer
	alert    *data.Alert
	values   []float64
	maxValue float64
	minValue float64
	scale    int
	gradient []ui.Color
	palette  console.Palette
}

func NewSparkLine(c config.SparkLineConfig, palette console.Palette) *SparkLine {

	line := &SparkLine{
		Block:    component.NewBlock(c.Title, true, palette),
		Consumer: data.NewConsumer(),
		values:   []float64{},
		scale:    *c.Scale,
		gradient: *c.Gradient,
		palette:  palette,
	}

	go func() {
		for {
			select {
			case sample := <-line.SampleChannel:
				line.consumeSample(sample)
			case alert := <-line.AlertChannel:
				line.alert = alert
			}
		}
	}()

	return line
}

func (s *SparkLine) consumeSample(sample *data.Sample) {
	float, err := strconv.ParseFloat(sample.Value, 64)

	if err != nil {
		s.AlertChannel <- &data.Alert{
			Title: "FAILED TO PARSE A NUMBER",
			Text:  err.Error(),
			Color: sample.Color,
		}
		return
	}

	s.values = append(s.values, float)
	max, min := s.values[0], s.values[0]

	for i := len(s.values) - 1; i >= 0; i-- {
		if len(s.values)-i > s.Dx() {
			break
		}
		if s.values[i] > max {
			max = s.values[i]
		}
		if s.values[i] < min {
			min = s.values[i]
		}
	}

	s.maxValue = max
	s.minValue = min

	// perform cleanup once in a while
	if len(s.values)%100 == 0 {
		s.values = append(s.values[:0], s.values[len(s.values)-s.Dx()+1:]...)
	}
}

func (s *SparkLine) Draw(buffer *ui.Buffer) {

	textStyle := ui.NewStyle(s.palette.BaseColor)

	height := s.Dy() - 2
	minValue := util.FormatValue(s.minValue, s.scale)
	maxValue := util.FormatValue(s.maxValue, s.scale)
	curValue := util.FormatValue(s.values[len(s.values)-1], s.scale)

	indent := 2 + util.Max([]int{
		len(minValue), len(maxValue), len(curValue),
	})

	for i := len(s.values) - 1; i >= 0; i-- {

		n := len(s.values) - i

		if n > s.Dx()-indent-3 {
			break
		}

		top := int((s.values[i] / s.maxValue) * float64(height))

		if top == 0 {
			top = 1
		}

		for j := 1; j <= top; j++ {
			buffer.SetCell(ui.NewCell(console.SymbolVerticalBar, ui.NewStyle(console.GetGradientColor(s.gradient, j-1, height))), image.Pt(s.Inner.Max.X-n-indent, s.Inner.Max.Y-j))
			if i == len(s.values)-1 && j == top {
				buffer.SetString(curValue, textStyle, image.Pt(s.Inner.Max.X-n-indent+2, s.Inner.Max.Y-j))
				buffer.SetString(minValue, textStyle, image.Pt(s.Inner.Max.X-n-indent+2, s.Max.Y-2))
				buffer.SetString(maxValue, textStyle, image.Pt(s.Inner.Max.X-n-indent+2, s.Min.Y+1))
			}
		}
	}

	s.Block.Draw(buffer)
}
