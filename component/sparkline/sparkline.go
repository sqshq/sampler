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
	color    ui.Color
	palette  console.Palette
}

func NewSparkLine(c config.SparkLineConfig, palette console.Palette) *SparkLine {

	line := &SparkLine{
		Block:    component.NewBlock(c.Title, true, palette),
		Consumer: data.NewConsumer(),
		values:   []float64{},
		scale:    *c.Scale,
		color:    *c.Item.Color,
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
			Title: "FAILED TO PARSE NUMBER",
			Text:  err.Error(),
			Color: sample.Color,
		}
		return
	}

	s.values = append(s.values, float)
	// TODO cleanup old ones

	for i := len(s.values) - 1; i >= 0; i-- {
		if len(s.values)-i > s.Dx() {
			break
		}
		if s.values[i] > s.maxValue {
			s.maxValue = s.values[i]
		}
		if s.values[i] < s.minValue {
			s.minValue = s.values[i]
		}
	}
}

// TODO make sure that 0 value is still printed
// TODO make sure that cur value is printed on the same Y as sparkline (include in for loop for last iteratiton)
// TODO gradient color
func (s *SparkLine) Draw(buffer *ui.Buffer) {

	textStyle := ui.NewStyle(s.palette.BaseColor)
	lineStyle := ui.NewStyle(s.color)

	minValue := util.FormatValue(s.minValue, s.scale)
	maxValue := util.FormatValue(s.maxValue, s.scale)
	curValue := util.FormatValue(s.values[len(s.values)-1], s.scale)

	buffer.SetString(minValue, textStyle, image.Pt(s.Min.X+2, s.Max.Y-2))
	buffer.SetString(maxValue, textStyle, image.Pt(s.Min.X+2, s.Min.Y+2))

	curY := int((s.values[len(s.values)-1]/s.maxValue)*float64(s.Dy())) - 1
	buffer.SetString(curValue, textStyle, image.Pt(s.Max.X-len(curValue)-2, s.Max.Y-util.Max([]int{curY, 2})))

	indent := 2 + util.Max([]int{
		len(minValue), len(maxValue), len(curValue),
	})

	for i := len(s.values) - 1; i >= 0; i-- {

		n := len(s.values) - i

		if n > s.Dx()-indent*2-2 {
			break
		}

		for j := 1; j < int((s.values[i]/s.maxValue)*float64(s.Dy()-2))+2; j++ {
			buffer.SetString("▪", lineStyle, image.Pt(s.Inner.Max.X-n-indent, s.Inner.Max.Y-j))
		}
	}

	s.Block.Draw(buffer)
}
