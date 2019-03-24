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
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"github.com/sqshq/sampler/event"
)

type Starter struct {
	lout   *layout.Layout
	player *asset.AudioPlayer
	opt    config.Options
}

func (s *Starter) start(drawable ui.Drawable, consumer *data.Consumer, conponentConfig config.ComponentConfig, itemsConfig []config.Item, triggersConfig []config.TriggerConfig) {
	cpt := component.NewComponent(drawable, consumer, conponentConfig)
	triggers := data.NewTriggers(triggersConfig, consumer, s.opt, s.player)
	data.NewSampler(consumer, data.NewItems(itemsConfig), triggers, s.opt, *conponentConfig.RefreshRateMs)
	s.lout.AddComponent(cpt)
}

func main() {

	cfg, opt := config.Load()

	console.Init()
	defer console.Close()

	player := asset.NewAudioPlayer()
	defer player.Close()

	palette := console.GetPalette(*cfg.Theme)
	width, height := ui.TerminalDimensions()

	lout := layout.NewLayout(width, height, component.NewStatusLine(opt.ConfigFile, palette), component.NewMenu(palette))

	starter := &Starter{lout, player, opt}

	for _, c := range cfg.RunCharts {
		cpt := runchart.NewRunChart(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}

	for _, c := range cfg.AsciiBoxes {
		cpt := asciibox.NewAsciiBox(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers)
	}

	for _, c := range cfg.BarCharts {
		cpt := barchart.NewBarChart(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}

	for _, c := range cfg.Gauges {
		cpt := gauge.NewGauge(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}

	handler := event.NewHandler(lout, opt)
	handler.HandleEvents()
}
