package layout

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/runchart"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"math"
)

type Layout struct {
	ui.Block
	Components       []*component.Component
	statusbar        *component.StatusBar
	menu             *component.Menu
	ChangeModeEvents chan Mode
	mode             Mode
	selection        int
}

type Mode rune

const (
	ModeDefault          Mode = 0
	ModePause            Mode = 1
	ModeComponentSelect  Mode = 2
	ModeMenuOptionSelect Mode = 3
	ModeComponentMove    Mode = 4
	ModeComponentResize  Mode = 5
	ModeChartPinpoint    Mode = 6
)

const (
	columnsCount    = 80
	rowsCount       = 40
	minDimension    = 3
	statusbarHeight = 1
)

func NewLayout(width, height int, statusline *component.StatusBar, menu *component.Menu) *Layout {

	block := *ui.NewBlock()
	block.SetRect(0, 0, width, height)
	statusline.SetRect(0, height-statusbarHeight, width, height)

	return &Layout{
		Block:            block,
		Components:       make([]*component.Component, 0),
		statusbar:        statusline,
		menu:             menu,
		mode:             ModeDefault,
		selection:        0,
		ChangeModeEvents: make(chan Mode, 10),
	}
}

func (l *Layout) AddComponent(cpt *component.Component, Type config.ComponentType) {
	l.Components = append(l.Components, cpt)
}

func (l *Layout) changeMode(m Mode) {
	l.mode = m
	l.ChangeModeEvents <- m
}

func (l *Layout) HandleConsoleEvent(e string) {

	selected := l.getSelection()

	switch e {
	case console.KeyPause:
		if l.mode == ModePause {
			l.changeMode(ModeDefault)
		} else {
			if selected.Type == config.TypeRunChart {
				selected.CommandChannel <- data.Command{Type: runchart.CommandDisableSelection}
			}
			l.menu.Idle()
			l.changeMode(ModePause)
		}
	case console.KeyEnter:
		switch l.mode {
		case ModeComponentSelect:
			l.menu.Choose()
			l.changeMode(ModeMenuOptionSelect)
		case ModeMenuOptionSelect:
			option := l.menu.GetSelectedOption()
			switch option {
			case component.MenuOptionMove:
				l.changeMode(ModeComponentMove)
				l.menu.MoveOrResize()
			case component.MenuOptionResize:
				l.changeMode(ModeComponentResize)
				l.menu.MoveOrResize()
			case component.MenuOptionPinpoint:
				l.changeMode(ModeChartPinpoint)
				l.menu.Idle()
				selected.CommandChannel <- data.Command{Type: runchart.CommandMoveSelection, Value: 0}
			case component.MenuOptionResume:
				l.changeMode(ModeDefault)
				l.menu.Idle()
			}
		case ModeComponentMove:
			fallthrough
		case ModeComponentResize:
			l.menu.Idle()
			l.changeMode(ModeDefault)
		}
	case console.KeyEsc:
		switch l.mode {
		case ModeChartPinpoint:
			selected.CommandChannel <- data.Command{Type: runchart.CommandDisableSelection}
			fallthrough
		case ModeComponentSelect:
			fallthrough
		case ModeMenuOptionSelect:
			l.menu.Idle()
			l.changeMode(ModeDefault)
		}
	case console.KeyLeft:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeChartPinpoint:
			selected.CommandChannel <- data.Command{Type: runchart.CommandMoveSelection, Value: -1}
		case ModeComponentSelect:
			l.moveSelection(e)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeComponentMove:
			selected.Move(-1, 0)
		case ModeComponentResize:
			selected.Resize(-1, 0)
		}
	case console.KeyRight:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeChartPinpoint:
			selected.CommandChannel <- data.Command{Type: runchart.CommandMoveSelection, Value: 1}
		case ModeComponentSelect:
			l.moveSelection(e)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeComponentMove:
			selected.Move(1, 0)
		case ModeComponentResize:
			selected.Resize(1, 0)
		}
	case console.KeyUp:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeComponentSelect:
			l.moveSelection(e)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeMenuOptionSelect:
			l.menu.Up()
		case ModeComponentMove:
			selected.Move(0, -1)
		case ModeComponentResize:
			selected.Resize(0, -1)
		}
	case console.KeyDown:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeComponentSelect:
			l.moveSelection(e)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeMenuOptionSelect:
			l.menu.Down()
		case ModeComponentMove:
			selected.Move(0, 1)
		case ModeComponentResize:
			selected.Resize(0, 1)
		}
	}
}

