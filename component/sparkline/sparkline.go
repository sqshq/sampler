package sparkline

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/util"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"sync"
)

// SparkLine displays general shape of a measurement variation over time
type SparkLine struct {
	*ui.Block
	*data.Consumer
	label    string
	values   []float64
	maxValue float64
	minValue float64
	minFloor *float64
	maxFloor *float64
	scale    int
	gradient []ui.Color
	palette  console.Palette
	mutex    *sync.Mutex
}

func NewSparkLine(c config.SparkLineConfig, palette console.Palette) *SparkLine {

	line := &SparkLine{
		Block:    component.NewComponentBlock(c.Title, c.IsTitleVisible(), *c.Border, palette),
		Consumer: data.NewConsumer(),
		label:    *c.Item.Label,
		values:   []float64{},
		scale:    *c.Scale,
		minFloor: c.Min,
		maxFloor: c.Max,
		gradient: *c.Gradient,
		palette:  palette,
		mutex:    &sync.Mutex{},
	}

	go func() {
		for {
			select {
			case sample := <-line.SampleChannel:
				line.consumeSample(sample)
			case alert := <-line.AlertChannel:
				line.Alert = alert
			}
		}
	}()

	return line
}

func (s *SparkLine) consumeSample(sample *data.Sample) {

	float, err := util.ParseFloat(sample.Value)
	if err != nil {
		s.HandleConsumeFailure("Failed to parse a number", err, sample)
		return
	}

	s.HandleConsumeSuccess()

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

	if s.maxFloor != nil && *s.maxFloor > max {
		max = *s.maxFloor
	}
	if s.minFloor != nil && *s.minFloor < min {
		min = *s.minFloor
	}

	s.maxValue = max
	s.minValue = min

	if len(s.values)%100 == 0 {
		s.mutex.Lock()
		s.trimOutOfRangeValues(s.Dx())
		s.mutex.Unlock()
	}
}

func (s *SparkLine) trimOutOfRangeValues(maxSize int) {
	if maxSize < len(s.values) {
		s.values = append(s.values[:0], s.values[len(s.values)-maxSize:]...)
	}
}

func (s *SparkLine) Draw(buffer *ui.Buffer) {

	s.mutex.Lock()

	textStyle := ui.NewStyle(s.palette.BaseColor)

	height := s.Inner.Dy()
	minValue := util.FormatValue(s.minValue, s.scale)
	maxValue := util.FormatValue(s.maxValue, s.scale)
	curValue := util.FormatValue(0, s.scale)

	if len(s.values) > 0 {
		curValue = util.FormatValue(s.values[len(s.values)-1], s.scale)
	}

	overrideColor, hasOverride := s.ColorOverride(s.label)

	indent := 2 + util.Max([]int{
		len(minValue), len(maxValue), len(curValue),
	})

	for i := len(s.values) - 1; i >= 0; i-- {

		n := len(s.values) - i

		if n > s.Dx()-indent-3 {
			break
		}

		top := 0

		if s.maxValue != s.minValue {
			top = int((s.values[i] - s.minValue) * float64(height) / (s.maxValue - s.minValue))
		}

		for j := 0; j <= top; j++ {
			color := console.GetGradientColor(s.gradient, j, height)
			if hasOverride {
				color = overrideColor
			}
			buffer.SetCell(ui.NewCell(console.SymbolVerticalBar, ui.NewStyle(color)), image.Pt(s.Inner.Max.X-n-indent, s.Inner.Max.Y-j-1))
		}

		if i == len(s.values)-1 {
			buffer.SetString(curValue, textStyle, image.Pt(s.Inner.Max.X-n-indent+2, s.Inner.Max.Y-top-1))
			if s.maxValue != s.minValue {
				buffer.SetString(minValue, textStyle, image.Pt(s.Inner.Max.X-n-indent+2, s.Inner.Max.Y-1))
				buffer.SetString(maxValue, textStyle, image.Pt(s.Inner.Max.X-n-indent+2, s.Inner.Min.Y))
			}
		}
	}

	s.mutex.Unlock()

	s.Block.Draw(buffer)
	component.RenderAlert(s.Alert, s.Rectangle, buffer)
}
