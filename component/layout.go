package component

import (
	"github.com/sqshq/sampler/component/runchart"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	ui "github.com/sqshq/termui"
)

type Layout struct {
	ui.Block
	Components       []Component
	ChangeModeEvents chan Mode
	statusbar        *StatusBar
	menu             *Menu
	mode             Mode
	selection        int
}

type Mode rune

const (
	ModeDefault          Mode = 0
	ModeComponentSelect  Mode = 1
	ModeMenuOptionSelect Mode = 2
	ModeComponentMove    Mode = 3
	ModeComponentResize  Mode = 4
	ModeChartPinpoint    Mode = 5
)

const (
	columnsCount    = 80
	rowsCount       = 60
	statusbarHeight = 1
)

func NewLayout(width, height int, statusline *StatusBar, menu *Menu) *Layout {

	block := *ui.NewBlock()
	block.SetRect(0, 0, width, height)
	statusline.SetRect(0, height-statusbarHeight, width, height)

	return &Layout{
		Block:            block,
		Components:       make([]Component, 0),
		statusbar:        statusline,
		menu:             menu,
		mode:             ModeDefault,
		selection:        0,
		ChangeModeEvents: make(chan Mode, 10),
	}
}

func (l *Layout) AddComponent(Type config.ComponentType, drawable ui.Drawable, title string, position config.Position, size config.Size, refreshRateMs int) {
	l.Components = append(l.Components, Component{Type, drawable, title, position, size, refreshRateMs})
}

func (l *Layout) GetComponents(Type config.ComponentType) []ui.Drawable {

	var components []ui.Drawable

	for _, component := range l.Components {
		if component.Type == Type {
			components = append(components, component.Drawable)
		}
	}

	return components
}

func (l *Layout) changeMode(m Mode) {
	l.mode = m
	l.ChangeModeEvents <- m
}

func (l *Layout) HandleConsoleEvent(e string) {
	switch e {
	case console.KeyEnter:
		switch l.mode {
		case ModeComponentSelect:
			l.menu.choose()
			l.changeMode(ModeMenuOptionSelect)
		case ModeMenuOptionSelect:
			option := l.menu.getSelectedOption()
			switch option {
			case MenuOptionMove:
				l.changeMode(ModeComponentMove)
				l.menu.moveOrResize()
			case MenuOptionResize:
				l.changeMode(ModeComponentResize)
				l.menu.moveOrResize()
			case MenuOptionPinpoint:
				l.changeMode(ModeChartPinpoint)
				l.menu.idle()
				chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
				chart.MoveSelection(0)
			case MenuOptionResume:
				l.changeMode(ModeDefault)
				l.menu.idle()
			}
		case ModeComponentMove:
			fallthrough
		case ModeComponentResize:
			l.menu.idle()
			l.changeMode(ModeDefault)
		}
	case console.KeyEsc:
		switch l.mode {
		case ModeChartPinpoint:
			chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
			chart.DisableSelection()
			fallthrough
		case ModeComponentSelect:
			fallthrough
		case ModeMenuOptionSelect:
			l.menu.idle()
			l.changeMode(ModeDefault)
		}
	case console.KeyLeft:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.highlight(l.getComponent(l.selection))
		case ModeChartPinpoint:
			chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
			chart.MoveSelection(-1)
		case ModeComponentSelect:
			if l.selection > 0 {
				l.selection--
			}
			l.menu.highlight(l.getComponent(l.selection))
		case ModeComponentMove:
			l.getSelectedComponent().Move(-1, 0)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(-1, 0)
		}
	case console.KeyRight:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.highlight(l.getComponent(l.selection))
		case ModeChartPinpoint:
			chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
			chart.MoveSelection(1)
		case ModeComponentSelect:
			if l.selection < len(l.Components)-1 {
				l.selection++
			}
			l.menu.highlight(l.getComponent(l.selection))
		case ModeComponentMove:
			l.getSelectedComponent().Move(1, 0)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(1, 0)
		}
	case console.KeyUp:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.highlight(l.getComponent(l.selection))
		case ModeComponentSelect:
			if l.selection > 0 {
				l.selection--
			}
			l.menu.highlight(l.getComponent(l.selection))
		case ModeMenuOptionSelect:
			l.menu.up()
		case ModeComponentMove:
			l.getSelectedComponent().Move(0, -1)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(0, -1)
		}
	case console.KeyDown:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.highlight(l.getComponent(l.selection))
		case ModeComponentSelect:
			if l.selection < len(l.Components)-1 {
				l.selection++
			}
			l.menu.highlight(l.getComponent(l.selection))
		case ModeMenuOptionSelect:
			l.menu.down()
		case ModeComponentMove:
			l.getSelectedComponent().Move(0, 1)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(0, 1)
		}
	}
}

func (l *Layout) ChangeDimensions(width, height int) {
	l.SetRect(0, 0, width, height)
}

// TODO func to get prev/next component navigating left/right/top/bottom
func (l *Layout) getComponent(i int) Component {
	return l.Components[i]
}

func (l *Layout) getSelectedComponent() *Component {
	return &l.Components[l.selection]
}

func (l *Layout) Draw(buffer *ui.Buffer) {

	columnWidth := float64(l.GetRect().Dx()) / columnsCount
	rowHeight := float64(l.GetRect().Dy()-statusbarHeight) / rowsCount

	for _, component := range l.Components {

		x1 := float64(component.Position.X) * columnWidth
		y1 := float64(component.Position.Y) * rowHeight
		x2 := x1 + float64(component.Size.X)*columnWidth
		y2 := y1 + float64(component.Size.Y)*rowHeight

		component.Drawable.SetRect(int(x1), int(y1), int(x2), int(y2))
		component.Drawable.Draw(buffer)
	}

	l.statusbar.SetRect(
		0, l.GetRect().Dy()-statusbarHeight,
		l.GetRect().Dx(), l.GetRect().Dy())

	l.menu.Draw(buffer)
	l.statusbar.Draw(buffer)
}
