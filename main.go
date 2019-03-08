package main

import (
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
	"github.com/sqshq/sampler/trigger"
	ui "github.com/sqshq/termui"
)

func main() {

	cfg, flg := config.Load()
	csl := console.Console{}
	csl.Init()
	defer csl.Close()

	width, height := ui.TerminalDimensions()
	lout := layout.NewLayout(width, height, component.NewStatusLine(flg.ConfigFileName), component.NewMenu())

	for _, c := range cfg.RunCharts {

		legend := runchart.Legend{Enabled: c.Legend.Enabled, Details: c.Legend.Details}
		chart := runchart.NewRunChart(c, legend)
		lout.AddComponent(config.TypeRunChart, chart, c.Title, c.Position, c.Size, *c.RefreshRateMs)
		triggers := trigger.NewTriggers(c.Triggers)

		for _, i := range c.Items {
			item := data.Item{Label: *i.Label, Script: i.Script, Color: *i.Color}
			chart.AddLine(item.Label, item.Color)
			data.NewSampler(chart.Consumer, item, triggers, *c.RefreshRateMs)
		}
	}

	for _, a := range cfg.AsciiBoxes {
		box := asciibox.NewAsciiBox(a)
		item := data.Item{Label: *a.Label, Script: a.Script, Color: *a.Color}
		triggers := trigger.NewTriggers(a.Triggers)
		lout.AddComponent(config.TypeAsciiBox, box, a.Title, a.Position, a.Size, *a.RefreshRateMs)
		data.NewSampler(box.Consumer, item, triggers, *a.RefreshRateMs)
	}

	for _, b := range cfg.BarCharts {

		chart := barchart.NewBarChart(b.Title, *b.Scale)
		triggers := trigger.NewTriggers(b.Triggers)
		lout.AddComponent(config.TypeBarChart, chart, b.Title, b.Position, b.Size, *b.RefreshRateMs)

		for _, i := range b.Items {
			item := data.Item{Label: *i.Label, Script: i.Script, Color: *i.Color}
			chart.AddBar(*i.Label, *i.Color)
			data.NewSampler(chart.Consumer, item, triggers, *b.RefreshRateMs)
		}
	}

	for _, gc := range cfg.Gauges {

		g := gauge.NewGauge(gc.Title, *gc.Scale, *gc.Color)
		triggers := trigger.NewTriggers(gc.Triggers)
		lout.AddComponent(config.TypeGauge, g, gc.Title, gc.Position, gc.Size, *gc.RefreshRateMs)

		for _, i := range gc.Items {
			item := data.Item{Label: *i.Label, Script: i.Script}
			data.NewSampler(g.Consumer, item, triggers, *gc.RefreshRateMs)
		}
	}

	handler := event.NewHandler(lout)
	handler.HandleEvents()
}
