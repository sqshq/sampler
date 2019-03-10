package main

import (
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
	ui "github.com/sqshq/termui"
)

func main() {

	cfg, flg := config.Load()
	csl := console.Console{}
	csl.Init()
	defer csl.Close()

	width, height := ui.TerminalDimensions()
	lout := layout.NewLayout(width, height, component.NewStatusLine(flg.ConfigFileName), component.NewMenu())

	player := asset.NewAudioPlayer()
	defer player.Close()

	for _, c := range cfg.RunCharts {

		legend := runchart.Legend{Enabled: c.Legend.Enabled, Details: c.Legend.Details}
		chart := runchart.NewRunChart(c, legend)
		lout.AddComponent(config.TypeRunChart, chart, c.Title, c.Position, c.Size, *c.RefreshRateMs)
		triggers := data.NewTriggers(c.Triggers, chart.Consumer, player)
		items := data.NewItems(c.Items)
		data.NewSampler(chart.Consumer, items, triggers, *c.RefreshRateMs)

		for _, i := range c.Items {
			chart.AddLine(*i.Label, *i.Color)
		}
	}

	for _, a := range cfg.AsciiBoxes {
		box := asciibox.NewAsciiBox(a)
		item := data.Item{Label: *a.Label, Script: a.Script, Color: a.Color}
		lout.AddComponent(config.TypeAsciiBox, box, a.Title, a.Position, a.Size, *a.RefreshRateMs)
		triggers := data.NewTriggers(a.Triggers, box.Consumer, player)
		data.NewSampler(box.Consumer, []data.Item{item}, triggers, *a.RefreshRateMs)
	}

	for _, b := range cfg.BarCharts {

		chart := barchart.NewBarChart(b.Title, *b.Scale)
		triggers := data.NewTriggers(b.Triggers, chart.Consumer, player)
		lout.AddComponent(config.TypeBarChart, chart, b.Title, b.Position, b.Size, *b.RefreshRateMs)
		items := data.NewItems(b.Items)
		data.NewSampler(chart.Consumer, items, triggers, *b.RefreshRateMs)

		for _, i := range b.Items {
			chart.AddBar(*i.Label, *i.Color)
		}
	}

	for _, gc := range cfg.Gauges {
		g := gauge.NewGauge(gc.Title, *gc.Scale, *gc.Color)
		triggers := data.NewTriggers(gc.Triggers, g.Consumer, player)
		lout.AddComponent(config.TypeGauge, g, gc.Title, gc.Position, gc.Size, *gc.RefreshRateMs)
		items := data.NewItems(gc.Items)
		data.NewSampler(g.Consumer, items, triggers, *gc.RefreshRateMs)
	}

	handler := event.NewHandler(lout)
	handler.HandleEvents()
}
