package widgets

import (
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/widgets/runchart"
	ui "github.com/sqshq/termui"
)

type Layout struct {
	ui.Block
	components []Component
	menu       *Menu
	mode       Mode
	selection  int
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
	columnsCount = 30
	rowsCount    = 30
)

func NewLayout(width, height int, menu *Menu) *Layout {

	block := *ui.NewBlock()
	block.SetRect(0, 0, width, height)

	return &Layout{
		Block:      block,
		components: make([]Component, 0),
		menu:       menu,
		mode:       ModeDefault,
		selection:  0,
	}
}

func (l *Layout) AddComponent(drawable ui.Drawable, title string, position Position, size Size, Type ComponentType) {
	l.components = append(l.components, Component{drawable, title, position, size, Type})
}

func (l *Layout) GetComponents(Type ComponentType) []ui.Drawable {

	var components []ui.Drawable

	for _, component := range l.components {
		if component.Type == Type {
			components = append(components, component.Drawable)
		}
	}

	return components
}

func (l *Layout) HandleConsoleEvent(e string) {
	switch e {
	case console.KeyEnter:
		switch l.mode {
		case ModeComponentSelect:
			l.menu.choose()
			l.mode = ModeMenuOptionSelect
		case ModeMenuOptionSelect:
			option := l.menu.getSelectedOption()
			switch option {
			case MenuOptionMove:
				l.mode = ModeComponentMove
				l.menu.moveOrResize()
			case MenuOptionResize:
				l.mode = ModeComponentResize
				l.menu.moveOrResize()
			case MenuOptionPinpoint:
				l.mode = ModeChartPinpoint
				l.menu.idle()
				chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
				chart.MoveSelection(0)
			case MenuOptionResume:
				l.mode = ModeDefault
				l.menu.idle()
			}
		case ModeComponentMove:
			fallthrough
		case ModeComponentResize:
			l.menu.idle()
			l.mode = ModeDefault
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
			l.mode = ModeDefault
		}
	case console.KeyLeft:
		switch l.mode {
		case ModeDefault:
			l.mode = ModeComponentSelect
			l.selection = 0
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
			l.mode = ModeComponentSelect
			l.selection = 0
			l.menu.highlight(l.getComponent(l.selection))
		case ModeChartPinpoint:
			chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
			chart.MoveSelection(1)
		case ModeComponentSelect:
			if l.selection < len(l.components)-1 {
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
		case ModeMenuOptionSelect:
			l.menu.up()
		case ModeComponentMove:
			l.getSelectedComponent().Move(0, -1)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(0, -1)
		}
	case console.KeyDown:
		switch l.mode {
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
	return l.components[i]
}

func (l *Layout) getSelectedComponent() *Component {
	return &l.components[l.selection]
}

func (l *Layout) Draw(buffer *ui.Buffer) {

	columnWidth := float64(l.GetRect().Dx()) / columnsCount
	rowHeight := float64(l.GetRect().Dy()) / rowsCount

	for _, component := range l.components {

		x1 := float64(component.Position.X) * columnWidth
		y1 := float64(component.Position.Y) * rowHeight
		x2 := x1 + float64(component.Size.X)*columnWidth
		y2 := y1 + float64(component.Size.Y)*rowHeight

		component.Drawable.SetRect(int(x1), int(y1), int(x2), int(y2))
		component.Drawable.Draw(buffer)
	}

	l.menu.Draw(buffer)
}
