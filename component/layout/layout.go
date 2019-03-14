package layout

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/runchart"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"image"
	"math"
)

type Layout struct {
	ui.Block
	Components       []component.Component
	ChangeModeEvents chan Mode
	statusbar        *component.StatusBar
	menu             *component.Menu
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
		Components:       make([]component.Component, 0),
		statusbar:        statusline,
		menu:             menu,
		mode:             ModeDefault,
		selection:        0,
		ChangeModeEvents: make(chan Mode, 10),
	}
}

func (l *Layout) AddComponent(Type config.ComponentType, drawable ui.Drawable, title string, position config.Position, size config.Size, refreshRateMs int) {
	l.Components = append(l.Components, component.Component{
		Type:          Type,
		Drawable:      drawable,
		Title:         title,
		Position:      position,
		Size:          size,
		RefreshRateMs: refreshRateMs,
	})
}

func (l *Layout) GetComponents(Type config.ComponentType) []component.Component {

	var components []component.Component

	for _, c := range l.Components {
		if c.Type == Type {
			components = append(components, c)
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
	case console.KeyPause:
		if l.mode == ModePause {
			l.changeMode(ModeDefault)
		} else {
			if l.getSelectedComponent().Type == config.TypeRunChart {
				chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
				chart.DisableSelection()
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
				chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
				chart.MoveSelection(0)
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
			chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
			chart.DisableSelection()
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
			chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
			chart.MoveSelection(-1)
		case ModeComponentSelect:
			l.moveSelection(e)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeComponentMove:
			l.getSelectedComponent().Move(-1, 0)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(-1, 0)
		}
	case console.KeyRight:
		switch l.mode {
		case ModeDefault:
			l.changeMode(ModeComponentSelect)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeChartPinpoint:
			chart := l.getSelectedComponent().Drawable.(*runchart.RunChart)
			chart.MoveSelection(1)
		case ModeComponentSelect:
			l.moveSelection(e)
			l.menu.Highlight(l.getComponent(l.selection))
		case ModeComponentMove:
			l.getSelectedComponent().Move(1, 0)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(1, 0)
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
			l.getSelectedComponent().Move(0, -1)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(0, -1)
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
			l.getSelectedComponent().Move(0, 1)
		case ModeComponentResize:
			l.getSelectedComponent().Resize(0, 1)
		}
	}
}

func (l *Layout) ChangeDimensions(width, height int) {
	l.SetRect(0, 0, width, height)
}

func (l *Layout) getComponent(i int) *component.Component {
	return &l.Components[i]
}

func (l *Layout) getSelectedComponent() *component.Component {
	return &l.Components[l.selection]
}

func (l *Layout) moveSelection(direction string) {

	previouslySelected := *l.getSelectedComponent()
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
			previouslySelectedCornerPoint = component.GetRectLeftAgeCenter(previouslySelected.Drawable.GetRect())
			newlySelectedCornerPoint = component.GetRectRightAgeCenter(l.getComponent(newlySelectedIndex).Drawable.GetRect())
			currentCornerPoint = component.GetRectRightAgeCenter(current.Drawable.GetRect())
		case console.KeyRight:
			previouslySelectedCornerPoint = component.GetRectRightAgeCenter(previouslySelected.Drawable.GetRect())
			newlySelectedCornerPoint = component.GetRectLeftAgeCenter(l.getComponent(newlySelectedIndex).Drawable.GetRect())
			currentCornerPoint = component.GetRectLeftAgeCenter(current.Drawable.GetRect())
		case console.KeyUp:
			previouslySelectedCornerPoint = component.GetRectTopAgeCenter(previouslySelected.Drawable.GetRect())
			newlySelectedCornerPoint = component.GetRectBottomAgeCenter(l.getComponent(newlySelectedIndex).Drawable.GetRect())
			currentCornerPoint = component.GetRectBottomAgeCenter(current.Drawable.GetRect())
		case console.KeyDown:
			previouslySelectedCornerPoint = component.GetRectBottomAgeCenter(previouslySelected.Drawable.GetRect())
			newlySelectedCornerPoint = component.GetRectTopAgeCenter(l.getComponent(newlySelectedIndex).Drawable.GetRect())
			currentCornerPoint = component.GetRectTopAgeCenter(current.Drawable.GetRect())
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

		c.Drawable.SetRect(int(x1), int(y1), int(x2), int(y2))
		c.Drawable.Draw(buffer)
	}

	l.statusbar.SetRect(
		0, l.GetRect().Dy()-statusbarHeight,
		l.GetRect().Dx(), l.GetRect().Dy())

	l.statusbar.Draw(buffer)
	l.menu.Draw(buffer)
}
