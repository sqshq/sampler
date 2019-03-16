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

func main() {

	cfg, flg := config.Load()

	console.Init()
	defer console.Close()

	player := asset.NewAudioPlayer()
	defer player.Close()

	width, height := ui.TerminalDimensions()
	lout := layout.NewLayout(width, height, component.NewStatusLine(flg.ConfigFileName), component.NewMenu())

	for _, c := range cfg.RunCharts {
		chart := runchart.NewRunChart(c)
		triggers := data.NewTriggers(c.Triggers, chart.Consumer, player)
		data.NewSampler(chart.Consumer, data.NewItems(c.Items), triggers, *c.RefreshRateMs)
		lout.AddComponent(chart.Component, config.TypeRunChart)
	}

	for _, a := range cfg.AsciiBoxes {
		box := asciibox.NewAsciiBox(a)
		triggers := data.NewTriggers(a.Triggers, box.Consumer, player)
		data.NewSampler(box.Consumer, data.NewItems([]config.Item{a.Item}), triggers, *a.RefreshRateMs)
		lout.AddComponent(box.Component, config.TypeAsciiBox)
	}

	for _, b := range cfg.BarCharts {
		chart := barchart.NewBarChart(b)
		triggers := data.NewTriggers(b.Triggers, chart.Consumer, player)
		data.NewSampler(chart.Consumer, data.NewItems(b.Items), triggers, *b.RefreshRateMs)
		lout.AddComponent(chart.Component, config.TypeBarChart)
	}

	for _, gc := range cfg.Gauges {
		g := gauge.NewGauge(gc)
		triggers := data.NewTriggers(gc.Triggers, g.Consumer, player)
		data.NewSampler(g.Consumer, data.NewItems(gc.Items), triggers, *gc.RefreshRateMs)
		lout.AddComponent(g.Component, config.TypeGauge)
	}

	handler := event.NewHandler(lout)
	handler.HandleEvents()
}
