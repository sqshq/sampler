package main

import (
	ui "github.com/sqshq/termui"
	"github.com/sqshq/termui/widgets"
	"github.com/sqshq/vcmd/app/config"
	"log"
	"math"
)

func main() {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	cfg := config.Load("/Users/sqshq/config.yml")

	p1 := widgets.NewPlot()
	p1.Title = "dot-mode line Chart"
	p1.Marker = widgets.MarkerDot
	p1.Data = [][]float64{
		{0, 1, 2, 3},
		{4, 5, 6, 7},
	}
	p1.SetRect(50, 0, 100, 10)
	p1.DataLabels = []string{"hello"}
	p1.DotRune = '+'
	p1.AxesColor = ui.ColorWhite
	p1.LineColors[0] = ui.ColorCyan
	p1.DrawDirection = widgets.DrawLeft

	for _, linechart := range cfg.LineCharts {
		for _, data := range linechart.Data {
			value, _ := data.NextValue()
			log.Printf("%s: %s - %v", linechart.Title, data.Label, value)
		}
	}

	ui.Render(p1)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}

func printChart() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	sinData := func() [][]float64 {
		n := 220
		data := make([][]float64, 2)
		data[0] = make([]float64, n)
		data[1] = make([]float64, n)
		for i := 0; i < n; i++ {
			data[0][i] = 1 + math.Sin(float64(i)/5)
			data[1][i] = 1 + math.Cos(float64(i)/5)
		}
		return data
	}()

	p0 := widgets.NewPlot()
	p0.Title = "braille-mode Line Chart"
	p0.Data = sinData
	p0.SetRect(0, 0, 50, 15)
	p0.AxesColor = ui.ColorWhite
	p0.LineColors[0] = ui.ColorGreen

	p1 := widgets.NewPlot()
	p1.Title = "dot-mode line Chart"
	p1.Marker = widgets.MarkerDot
	p1.Data = [][]float64{{1, 2, 3, 4, 5}}
	p1.SetRect(50, 0, 75, 10)
	p1.DotRune = '+'
	p1.AxesColor = ui.ColorWhite
	p1.LineColors[0] = ui.ColorYellow
	p1.DrawDirection = widgets.DrawLeft

	p2 := widgets.NewPlot()
	p2.Title = "dot-mode Scatter Plot"
	p2.Marker = widgets.MarkerDot
	p2.Data = make([][]float64, 2)
	p2.Data[0] = []float64{1, 2, 3, 4, 5}
	p2.Data[1] = sinData[1][4:]
	p2.SetRect(0, 15, 50, 30)
	p2.AxesColor = ui.ColorWhite
	p2.LineColors[0] = ui.ColorCyan
	p2.Type = widgets.ScatterPlot

	p3 := widgets.NewPlot()
	p3.Title = "braille-mode Scatter Plot"
	p3.Data = make([][]float64, 2)
	p3.Data[0] = []float64{1, 2, 3, 4, 5}
	p3.Data[1] = sinData[1][4:]
	p3.SetRect(45, 15, 80, 30)
	p3.AxesColor = ui.ColorWhite
	p3.LineColors[0] = ui.ColorCyan
	p3.Marker = widgets.MarkerBraille
	p3.Type = widgets.ScatterPlot

	ui.Render(p0, p1, p2, p3)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
