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
	"sync"
	"unicode/utf8"
)

const (
	MinValueLabel = config.GaugeMinLabel
	MaxValueLabel = config.GaugeMaxLabel
	CurValueLabel = config.GaugeCurLabel
)

// subGauge est une jauge à l'intérieur du bloc
type subGauge struct {
	title       string // affiché à droite de la jauge (vide en mode simple)
	curLabel    string // libellé de l'échantillon de valeur courante
	minLabel    string // libellé du min propre à cette jauge ("" si absent)
	maxLabel    string // libellé du max propre à cette jauge ("" si absent)
	curValue    float64
	minValue    float64
	maxValue    float64
	color       ui.Color
	scale       int
	percentOnly bool
}

// Gauge displays one or several cur values between specified min and max values
type Gauge struct {
	*ui.Block
	*data.Consumer
	gauges   []*subGauge
	multi    bool
	minValue float64 // min global du bloc
	maxValue float64 // max global du bloc
	palette  console.Palette
	mutex    *sync.Mutex
}

func NewGauge(c config.GaugeConfig, palette console.Palette) *Gauge {

	g := Gauge{
		Block:    component.NewComponentBlock(c.Title, c.IsTitleVisible(), *c.Border, palette),
		Consumer: data.NewConsumer(),
		multi:    c.IsMulti(),
		palette:  palette,
		mutex:    &sync.Mutex{},
	}

	if !g.multi {
		g.gauges = []*subGauge{{
			curLabel:    CurValueLabel,
			color:       *c.Color,
			scale:       *c.Scale,
			percentOnly: *c.PercentOnly,
		}}
	}

	for _, it := range c.Items {
		s := &subGauge{
			title:       *it.Label,
			curLabel:    *it.Label,
			color:       *it.Color,
			scale:       *c.Scale,
			percentOnly: *c.PercentOnly,
		}
		if it.Scale != nil {
			s.scale = *it.Scale
		}
		if it.PercentOnly != nil {
			s.percentOnly = *it.PercentOnly
		}
		if it.Min != nil {
			s.minLabel = config.GaugeItemMinLabel(*it.Label)
		}
		if it.Max != nil {
			s.maxLabel = config.GaugeItemMaxLabel(*it.Label)
		}
		g.gauges = append(g.gauges, s)
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

	g.mutex.Lock()
	defer g.mutex.Unlock()

	switch sample.Label {
	case MinValueLabel:
		g.minValue = float
		return
	case MaxValueLabel:
		g.maxValue = float
		return
	}

	for _, s := range g.gauges {
		switch sample.Label {
		case s.curLabel:
			s.curValue = float
		case s.minLabel:
			s.minValue = float
		case s.maxLabel:
			s.maxValue = float
		}
	}
}

func (g *Gauge) Draw(buffer *ui.Buffer) {

	g.Block.Draw(buffer)

	g.mutex.Lock()

	if !g.multi {
		area := image.Rect(g.Inner.Min.X+1, g.Inner.Min.Y, g.Inner.Max.X, g.Inner.Max.Y)
		g.drawGauge(buffer, g.gauges[0], area)
	} else {
		titleWidth := 0
		for _, s := range g.gauges {
			if w := utf8.RuneCountInString(s.title); w > titleWidth {
				titleWidth = w
			}
		}

		count := len(g.gauges)
		y1 := g.Inner.Min.Y
		for i, s := range g.gauges {
			height := g.Inner.Dy() / count
			if i < g.Inner.Dy()%count {
				height++
			}
			if height <= 0 {
				break
			}
			y2 := y1 + height
			buffer.SetString(s.title, ui.NewStyle(g.palette.BaseColor),
				image.Pt(g.Inner.Min.X+1, y1+(y2-y1-1)/2))
			area := image.Rect(g.Inner.Min.X+titleWidth+2, y1, g.Inner.Max.X, y2)
			g.drawGauge(buffer, s, area)
			y1 = y2
		}
	}

	g.mutex.Unlock()

	component.RenderAlert(g.Alert, g.Rectangle, buffer)
}

func (g *Gauge) drawGauge(buffer *ui.Buffer, s *subGauge, area image.Rectangle) {

	if area.Dx() <= 0 || area.Dy() <= 0 {
		return
	}

	minValue, maxValue := g.minValue, g.maxValue
	if s.minLabel != "" {
		minValue = s.minValue
	}
	if s.maxLabel != "" {
		maxValue = s.maxValue
	}

	percent := calculatePercent(s.curValue, minValue, maxValue)

	var label string
	if s.percentOnly {
		label = fmt.Sprintf(" %v%% ", util.FormatValue(percent, s.scale))
	} else {
		label = fmt.Sprintf(" %v%% (%v) ", util.FormatValue(percent, s.scale), util.FormatValue(s.curValue, s.scale))
	}

	// draw bar
	barWidth := int((percent / 100) * float64(area.Dx()))
	if barWidth < 1 {
		barWidth = 1
	} else if barWidth > area.Dx() {
		barWidth = area.Dx()
	}
	buffer.Fill(
		ui.NewCell(console.SymbolVerticalBar, ui.NewStyle(g.ColorFor(s.curLabel, s.color))),
		image.Rect(area.Min.X, area.Min.Y, area.Min.X+barWidth, area.Max.Y),
	)

	// draw label, centré dans la zone de la jauge
	labelXCoordinate := area.Min.X + (area.Dx() / 2) - int(float64(len(label))/2)
	labelYCoordinate := area.Min.Y + ((area.Dy() - 1) / 2)
	for i, char := range []rune(label) {
		x := labelXCoordinate + i
		if x < area.Min.X || x >= area.Max.X {
			continue
		}
		buffer.SetCell(ui.NewCell(char, ui.NewStyle(g.palette.BaseColor)), image.Pt(x, labelYCoordinate))
	}
}

func calculatePercent(cur, min, max float64) float64 {
	if cur != min && max != min {
		return (100 * (cur - min)) / (max - min)
	}
	return 0
}
