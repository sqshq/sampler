package event

import (
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

func (self *Handler) HandleEvents() {

	pause := false

	for {
		select {
		case <-self.RenderEvents:
			if !pause {
				ui.Render(self.Layout)
			}
		case e := <-self.ConsoleEvents:
			switch e.ID {
			case console.KeyQuit, console.KeyExit:
				return
			case console.KeyPause:
				pause = !pause
			case console.KeyResize:
				payload := e.Payload.(ui.Resize)
				self.Layout.ChangeDimensions(payload.Width, payload.Height)
			//case "a":
			//	self.Layout.GetComponent(0).DisableSelection()
			default:
				self.Layout.HandleConsoleEvent(e.ID)
			}
		}
	}
}
