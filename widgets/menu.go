package widgets

import (
	"github.com/sqshq/sampler/console"
	ui "github.com/sqshq/termui"
	"image"
)

type Menu struct {
	ui.Block
	options   []MenuOption
	component Component
	mode      MenuMode
	option    MenuOption
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

func NewMenu() *Menu {
	block := *ui.NewBlock()
	block.Border = true
	block.BorderStyle = ui.NewStyle(console.ColorDarkGrey)
	return &Menu{
		Block:   block,
		options: []MenuOption{MenuOptionMove, MenuOptionResize, MenuOptionPinpoint, MenuOptionResume},
		mode:    MenuModeIdle,
		option:  MenuOptionMove,
	}
}

func (m *Menu) getSelectedOption() MenuOption {
	return m.option
}

func (m *Menu) highlight(component Component) {
	m.component = component
	m.updateDimensions()
	m.mode = MenuModeHighlight
	m.Title = component.Title
}

func (m *Menu) choose() {
	m.mode = MenuModeOptionSelect
}

func (m *Menu) idle() {
	m.mode = MenuModeIdle
}

func (m *Menu) up() {
	for i := 1; i < len(m.options); i++ {
		if m.options[i] == m.option {
			m.option = m.options[i-1]
			break
		}
	}
}

func (m *Menu) down() {
	for i := 0; i < len(m.options)-1; i++ {
		if m.options[i] == m.option {
			m.option = m.options[i+1]
			break
		}
	}
}

func (m *Menu) moveOrResize() {
	m.mode = MenuModeMoveAndResize
}

func (m *Menu) Draw(buffer *ui.Buffer) {

	if m.mode == MenuModeIdle {
		return
	}

	m.updateDimensions()

	buffer.Fill(
		ui.NewCell(' ', ui.NewStyle(ui.ColorClear, ui.ColorBlack)),
		m.GetRect(),
	)

	switch m.mode {
	case MenuModeHighlight:
		m.renderHighlight(buffer)
	case MenuModeMoveAndResize:
		m.renderMoveAndResize(buffer)
	case MenuModeOptionSelect:
		m.renderOptions(buffer)
	}

	m.drawInnerBorder(buffer)
	m.Block.Draw(buffer)
}

func (m *Menu) renderHighlight(buffer *ui.Buffer) {

	m.printAllDirectionsArrowSign(buffer, -2)

	arrowsText := "Use arrows for selection"
	arrowsTextPoint := getMiddlePoint(m.Block, arrowsText, 2)
	if arrowsTextPoint.In(m.Rectangle) {
		buffer.SetString(
			arrowsText,
			ui.NewStyle(console.ColorDarkGrey),
			arrowsTextPoint,
		)
	}

	optionsText := "<ENTER> to view options"
	optionsTextPoint := getMiddlePoint(m.Block, optionsText, 3)
	if optionsTextPoint.In(m.Rectangle) {
		buffer.SetString(
			optionsText,
			ui.NewStyle(console.ColorDarkGrey),
			getMiddlePoint(m.Block, optionsText, 3),
		)
	}

	resumeText := "<ESC> to resume"
	resumeTextPoint := getMiddlePoint(m.Block, resumeText, 4)
	if resumeTextPoint.In(m.Rectangle) {
		buffer.SetString(
			resumeText,
			ui.NewStyle(console.ColorDarkGrey),
			resumeTextPoint,
		)
	}
}

func (m *Menu) renderMoveAndResize(buffer *ui.Buffer) {

	m.printAllDirectionsArrowSign(buffer, -2)

	saveText := "<ENTER> to save changes"
	saveTextPoint := getMiddlePoint(m.Block, saveText, 4)
	if saveTextPoint.In(m.Rectangle) {
		buffer.SetString(
			saveText,
			ui.NewStyle(console.ColorDarkGrey),
			saveTextPoint,
		)
	}
}

func (m *Menu) printAllDirectionsArrowSign(buffer *ui.Buffer, y int) {

	arrows := []string{
		"  ↑  ",
		"←· →",
		"  ↓  ",
	}

	for i, a := range arrows {
		buffer.SetString(
			a,
			ui.NewStyle(console.ColorOlive),
			getMiddlePoint(m.Block, a, i+y),
		)
	}
}

func (m *Menu) renderOptions(buffer *ui.Buffer) {

	// TODO extract styles to console.Palette
	highlightedStyle := ui.NewStyle(console.ColorOlive, console.ColorBlack, ui.ModifierReverse)
	regularStyle := ui.NewStyle(console.ColorWhite)

	offset := 1
	for _, option := range m.options {

		style := regularStyle
		if m.option == option {
			style = highlightedStyle
		}

		if option != MenuOptionPinpoint || m.component.Type == TypeRunChart {
			offset += 2
			point := getMiddlePoint(m.Block, string(option), offset-5)
			if point.In(m.GetRect()) {
				buffer.SetString(string(option), style, point)
			}
		}
	}
}

func (m *Menu) updateDimensions() {
	r := m.component.Drawable.GetRect()
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

func getMiddlePoint(block ui.Block, text string, offset int) image.Point {
	return image.Pt(block.Min.X+block.Dx()/2-len(text)/2, block.Max.Y-block.Dy()/2+offset)
}
