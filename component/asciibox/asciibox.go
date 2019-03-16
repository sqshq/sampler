package asciibox

import (
	ui "github.com/gizak/termui/v3"
	fl "github.com/mbndr/figlet4go"
	"github.com/sqshq/sampler/asset"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/config"
	"image"
)

type AsciiBox struct {
	*component.Component
	text    string
	ascii   string
	style   ui.Style
	render  *fl.AsciiRender
	options *fl.RenderOptions
}

const asciiFontExtension = ".flf"

func NewAsciiBox(c config.AsciiBoxConfig) *AsciiBox {

	options := fl.NewRenderOptions()
	options.FontName = string(*c.Font)

	fontStr, err := asset.Asset(options.FontName + asciiFontExtension)
	if err != nil {
		panic("Can't load the font: " + err.Error())
	}

	render := fl.NewAsciiRender()
	_ = render.LoadBindataFont(fontStr, options.FontName)

	box := AsciiBox{
		Component: component.NewComponent(c.ComponentConfig, config.TypeAsciiBox),
		style:     ui.NewStyle(*c.Color),
		render:    render,
		options:   options,
	}

	go func() {
		for {
			select {
			case sample := <-box.SampleChannel:
				box.text = sample.Value
				box.ascii, _ = box.render.RenderOpts(sample.Value, box.options)
			}
		}
	}()

	return &box
}

func (a *AsciiBox) Draw(buffer *ui.Buffer) {

	buffer.Fill(ui.NewCell(' ', ui.NewStyle(ui.ColorBlack)), a.GetRect())
	a.Block.Draw(buffer)

	point := a.Inner.Min
	cells := ui.ParseStyles(a.ascii, a.style)

	for i := 0; i < len(cells) && point.Y < a.Inner.Max.Y; i++ {
		if cells[i].Rune == '\n' {
			point = image.Pt(a.Inner.Min.X, point.Y+1)
		} else if point.In(a.Inner) {
			buffer.SetCell(cells[i], point)
			point = point.Add(image.Pt(1, 0))
		}
	}
}
