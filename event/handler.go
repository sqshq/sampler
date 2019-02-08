package event

import (
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
			case EventQuit, EventExit:
				return
			case EventPause:
				pause = !pause
			case EventResize:
				payload := e.Payload.(ui.Resize)
				self.Layout.ChangeDimensions(payload.Width, payload.Height)
			case "a":
				self.Layout.GetComponent(0).DisableSelection()
			case EventKeyboardLeft:
				self.Layout.GetComponent(0).MoveSelection(-1)
			case EventKeyboardRight:
				self.Layout.GetComponent(0).MoveSelection(+1)
			case EventKeyboardDown:
				//layout.GetItem(0).Move(0, 1)
			case EventKeyboardUp:
				//layout.GetItem(0).Move(0, -1)
			}
		}
	}
}
