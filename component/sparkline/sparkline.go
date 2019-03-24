package sparkline

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"strconv"
)

type SparkLine struct {
	*ui.Block
	*data.Consumer
	alert   *data.Alert
	values  []float64
	palette console.Palette
}

func NewSparkLine(c config.SparkLineConfig, palette console.Palette) *SparkLine {

	line := &SparkLine{
		Block:    component.NewBlock(c.Title, true, palette),
		Consumer: data.NewConsumer(),
		values:   []float64{},
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
}

func (s *SparkLine) Draw(buffer *ui.Buffer) {
	s.Block.Draw(buffer)
}
