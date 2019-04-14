package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"image"
)

type Menu struct {
	*ui.Block
	options   []MenuOption
	component Component
	mode      MenuMode
	option    MenuOption
	palette   console.Palette
}

type MenuMode rune

const (
	MenuModeIdle          MenuMode = 0
	MenuModeHighlight     MenuMode = 1
	MenuModeOptionSelect  MenuMode = 2
	MenuModeMoveAndResize MenuMode = 3
)

type MenuOption string

const (
	MenuOptionMove     MenuOption = "MOVE"
	MenuOptionResize   MenuOption = "RESIZE"
	MenuOptionPinpoint MenuOption = "PINPOINT"
	MenuOptionResume   MenuOption = "RESUME"
)

const (
	minimalMenuHeight = 8
)

func NewMenu(palette console.Palette) *Menu {
	return &Menu{
		Block:   NewBlock("", true, palette),
		options: []MenuOption{MenuOptionMove, MenuOptionResize, MenuOptionPinpoint, MenuOptionResume},
		mode:    MenuModeIdle,
		option:  MenuOptionMove,
		palette: palette,
	}
}

func (m *Menu) GetSelectedOption() MenuOption {
	return m.option
}

func (m *Menu) Highlight(component *Component) {
	m.component = *component
	m.updateDimensions()
	m.mode = MenuModeHighlight
	m.Title = component.Title
}

func (m *Menu) Choose() {
	m.mode = MenuModeOptionSelect
}

func (m *Menu) Idle() {
	m.mode = MenuModeIdle
}

func (m *Menu) Up() {
	for i := 1; i < len(m.options); i++ {
		if m.options[i] == m.option {
			m.option = m.options[i-1]
			break
		}
	}
	if m.option == MenuOptionPinpoint && m.component.Type != config.TypeRunChart {
		m.Up()
	}
}

func (m *Menu) Down() {
	for i := 0; i < len(m.options)-1; i++ {
		if m.options[i] == m.option {
			m.option = m.options[i+1]
			break
		}
	}
	if m.option == MenuOptionPinpoint && m.component.Type != config.TypeRunChart {
		m.Down()
	}
}

func (m *Menu) MoveOrResize() {
	m.mode = MenuModeMoveAndResize
}

func (m *Menu) Draw(buffer *ui.Buffer) {

	if m.mode == MenuModeIdle {
		return
	}

	m.updateDimensions()
	buffer.Fill(ui.NewCell(' ', ui.NewStyle(m.palette.ReverseColor)), m.GetRect())

	if m.Dy() > minimalMenuHeight {
		m.drawInnerBorder(buffer)
	}

	m.Block.Draw(buffer)

	switch m.mode {
	case MenuModeHighlight:
		m.renderHighlight(buffer)
	case MenuModeMoveAndResize:
		m.renderMoveAndResize(buffer)
	case MenuModeOptionSelect:
		m.renderOptions(buffer)
	}
}

func (m *Menu) renderHighlight(buffer *ui.Buffer) {

	arrowsText := "Use mouse or arrows for selection"
	optionsText := "<ENTER> to view options"
	resumeText := "<ESC> to resume"

	if m.Dy() <= minimalMenuHeight {
		buffer.SetString(
			optionsText,
			ui.NewStyle(console.ColorDarkGrey),
			getMiddlePoint(m.Block.Rectangle, optionsText, -1),
		)
		return
	}

	m.printAllDirectionsArrowSign(buffer, -2)

	arrowsTextPoint := getMiddlePoint(m.Block.Rectangle, arrowsText, 2)
	if arrowsTextPoint.Y+1 < m.Inner.Max.Y {
		buffer.SetString(
			arrowsText,
			ui.NewStyle(console.ColorDarkGrey),
			arrowsTextPoint,
		)
	}

	optionsTextPoint := getMiddlePoint(m.Block.Rectangle, optionsText, 3)
	if optionsTextPoint.Y+1 < m.Inner.Max.Y {
		buffer.SetString(
			optionsText,
			ui.NewStyle(console.ColorDarkGrey),
			getMiddlePoint(m.Block.Rectangle, optionsText, 3),
		)
	}

	resumeTextPoint := getMiddlePoint(m.Block.Rectangle, resumeText, 4)
	if resumeTextPoint.Y+1 < m.Inner.Max.Y {
		buffer.SetString(
			resumeText,
			ui.NewStyle(console.ColorDarkGrey),
			resumeTextPoint,
		)
	}
}

