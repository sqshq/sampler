package component

import (
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	ui "github.com/sqshq/termui"
	"image"
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

	buffer.Fill(ui.NewCell(' ', ui.NewStyle(console.ColorBlack)), area)
	buffer.SetString(a.alert.Title, ui.NewStyle(console.ColorWhite), getMiddlePoint(area, a.alert.Title, -1))
	buffer.SetString(a.alert.Text, ui.NewStyle(console.ColorWhite), getMiddlePoint(area, a.alert.Text, 0))
}
