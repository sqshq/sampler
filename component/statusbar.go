package component

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/metadata"
	"image"
)

const (
	pauseText      = "  P A U S E D  "
	bindingsIndent = 3
)

type StatusBar struct {
	*ui.Block
	keyBindings []string
	text        string
	pause       bool
}

func NewStatusBar(configFileName string, palette console.Palette, license *metadata.License) *StatusBar {

	block := *ui.NewBlock()
	block.Border = false
	text := fmt.Sprintf(" %s %s | ", console.AppTitle, console.AppVersion)

	if license == nil || !license.Valid || license.Type == nil {
		text += console.AppLicenseWarning
	} else if *license.Type == metadata.TypePersonal {
		text += fmt.Sprintf("%s | personal license: %s", configFileName, *license.Username)
	} else if license.Username != nil {
		text += fmt.Sprintf("%s | licensed to %s", configFileName, *license.Username)
		if license.Company != nil {
			text += fmt.Sprintf(", %s", *license.Company)
		}
	} else {
		text += fmt.Sprintf("%s | licensed to %s", configFileName, *license.Company)
	}

	return &StatusBar{
		Block: NewBlock("", false, palette),
		text:  text,
		keyBindings: []string{
			"(q) quit",
			"(p) pause",
			"(<->) selection",
			"(ESC) reset alerts",
		},
	}
}

func (s *StatusBar) Draw(buffer *ui.Buffer) {

	buffer.Fill(ui.NewCell(' ', ui.NewStyle(console.ColorClear, console.GetMenuColorReverse())), s.GetRect())

	indent := bindingsIndent
	for _, binding := range s.keyBindings {
		buffer.SetString(binding, ui.NewStyle(console.GetMenuColor(), console.GetMenuColorReverse()), image.Pt(s.Max.X-len(binding)-indent, s.Min.Y))
		indent += bindingsIndent + len(binding)
	}

	buffer.SetString(s.text, ui.NewStyle(console.GetMenuColor(), console.GetMenuColorReverse()), s.Min)

	if s.pause {
		buffer.SetString(pauseText, ui.NewStyle(console.GetMenuColorReverse(), console.GetMenuColor()), image.Pt(s.Max.X-s.Dx()/2-len(pauseText)/2, s.Min.Y))
	}

	s.Block.Draw(buffer)
}

func (s *StatusBar) TogglePause() {
	s.pause = !s.pause
}
