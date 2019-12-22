package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/asset"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/asciibox"
	"github.com/sqshq/sampler/component/barchart"
	"github.com/sqshq/sampler/component/gauge"
	"github.com/sqshq/sampler/component/layout"
	"github.com/sqshq/sampler/component/runchart"
	"github.com/sqshq/sampler/component/sparkline"
	"github.com/sqshq/sampler/component/textbox"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"github.com/sqshq/sampler/event"
	"time"
)

type Starter struct {
	player  *asset.AudioPlayer
	lout    *layout.Layout
	palette console.Palette
	opt     config.Options
	cfg     config.Config
}

func (s *Starter) startAll() []*data.Sampler {
	samplers := make([]*data.Sampler, 0)
	for _, c := range s.cfg.RunCharts {
		cpt := runchart.NewRunChart(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers))
	}
	for _, c := range s.cfg.SparkLines {
		cpt := sparkline.NewSparkLine(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers))
	}
	for _, c := range s.cfg.BarCharts {
		cpt := barchart.NewBarChart(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers))
	}
	for _, c := range s.cfg.Gauges {
		cpt := gauge.NewGauge(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Cur, c.Min, c.Max}, c.Triggers))
	}
	for _, c := range s.cfg.AsciiBoxes {
		cpt := asciibox.NewAsciiBox(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers))
	}
	for _, c := range s.cfg.TextBoxes {
		cpt := textbox.NewTextBox(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers))
	}
	return samplers
}

func (s *Starter) start(drawable ui.Drawable, consumer *data.Consumer, componentConfig config.ComponentConfig, itemsConfig []config.Item, triggersConfig []config.TriggerConfig) *data.Sampler {
	cpt := component.NewComponent(drawable, consumer, componentConfig)
	triggers := data.NewTriggers(triggersConfig, consumer, s.opt, s.player)
	items := data.NewItems(itemsConfig, *componentConfig.RateMs)
	s.lout.AddComponent(cpt)
	time.Sleep(10 * time.Millisecond) // desync coroutines
	return data.NewSampler(consumer, items, triggers, s.opt, s.cfg.Variables, *componentConfig.RateMs)
}

func main() {

	cfg, opt := config.LoadConfig()

	console.Init()
	defer console.Close()

	player := asset.NewAudioPlayer()
	if player != nil {
		defer player.Close()
	}

	palette := console.GetPalette(*cfg.Theme)
	lout := layout.NewLayout(component.NewStatusBar(*opt.ConfigFile, palette), component.NewMenu(palette))

	starter := &Starter{player, lout, palette, opt, *cfg}
	samplers := starter.startAll()

	handler := event.NewHandler(samplers, opt, lout)
	handler.HandleEvents()
}
