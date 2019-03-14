package component

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"strings"
)

type Alerter struct {
	channel <-chan data.Alert
	alert   *data.Alert
}

func NewAlerter(channel <-chan data.Alert) *Alerter {
	alerter := Alerter{channel: channel}
	alerter.consume()
	return &alerter
}

func (a *Alerter) consume() {
	go func() {
		for {
			select {
			case alert := <-a.channel:
				a.alert = &alert
			}
		}
	}()
}

func (a *Alerter) RenderAlert(buffer *ui.Buffer, area image.Rectangle) {

	if a.alert == nil {
		return
	}

	color := console.ColorWhite

	if a.alert.Color != nil {
		color = *a.alert.Color
	}

	width := max(len(a.alert.Title), len(a.alert.Text)) + 10

	if width > area.Dx() {
		width = area.Dx()
	}

	cells := ui.WrapCells(ui.ParseStyles(fmt.Sprintf("%s\n%s\n",
		strings.ToUpper(a.alert.Title), a.alert.Text), ui.NewStyle(console.ColorWhite)), uint(width))

	var lines []string
	line := ""

	for i := 0; i < len(cells); i++ {
		if cells[i].Rune == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(cells[i].Rune)
		}
	}

	block := *ui.NewBlock()
	block.SetRect(getRectCoordinates(area, width, len(lines)))
	block.BorderStyle = ui.Style{Fg: color}
	block.Draw(buffer)

	buffer.Fill(ui.NewCell(' ', ui.NewStyle(console.ColorBlack)), block.Inner)

	for i := 0; i < len(lines); i++ {
		buffer.SetString(lines[i],
			ui.NewStyle(color), getMiddlePoint(block.Inner, lines[i], i-1))
	}
}

//TODO move to utils
func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func getRectCoordinates(area image.Rectangle, width int, height int) (int, int, int, int) {
	x1 := area.Min.X + area.Dx()/2 - width/2
	y1 := area.Min.Y + area.Dy()/2 - height
	return x1, y1, x1 + width, y1 + height + 2
}