func (m *Menu) renderMoveAndResize(buffer *ui.Buffer) {

	saveText := "<ENTER> to save changes"

	if m.Dy() <= minimalMenuHeight {
		buffer.SetString(saveText, ui.NewStyle(console.ColorDarkGrey), getMiddlePoint(m.Block.Rectangle, saveText, -1))
		return
	}

	m.printAllDirectionsArrowSign(buffer, -1)
	buffer.SetString(saveText, ui.NewStyle(console.ColorDarkGrey), getMiddlePoint(m.Block.Rectangle, saveText, 3))
}

func (m *Menu) printAllDirectionsArrowSign(buffer *ui.Buffer, y int) {

	arrows := []string{
		"  ↑  ",
		"←   →",
		"  ↓  ",
	}

	for i, a := range arrows {
		printString(
			a,
			ui.NewStyle(console.ColorOlive),
			getMiddlePoint(m.Block.Rectangle, a, i+y),
			buffer,
		)
	}
}

func (m *Menu) renderOptions(buffer *ui.Buffer) {

	highlightedStyle := ui.NewStyle(m.palette.ReverseColor, console.ColorOlive)
	regularStyle := ui.NewStyle(m.palette.BaseColor, m.palette.ReverseColor)

	offset := 1
	for _, option := range m.options {

		style := regularStyle
		if m.option == option {
			style = highlightedStyle
		}

		if option != MenuOptionPinpoint || m.component.Type == config.TypeRunChart {
			offset += 2
			point := getMiddlePoint(m.Block.Rectangle, string(option), offset-6)
			buffer.SetString(string(option), style, point)
		}
	}
}

func (m *Menu) updateDimensions() {
	r := m.component.GetRect()
	m.SetRect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

func (m *Menu) drawInnerBorder(buffer *ui.Buffer) {

	verticalCell := ui.Cell{ui.VERTICAL_LINE, m.BorderStyle}
	horizontalCell := ui.Cell{ui.HORIZONTAL_LINE, m.BorderStyle}

	// draw lines
	buffer.Fill(horizontalCell, image.Rect(m.Min.X+2, m.Min.Y+2, m.Max.X-2, m.Min.Y))
	buffer.Fill(horizontalCell, image.Rect(m.Min.X+2, m.Max.Y-2, m.Max.X-2, m.Max.Y))
	buffer.Fill(verticalCell, image.Rect(m.Min.X+2, m.Min.Y+1, m.Min.X+3, m.Max.Y-1))
	buffer.Fill(verticalCell, image.Rect(m.Max.X-2, m.Min.Y, m.Max.X-3, m.Max.Y))

	// draw corners
	buffer.SetCell(ui.Cell{ui.TOP_LEFT, m.BorderStyle}, image.Pt(m.Min.X+2, m.Min.Y+1))
	buffer.SetCell(ui.Cell{ui.TOP_RIGHT, m.BorderStyle}, image.Pt(m.Max.X-3, m.Min.Y+1))
	buffer.SetCell(ui.Cell{ui.BOTTOM_LEFT, m.BorderStyle}, image.Pt(m.Min.X+2, m.Max.Y-2))
	buffer.SetCell(ui.Cell{ui.BOTTOM_RIGHT, m.BorderStyle}, image.Pt(m.Max.X-3, m.Max.Y-2))
}

// TODO move to utils
func getMiddlePoint(rectangle image.Rectangle, text string, offset int) image.Point {
	return image.Pt(rectangle.Min.X+rectangle.Dx()/2-len(text)/2, rectangle.Max.Y-rectangle.Dy()/2+offset)
}

// TODO move to utils
func printString(s string, style ui.Style, p image.Point, buffer *ui.Buffer) {
	for i, char := range s {
		buffer.SetCell(ui.Cell{Rune: char, Style: style}, image.Pt(p.X+i, p.Y))
	}
}
