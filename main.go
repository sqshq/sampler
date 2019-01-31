package main

import (
	ui "github.com/sqshq/termui"
	"github.com/sqshq/vcmd/config"
	"github.com/sqshq/vcmd/data"
	"github.com/sqshq/vcmd/layout"
	"github.com/sqshq/vcmd/widgets"
	"log"
	"time"
)

/*
 TODO validation
 - title uniquness and mandatory within a single type of widget
 - label uniqueness and mandatory (if > 1 data bullets)
*/
func main() {

	// todo error handling + validation
	cfg := config.Load("/Users/sqshq/Go/src/github.com/sqshq/vcmd/config.yml")

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()
	events := ui.PollEvents()

	pollers := make([]data.Poller, 0)
	lout := layout.NewLayout(ui.TerminalDimensions())

	for _, chartConfig := range cfg.RunCharts {

		chart := widgets.NewRunChart(chartConfig.Title)
		lout.AddItem(chart, chartConfig.Position, chartConfig.Size)

		for _, chartData := range chartConfig.DataConfig {
			pollers = append(pollers,
				data.NewPoller(chart, chartData.Script, chartData.Label, chartConfig.RefreshRateMs))
		}
	}

	ticker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case e := <-events:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				lout.ChangeDimensions(payload.Width, payload.Height)
			case "<MouseLeft>":
				//payload := e.Payload.(ui.Mouse)
				//x, y := payload.X, payload.Y
				//log.Printf("x: %v, y: %v", x, y)
			}
			switch e.Type {
			case ui.KeyboardEvent:
				switch e.ID {
				case "<Left>":
					// here we are going to move selection (special type of layout item)
					//lout.GetItem("").Move(-1, 0)
				case "<Right>":
					//lout.GetItem(0).Move(1, 0)
				case "<Down>":
					//lout.GetItem(0).Move(0, 1)
				case "<Up>":
					//lout.GetItem(0).Move(0, -1)
				case "p":
					for _, poller := range pollers {
						poller.TogglePause()
					}
				}
			}
		case <-ticker.C:
			ui.Render(lout)
		}
	}
}
