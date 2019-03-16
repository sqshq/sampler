package component

import (
	ui "github.com/gizak/termui/v3"
)

func NewBlock(title string, border bool) *ui.Block {
	block := ui.NewBlock()
	block.Title = title
	block.Border = border
	return block
}
