package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component/util"
	"github.com/sqshq/sampler/console"
)

type NagWindow struct {
	*ui.Block
	palette  console.Palette
	accepted bool
}

func NewNagWindow(palette console.Palette) *NagWindow {
	return &NagWindow{
		Block:    NewBlock("", false, palette),
		palette:  palette,
		accepted: false,
	}
}

func (n *NagWindow) Accept() {
	n.accepted = true
}

func (n *NagWindow) IsAccepted() bool {
	return n.accepted
}

func (n *NagWindow) Draw(buffer *ui.Buffer) {

	text := append(util.AsciiLogo, []string{
		"", "", "",
		"Thank you for using Sampler.",
		"It is always free, but you can sponsor the project and buy the personal or commercial license",
		"for priority support and technical assistance.",
		"",
		"Please visit www.sampler.dev",
	}...)

	for i, a := range text {
		util.PrintString(
			a,
			ui.NewStyle(n.palette.BaseColor),
			util.GetMiddlePoint(n.Block.Rectangle, a, i-15),
			buffer)
	}

	buffer.SetString(buttonOk, ui.NewStyle(n.palette.ReverseColor, n.palette.BaseColor),
		util.GetMiddlePoint(n.Block.Rectangle, buttonOk, 4))

	n.Block.Draw(buffer)
}