func (l *Layout) ChangeDimensions(width, height int) {
	l.SetRect(0, 0, width, height)
}

func (l *Layout) getComponent(i int) *component.Component {
	return l.Components[i]
}

func (l *Layout) getSelection() *component.Component {
	return l.Components[l.selection]
}

func (l *Layout) moveSelection(direction string) {

	previouslySelected := l.getSelection()
	newlySelectedIndex := l.selection

	for i, current := range l.Components {

		if current == previouslySelected {
			continue
		}

		if newlySelectedIndex < 0 {
			newlySelectedIndex = i
		}

		var previouslySelectedCornerPoint image.Point
		var newlySelectedCornerPoint image.Point
		var currentCornerPoint image.Point

		switch direction {
		case console.KeyLeft:
			previouslySelectedCornerPoint = component.GetRectLeftSideCenter(previouslySelected.GetRect())
			newlySelectedCornerPoint = component.GetRectRightSideCenter(l.getComponent(newlySelectedIndex).GetRect())
			currentCornerPoint = component.GetRectRightSideCenter(current.GetRect())
		case console.KeyRight:
			previouslySelectedCornerPoint = component.GetRectRightSideCenter(previouslySelected.GetRect())
			newlySelectedCornerPoint = component.GetRectLeftSideCenter(l.getComponent(newlySelectedIndex).GetRect())
			currentCornerPoint = component.GetRectLeftSideCenter(current.GetRect())
		case console.KeyUp:
			previouslySelectedCornerPoint = component.GetRectTopSideCenter(previouslySelected.GetRect())
			newlySelectedCornerPoint = component.GetRectBottomSideCenter(l.getComponent(newlySelectedIndex).GetRect())
			currentCornerPoint = component.GetRectBottomSideCenter(current.GetRect())
		case console.KeyDown:
			previouslySelectedCornerPoint = component.GetRectBottomSideCenter(previouslySelected.GetRect())
			newlySelectedCornerPoint = component.GetRectTopSideCenter(l.getComponent(newlySelectedIndex).GetRect())
			currentCornerPoint = component.GetRectTopSideCenter(current.GetRect())
		}

		if component.GetDistance(previouslySelectedCornerPoint, currentCornerPoint) <
			component.GetDistance(previouslySelectedCornerPoint, newlySelectedCornerPoint) {
			newlySelectedIndex = i
		}
	}

	l.selection = newlySelectedIndex
}

func (l *Layout) Draw(buffer *ui.Buffer) {

	columnWidth := float64(l.GetRect().Dx()) / float64(columnsCount)
	rowHeight := float64(l.GetRect().Dy()-statusbarHeight) / float64(rowsCount)

	for _, c := range l.Components {

		x1 := math.Floor(float64(c.Position.X) * columnWidth)
		y1 := math.Floor(float64(c.Position.Y) * rowHeight)
		x2 := x1 + math.Floor(float64(c.Size.X))*columnWidth
		y2 := y1 + math.Floor(float64(c.Size.Y))*rowHeight

		if x2-x1 < minDimension {
			x2 = x1 + minDimension
		}

		if y2-y1 < minDimension {
			y2 = y1 + minDimension
		}

		c.SetRect(int(x1), int(y1), int(x2), int(y2))
		c.Draw(buffer)
	}

	l.statusbar.SetRect(
		0, l.GetRect().Dy()-statusbarHeight,
		l.GetRect().Dx(), l.GetRect().Dy())

	l.statusbar.Draw(buffer)
	l.menu.Draw(buffer)
}
