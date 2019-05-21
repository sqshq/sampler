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
	"github.com/sqshq/sampler/storage"
	"time"
)

type Starter struct {
	lout   *layout.Layout
	player *asset.AudioPlayer
	opt    config.Options
	cfg    config.Config
}

func (s *Starter) start(drawable ui.Drawable, consumer *data.Consumer, componentConfig config.ComponentConfig, itemsConfig []config.Item, triggersConfig []config.TriggerConfig) {
	cpt := component.NewComponent(drawable, consumer, componentConfig)
	triggers := data.NewTriggers(triggersConfig, consumer, s.opt, s.player)
	items := data.NewItems(itemsConfig, *componentConfig.RateMs)
	data.NewSampler(consumer, items, triggers, s.opt, s.cfg.Variables, *componentConfig.RateMs)
	s.lout.AddComponent(cpt)
	time.Sleep(10 * time.Millisecond) // desync coroutines
}

func main() {

	cfg, opt := config.LoadConfig()

	console.Init()
	defer console.Close()

	player := asset.NewAudioPlayer()
	defer player.Close()

	palette := console.GetPalette(*cfg.Theme)
	width, height := ui.TerminalDimensions()

	license := storage.GetLicense()
	_ = storage.UpdateStatistics(cfg, width, height)

	lout := layout.NewLayout(width, height, component.NewStatusLine(opt.ConfigFile, palette, license), component.NewMenu(palette), component.NewIntro(palette))
	starter := &Starter{lout, player, opt, cfg}

	if license == nil {
		lout.RunIntro()
		storage.InitLicense()
	} else if !license.Purchased /* && random */ {
		// TODO lout.showNagWindow() with timeout and OK button
		// TODO verify license
		// TODO send stats
	}

	for _, c := range cfg.RunCharts {
		cpt := runchart.NewRunChart(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}

	for _, c := range cfg.SparkLines {
		cpt := sparkline.NewSparkLine(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers)
	}

	for _, c := range cfg.BarCharts {
		cpt := barchart.NewBarChart(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}

	for _, c := range cfg.Gauges {
		cpt := gauge.NewGauge(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Cur, c.Min, c.Max}, c.Triggers)
	}

	for _, c := range cfg.AsciiBoxes {
		cpt := asciibox.NewAsciiBox(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers)
	}

	for _, c := range cfg.TextBoxes {
		cpt := textbox.NewTextBox(c, palette)
		starter.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers)
	}

	handler := event.NewHandler(lout, opt)
	handler.HandleEvents()
}
