package component

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component/util"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"strings"
)

func RenderAlert(alert *data.Alert, area image.Rectangle, buffer *ui.Buffer) {

	if alert == nil {
		return
	}

	color := console.ColorWhite

	if alert.Color != nil {
		color = *alert.Color
	}

	width := util.Max([]int{len(alert.Title), len(alert.Text)}) + 10

	if width > area.Dx() {
		width = area.Dx()
	}

	cells := ui.WrapCells(ui.ParseStyles(fmt.Sprintf("%s\n%s\n",
		strings.ToUpper(alert.Title), alert.Text), ui.NewStyle(console.ColorWhite)), uint(width))

	var lines []string
	line := ""

	for i := 0; i < len(cells); i++ {
		if cells[i].Rune == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(cells[i].Rune)
		}
	}

	block := *ui.NewBlock()
	block.SetRect(util.GetRectCoordinates(area, width, len(lines)))
	block.BorderStyle = ui.Style{Fg: color, Bg: ui.ColorClear}
	block.Draw(buffer)

	buffer.Fill(ui.NewCell(' ', ui.NewStyle(console.ColorBlack)), block.Inner)

	for i := 0; i < len(lines); i++ {
		buffer.SetString(lines[i],
			ui.NewStyle(color), util.GetMiddlePoint(block.Inner, lines[i], i-1))
	}
}
