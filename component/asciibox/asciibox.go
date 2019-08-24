package asciibox

import (
	ui "github.com/gizak/termui/v3"
	fl "github.com/mbndr/figlet4go"
	"github.com/sqshq/sampler/asset"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"image"
	"strings"
)

// AsciiBox represents a component with ascii-style text
type AsciiBox struct {
	*ui.Block
	*data.Consumer
	alert   *data.Alert
	ascii   string
	style   ui.Style
	render  *fl.AsciiRender
	options *fl.RenderOptions
	palette console.Palette
}

const asciiFontExtension = ".flf"

func NewAsciiBox(c config.AsciiBoxConfig, palette console.Palette) *AsciiBox {

	options := fl.NewRenderOptions()
	options.FontName = string(*c.Font)

	fontStr, err := asset.Asset(options.FontName + asciiFontExtension)
	if err != nil {
		panic("Failed to load the font: " + err.Error())
	}

	render := fl.NewAsciiRender()
	_ = render.LoadBindataFont(fontStr, options.FontName)

	color := c.Color
	if color == nil {
		color = &palette.BaseColor
	}

	box := AsciiBox{
		Block:    component.NewBlock(c.Title, *c.Border, palette),
		Consumer: data.NewConsumer(),
		style:    ui.NewStyle(*color),
		render:   render,
		options:  options,
		palette:  palette,
	}

	go func() {
		for {
			select {
			case sample := <-box.SampleChannel:
				box.renderText(sample)
			case alert := <-box.AlertChannel:
				box.alert = alert
			}
		}
	}()

	return &box
}

func (a *AsciiBox) renderText(sample *data.Sample) {
	text := strings.TrimSpace(sample.Value)
	lines := strings.Split(text, "\n")
	a.ascii = ""
	for _, line := range lines {
		ascii, _ := a.render.RenderOpts(line, a.options)
		a.ascii += ascii + "\n"
	}
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

	component.RenderAlert(a.alert, a.Rectangle, buffer)
}
