package event

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component/layout"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"time"
)

const (
	refreshRateToRenderRateRatio = 0.5
)

type Handler struct {
	samplers      []*data.Sampler
	options       config.Options
	layout        *layout.Layout
	renderTicker  *time.Ticker
	consoleEvents <-chan ui.Event
	renderRate    time.Duration
}

func NewHandler(samplers []*data.Sampler, options config.Options, layout *layout.Layout) *Handler {
	renderRate := calcMinRenderRate(layout)
	return &Handler{
		samplers:      samplers,
		options:       options,
		layout:        layout,
		consoleEvents: ui.PollEvents(),
		renderTicker:  time.NewTicker(renderRate),
		renderRate:    renderRate,
	}
}

func (h *Handler) HandleEvents() {

	// initial render
	ui.Render(h.layout)

	for {
		select {
		case mode := <-h.layout.ChangeModeEvents:
			h.handleModeChange(mode)
		case <-h.renderTicker.C:
			ui.Render(h.layout)
		case e := <-h.consoleEvents:
			switch e.ID {
			case console.SignalClick:
				payload := e.Payload.(ui.Mouse)
				h.layout.HandleMouseClick(payload.X, payload.Y)
			case console.KeyQuit1, console.KeyQuit2, console.KeyQuit3:
				if h.layout.WerePositionsChanged() {
					h.updateConfigFile()
				}
				return
			case console.SignalResize:
				payload := e.Payload.(ui.Resize)
				h.layout.ChangeDimensions(payload.Width, payload.Height)
			default:
				h.layout.HandleKeyboardEvent(e.ID)
			}
		}
	}
}

func (h *Handler) handleModeChange(m layout.Mode) {

	// render the change before switching the tickers
	ui.Render(h.layout)
	h.renderTicker.Stop()

	switch m {
	case layout.ModeDefault:
		h.renderTicker = time.NewTicker(h.renderRate)
		h.pause(false)
	case layout.ModePause:
		h.pause(true)
		// proceed with stopped timer
	default:
		h.renderTicker = time.NewTicker(console.MinRenderInterval)
		h.pause(false)
	}
}

func (h *Handler) pause(pause bool) {
	for _, s := range h.samplers {
		s.Pause(pause)
	}
}

func (h *Handler) updateConfigFile() {
	var settings []config.ComponentSettings
	for _, c := range h.layout.Components {
		settings = append(settings,
			config.ComponentSettings{Type: c.Type, Title: c.Title, Size: c.Size, Location: c.Location})
	}
	config.Update(settings, h.options)
}

func calcMinRenderRate(layout *layout.Layout) time.Duration {

	minRateMs := layout.Components[0].RateMs
	for _, c := range layout.Components {
		if c.RateMs < minRateMs {
			minRateMs = c.RateMs
		}
	}

	renderRate := time.Duration(
		int(float64(minRateMs)*refreshRateToRenderRateRatio)) * time.Millisecond

	if renderRate < console.MinRenderInterval {
		return console.MinRenderInterval
	}

	if renderRate > console.MaxRenderInterval {
		return console.MaxRenderInterval
	}

	return renderRate
}
