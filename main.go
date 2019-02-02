package main

import (
	ui "github.com/sqshq/termui"
	"github.com/sqshq/vcmd/config"
	"github.com/sqshq/vcmd/console"
	"github.com/sqshq/vcmd/data"
	"github.com/sqshq/vcmd/event"
	"github.com/sqshq/vcmd/widgets"
	"time"
)

func main() {

	cfg := config.Load("/Users/sqshq/Go/src/github.com/sqshq/vcmd/config.yml")
	csl := console.Console{}
	csl.Init()
	defer csl.Close()

	layout := widgets.NewLayout(ui.TerminalDimensions())

	for _, chartConfig := range cfg.RunCharts {

		chart := widgets.NewRunChart(chartConfig.Title)
		layout.AddComponent(chart, chartConfig.Position, chartConfig.Size, widgets.TypeRunChart)

		for _, item := range chartConfig.Items {
			data.NewPoller(chart, item, chartConfig.RefreshRateMs)
		}
	}

	handler := event.Handler{
		Layout:        layout,
		RenderEvents:  time.NewTicker(console.RenderRate).C,
		ConsoleEvents: ui.PollEvents(),
	}

	handler.HandleEvents()
}
