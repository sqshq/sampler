package main

import (
	ui "github.com/sqshq/termui"
	"github.com/sqshq/vcmd/config"
	"github.com/sqshq/vcmd/data"
	"github.com/sqshq/vcmd/settings"
	"github.com/sqshq/vcmd/widgets"
	"log"
	"time"
)

func main() {

	print("\033]0;vcmd\007")

	cfg := config.Load("/Users/sqshq/Go/src/github.com/sqshq/vcmd/config.yml")

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()
	events := ui.PollEvents()

	pollers := make([]data.Poller, 0)
	lout := widgets.NewLayout(ui.TerminalDimensions())

	for _, chartConfig := range cfg.RunCharts {

		chart := widgets.NewRunChart(chartConfig.Title)
		lout.AddItem(chart, chartConfig.Position, chartConfig.Size)

		for _, item := range chartConfig.Items {
			pollers = append(pollers,
				data.NewPoller(chart, item, chartConfig.RefreshRateMs))
		}
	}

	ticker := time.NewTicker(30 * time.Millisecond)
	pause := false

	for {
		select {
		case e := <-events:
			switch e.ID {
			case settings.EventQuit, settings.EventExit:
				return
			case settings.EventResize:
				payload := e.Payload.(ui.Resize)
				lout.ChangeDimensions(payload.Width, payload.Height)
			case settings.EventMouseClick:
				//payload := e.Payload.(ui.Mouse)
				//x, y := payload.X, payload.Y
				//log.Printf("x: %v, y: %v", x, y)
			}
			switch e.Type {
			case ui.KeyboardEvent:
				switch e.ID {
				case settings.EventKeyboardLeft:
					// here we are going to move selection (special type of layout item)
					//lout.GetItem("").Move(-1, 0)
				case settings.EventKeyboardRight:
					//lout.GetItem(0).Move(1, 0)
				case settings.EventKeyboardDown:
					//lout.GetItem(0).Move(0, 1)
				case settings.EventKeyboardUp:
					//lout.GetItem(0).Move(0, -1)
				case settings.EventPause:
					pause = !pause
				}
			}
		case <-ticker.C:
			if !pause {
				ui.Render(lout)
			}
		}
	}
}
