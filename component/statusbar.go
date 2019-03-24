package component

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/console"
	"image"
)

const (
	bindingsIndent = 4
)

type StatusBar struct {
	*ui.Block
	keyBindings    []string
	configFileName string
}

func NewStatusLine(configFileName string, palette console.Palette) *StatusBar {
	block := *ui.NewBlock()
	block.Border = false
	return &StatusBar{
		Block:          NewBlock("", false, palette),
		configFileName: configFileName,
		keyBindings: []string{
			"(Q) quit",
			"(P) pause",
			"(<->) selection",
			"(ESC) reset alerts",
		},
	}
}

func (s *StatusBar) Draw(buffer *ui.Buffer) {
	buffer.Fill(ui.NewCell(' ', ui.NewStyle(console.ColorClear, console.MenuColorBackground)), s.GetRect())
	buffer.SetString(fmt.Sprintf(" %s %s @ %s", console.AppTitle, console.AppVersion, s.configFileName), ui.NewStyle(console.MenuColorText, console.MenuColorBackground), s.Min)

	indent := bindingsIndent
	for _, binding := range s.keyBindings {
		buffer.SetString(binding, ui.NewStyle(console.MenuColorText, console.MenuColorBackground), image.Pt(s.Max.X-len(binding)-indent, s.Min.Y))
		indent += bindingsIndent + len(binding)
	}

	s.Block.Draw(buffer)
}
