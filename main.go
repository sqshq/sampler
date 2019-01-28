package main

import (
	ui "github.com/sqshq/termui"
	"github.com/sqshq/vcmd/config"
	"github.com/sqshq/vcmd/widgets"
	"log"
	"time"
)

func main() {

	cfg := config.Load("/Users/sqshq/Go/src/github.com/sqshq/vcmd/config.yml")

	for _, linechart := range cfg.LineCharts {
		for _, data := range linechart.Data {
			value, _ := data.NextValue()
			log.Printf("%s: %s - %v", linechart.Title, data.Label, value)
		}
	}

	p1 := widgets.NewTimePlot()
	p1.Title = " CURL LATENCY STATISTICS (sec) "
	p1.SetRect(0, 20, 148, 40)
	p1.LineColors[0] = ui.ColorYellow
	p1.Marker = widgets.MarkerBraille

	if err := ui.Init(); err != nil {
		//log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()
	uiEvents := ui.PollEvents()

	dataTicker := time.NewTicker(200 * time.Millisecond)
	uiTicker := time.NewTicker(50 * time.Millisecond)

	pause := false

	go func() {
		for {
			select {
			case <-dataTicker.C:
				if !pause {
					value, err := cfg.LineCharts[0].Data[0].NextValue()
					if err != nil {
						log.Printf("failed to get value: %s", err)
						break
					}
					p1.AddValue(value)
				}
			}
		}
	}()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			}
			//case "<MouseLeft>":
			//	payload := e.Payload.(ui.Mouse)
			//	x, y := payload.X, payload.Y
			//	log.Printf("x: %v, y: %v", x, y)
			//}
			switch e.Type {
			case ui.KeyboardEvent: // handle all key presses
				//log.Printf("key: %v", e.ID)
				switch e.ID {
				// TODO refactor + control moving out of range
				case "<Left>":
					rect := p1.GetRect()
					min := rect.Min
					max := rect.Max
					p1.SetRect(min.X-1, min.Y, max.X-1, max.Y)
					ui.Clear()
				case "<Right>":
					rect := p1.GetRect()
					min := rect.Min
					max := rect.Max
					p1.SetRect(min.X+1, min.Y, max.X+1, max.Y)
					ui.Clear()
				case "<Down>":
					rect := p1.GetRect()
					min := rect.Min
					max := rect.Max
					p1.SetRect(min.X, min.Y+1, max.X, max.Y+1)
					ui.Clear()
				case "<Up>":
					rect := p1.GetRect()
					min := rect.Min
					max := rect.Max
					p1.SetRect(min.X, min.Y-1, max.X, max.Y-1)
					ui.Clear()
				case "p":
					pause = !pause
				}
			}
		case <-uiTicker.C:
			if !pause {
				ui.Render(p1)
			}
		}
	}
}
