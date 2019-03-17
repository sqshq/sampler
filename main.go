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
	flags  config.Flags
}

func (s *Starter) start(drawable ui.Drawable, consumer *data.Consumer, conponentConfig config.ComponentConfig, itemsConfig []config.Item, triggersConfig []config.TriggerConfig) {
	cpt := component.NewComponent(drawable, consumer, conponentConfig)
	triggers := data.NewTriggers(triggersConfig, consumer, s.player)
	data.NewSampler(consumer, data.NewItems(itemsConfig), triggers, *conponentConfig.RefreshRateMs)
	s.lout.AddComponent(cpt)
}

func main() {

	cfg, flg := config.Load()

	console.Init()
	defer console.Close()

	player := asset.NewAudioPlayer()
	defer player.Close()

	width, height := ui.TerminalDimensions()
	lout := layout.NewLayout(width, height, component.NewStatusLine(flg.ConfigFileName), component.NewMenu())

	starter := &Starter{lout, player, flg}

	for _, c := range cfg.RunCharts {
		cpt := runchart.NewRunChart(c)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}

	for _, c := range cfg.AsciiBoxes {
		cpt := asciibox.NewAsciiBox(c)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers)
	}

	for _, c := range cfg.BarCharts {
		cpt := barchart.NewBarChart(c)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}

	for _, c := range cfg.Gauges {
		cpt := gauge.NewGauge(c)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}

	handler := event.NewHandler(lout)
	handler.HandleEvents()
}
