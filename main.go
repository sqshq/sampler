package main

import (
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"github.com/sqshq/sampler/event"
	"github.com/sqshq/sampler/widgets"
	ui "github.com/sqshq/termui"
	"time"
)

func main() {

	cfg := config.Load("/Users/sqshq/Go/src/github.com/sqshq/sampler/config.yml")
	csl := console.Console{}
	csl.Init()
	defer csl.Close()

	width, height := ui.TerminalDimensions()
	layout := widgets.NewLayout(width, height, widgets.NewMenu())

	for _, chartConfig := range cfg.RunCharts {

		chart := widgets.NewRunChart(chartConfig.Title, chartConfig.Precision, chartConfig.RefreshRateMs)
		layout.AddComponent(chart, chartConfig.Title, chartConfig.Position, chartConfig.Size, widgets.TypeRunChart)

		for _, item := range chartConfig.Items {
			data.NewSampler(chart, item, chartConfig.RefreshRateMs)
		}
	}

	handler := event.Handler{
		Layout:        layout,
		RenderEvents:  time.NewTicker(console.RenderRate).C,
		ConsoleEvents: ui.PollEvents(),
	}

	handler.HandleEvents()
}
