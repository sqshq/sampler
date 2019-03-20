package event

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component/layout"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"time"
)

const (
	refreshRateToRenderRateRatio = 0.5
)

type Handler struct {
	layout        *layout.Layout
	renderTicker  *time.Ticker
	consoleEvents <-chan ui.Event
	renderRate    time.Duration
	options       config.Options
}

func NewHandler(layout *layout.Layout, options config.Options) *Handler {
	renderRate := calcMinRenderRate(layout)
	return &Handler{
		layout:        layout,
		consoleEvents: ui.PollEvents(),
		renderTicker:  time.NewTicker(renderRate),
		renderRate:    renderRate,
		options:       options,
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
			case console.KeyQuit, console.KeyExit:
				h.handleExit()
				return
			case console.SignalResize:
				payload := e.Payload.(ui.Resize)
				h.layout.ChangeDimensions(payload.Width, payload.Height)
			default:
				h.layout.HandleConsoleEvent(e.ID)
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
	case layout.ModePause:
		// proceed with stopped timer
	default:
		h.renderTicker = time.NewTicker(console.MinRenderInterval)
	}
}

func (h *Handler) handleExit() {
	var settings []config.ComponentSettings
	for _, c := range h.layout.Components {
		settings = append(settings,
			config.ComponentSettings{Type: c.Type, Title: c.Title, Size: c.Size, Location: c.Location})
	}
	config.Update(settings, h.options)
}

func calcMinRenderRate(layout *layout.Layout) time.Duration {

	minRefreshRateMs := layout.Components[0].RefreshRateMs
	for _, c := range layout.Components {
		if c.RefreshRateMs < minRefreshRateMs {
			minRefreshRateMs = c.RefreshRateMs
		}
	}

	renderRate := time.Duration(
		int(float64(minRefreshRateMs)*refreshRateToRenderRateRatio)) * time.Millisecond

	if renderRate < console.MinRenderInterval {
		return console.MinRenderInterval
	}

	if renderRate > console.MaxRenderInterval {
		return console.MaxRenderInterval
	}

	return renderRate
}
