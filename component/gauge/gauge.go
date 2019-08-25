package gauge

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/util"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
)

const (
	MinValueLabel = "min"
	MaxValueLabel = "max"
	CurValueLabel = "cur"
)

// Gauge displays cur value between specified min and max values
type Gauge struct {
	*ui.Block
	*data.Consumer
	minValue    float64
	maxValue    float64
	curValue    float64
	color       ui.Color
	scale       int
	percentOnly bool
	palette     console.Palette
}

func NewGauge(c config.GaugeConfig, palette console.Palette) *Gauge {

	g := Gauge{
		Block:       component.NewBlock(c.Title, true, palette),
		Consumer:    data.NewConsumer(),
		color:       *c.Color,
		scale:       *c.Scale,
		percentOnly: *c.PercentOnly,
		palette:     palette,
	}

	go func() {
		for {
			select {
			case sample := <-g.SampleChannel:
				g.ConsumeSample(sample)
			case alert := <-g.AlertChannel:
				g.Alert = alert
			}
		}
	}()

	return &g
}

func (g *Gauge) ConsumeSample(sample *data.Sample) {

	float, err := util.ParseFloat(sample.Value)
	if err != nil {
		g.HandleConsumeFailure("Failed to parse a number", err, sample)
		return
	}

	g.HandleConsumeSuccess()

	switch sample.Label {
	case MinValueLabel:
		g.minValue = float
		break
	case MaxValueLabel:
		g.maxValue = float
		break
	case CurValueLabel:
		g.curValue = float
		break
	}
}

func (g *Gauge) Draw(buffer *ui.Buffer) {

	g.Block.Draw(buffer)

	percent := calculatePercent(g)

	var label string
	if g.percentOnly {
		label = fmt.Sprintf(" %v%% ", util.FormatValue(percent, g.scale))
	} else {
		label = fmt.Sprintf(" %v%% (%v) ", util.FormatValue(percent, g.scale), util.FormatValue(g.curValue, g.scale))
	}

	// draw bar
	barWidth := int((percent / 100) * float64(g.Inner.Dx()))
	if barWidth == 0 {
		barWidth = 1
	} else if barWidth > g.Dx()-2 {
		barWidth = g.Dx() - 2
	}
	buffer.Fill(
		ui.NewCell(console.SymbolVerticalBar, ui.NewStyle(g.color)),
		image.Rect(g.Inner.Min.X+1, g.Inner.Min.Y, g.Inner.Min.X+barWidth, g.Inner.Max.Y),
	)

	// draw label
	labelXCoordinate := g.Inner.Min.X + (g.Inner.Dx() / 2) - int(float64(len(label))/2)
	labelYCoordinate := g.Inner.Min.Y + ((g.Inner.Dy() - 1) / 2)
	if labelYCoordinate < g.Inner.Max.Y {
		for i, char := range label {
			style := ui.NewStyle(g.palette.BaseColor)
			if labelXCoordinate+i+1 <= g.Inner.Min.X+barWidth {
				style = ui.NewStyle(g.palette.BaseColor, ui.ColorClear)
			}
			buffer.SetCell(ui.NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
		}
	}

	component.RenderAlert(g.Alert, g.Rectangle, buffer)
}

func calculatePercent(g *Gauge) float64 {
	if g.curValue != g.minValue && g.maxValue != g.minValue {
		return (100 * (g.curValue - g.minValue)) / (g.maxValue - g.minValue)
	}
	return 0
}
