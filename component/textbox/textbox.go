package textbox

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
)

// TextBox represents a component with regular text
type TextBox struct {
	*ui.Block
	*data.Consumer
	label  string
	alert  *data.Alert
	text   string
	border bool
	style  ui.Style
}

func NewTextBox(c config.TextBoxConfig, palette console.Palette) *TextBox {

	color := c.Color
	if color == nil {
		color = &palette.BaseColor
	}

	box := TextBox{
		Block:    component.NewComponentBlock(c.Title, c.IsTitleVisible(), *c.Border, palette),
		Consumer: data.NewConsumer(),
		label:    *c.Label,
		style:    ui.NewStyle(*color),
		border:   *c.Border,
	}

	go func() {
		for {
			select {
			case sample := <-box.SampleChannel:
				box.text = sample.Value
			case alert := <-box.AlertChannel:
				box.alert = alert
			}
		}
	}()

	return &box
}

func (t *TextBox) Draw(buffer *ui.Buffer) {

	t.Block.Draw(buffer)

	style := t.style
	if color, ok := t.ColorOverride(t.label); ok {
		style = ui.NewStyle(color)
	}

	cells := ui.ParseStyles(t.text, ui.Theme.Paragraph.Text)
	cells = ui.WrapCells(cells, uint(t.Inner.Dx()-2))

	rows := ui.SplitCells(cells, '\n')

	marginY := 0
	if t.border {
		marginY = 1
	}

	for y, row := range rows {
		if y+t.Inner.Min.Y+marginY >= t.Inner.Max.Y-marginY {
			break
		}
		row = ui.TrimCells(row, t.Inner.Dx()-2)
		for _, cx := range ui.BuildCellWithXArray(row) {
			x, cell := cx.X, cx.Cell
			cell.Style = style
			buffer.SetCell(cell, image.Pt(x+1, y+marginY).Add(t.Inner.Min))
		}
	}

	component.RenderAlert(t.alert, t.Rectangle, buffer)
}
