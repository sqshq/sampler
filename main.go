package main

import (
	"github.com/sqshq/sampler/asset"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/asciibox"
	"github.com/sqshq/sampler/component/barchart"
	"github.com/sqshq/sampler/component/runchart"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"github.com/sqshq/sampler/event"
	ui "github.com/sqshq/termui"
)

func main() {

	cfg := config.Load()
	csl := console.Console{}
	csl.Init()
	defer csl.Close()

	width, height := ui.TerminalDimensions()
	layout := component.NewLayout(width, height, component.NewMenu())

	for _, c := range cfg.RunCharts {

		legend := runchart.Legend{Enabled: c.Legend.Enabled, Details: c.Legend.Details}
		chart := runchart.NewRunChart(c.Title, *c.Scale, *c.RefreshRateMs, legend)
		layout.AddComponent(config.TypeRunChart, chart, c.Title, c.Position, c.Size, *c.RefreshRateMs)

		for _, item := range c.Items {
			chart.AddLine(*item.Label, *item.Color)
			data.NewSampler(chart, item, *c.RefreshRateMs)
		}
	}

	for _, a := range cfg.AsciiBoxes {
		box := asciibox.NewAsciiBox(a.Title, *a.Font, *a.Item.Color)
		layout.AddComponent(config.TypeAsciiBox, box, a.Title, a.Position, a.Size, *a.RefreshRateMs)
		data.NewSampler(box, a.Item, *a.RefreshRateMs)
	}

	for _, c := range cfg.BarCharts {

		chart := barchart.NewBarChart(c.Title, *c.Scale)
		layout.AddComponent(config.TypeBarChart, chart, c.Title, c.Position, c.Size, *c.RefreshRateMs)

		for _, item := range c.Items {
			chart.AddBar(*item.Label, *item.Color)
			data.NewSampler(chart, item, *c.RefreshRateMs)
		}
	}

	handler := event.NewHandler(layout)
	handler.HandleEvents()
}
