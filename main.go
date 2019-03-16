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
		cpt := component.NewComponent(chart, chart.Consumer, c.ComponentConfig, config.TypeRunChart)
		triggers := data.NewTriggers(c.Triggers, chart.Consumer, player)
		data.NewSampler(chart.Consumer, data.NewItems(c.Items), triggers, *c.RefreshRateMs)
		lout.AddComponent(cpt)
	}

	for _, a := range cfg.AsciiBoxes {
		box := asciibox.NewAsciiBox(a)
		cpt := component.NewComponent(box, box.Consumer, a.ComponentConfig, config.TypeRunChart)
		triggers := data.NewTriggers(a.Triggers, box.Consumer, player)
		data.NewSampler(box.Consumer, data.NewItems([]config.Item{a.Item}), triggers, *a.RefreshRateMs)
		lout.AddComponent(cpt)
	}

	for _, b := range cfg.BarCharts {
		chart := barchart.NewBarChart(b)
		cpt := component.NewComponent(chart, chart.Consumer, b.ComponentConfig, config.TypeRunChart)
		triggers := data.NewTriggers(b.Triggers, chart.Consumer, player)
		data.NewSampler(chart.Consumer, data.NewItems(b.Items), triggers, *b.RefreshRateMs)
		lout.AddComponent(cpt)
	}

	for _, gc := range cfg.Gauges {
		g := gauge.NewGauge(gc)
		cpt := component.NewComponent(g, g.Consumer, gc.ComponentConfig, config.TypeRunChart)
		triggers := data.NewTriggers(gc.Triggers, g.Consumer, player)
		data.NewSampler(g.Consumer, data.NewItems(gc.Items), triggers, *gc.RefreshRateMs)
		lout.AddComponent(cpt)
	}

	handler := event.NewHandler(lout)
	handler.HandleEvents()
}
