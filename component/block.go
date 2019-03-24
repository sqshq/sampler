package component

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/console"
)

func NewBlock(title string, border bool, palette console.Palette) *ui.Block {
	style := ui.Style{Fg: palette.BaseColor, Bg: ui.ColorClear}
	block := ui.NewBlock()
	block.Border = border
	block.BorderStyle = style
	block.TitleStyle = style
	if len(title) > 0 {
		block.Title = fmt.Sprintf(" %s ", title)
	}
	return block
}
