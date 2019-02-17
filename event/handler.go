package event

import (
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/widgets"
	ui "github.com/sqshq/termui"
	"time"
)

type Handler struct {
	Layout        *widgets.Layout
	RenderEvents  <-chan time.Time
	ConsoleEvents <-chan ui.Event
}

func (h *Handler) HandleEvents() {

	pause := false

	for {
		select {
		case <-h.RenderEvents:
			if !pause {
				ui.Render(h.Layout)
			}
		case e := <-h.ConsoleEvents:
			switch e.ID {
			case console.KeyQuit, console.KeyExit:
				h.handleExit()
				return
			case console.KeyPause:
				pause = !pause
			case console.SignalResize:
				payload := e.Payload.(ui.Resize)
				h.Layout.ChangeDimensions(payload.Width, payload.Height)
			default:
				h.Layout.HandleConsoleEvent(e.ID)
			}
		}
	}
}

func (h *Handler) handleExit() {
	var settings []config.ComponentSettings
	for _, c := range h.Layout.Components {
		settings = append(settings,
			config.ComponentSettings{Type: c.Type, Title: c.Title, Size: c.Size, Position: c.Position})
	}
	config.Update(settings)
}
