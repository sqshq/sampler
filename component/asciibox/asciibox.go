package asciibox

import (
	fl "github.com/mbndr/figlet4go"
	"github.com/sqshq/sampler/data"
	ui "github.com/sqshq/termui"
	"image"
)

type AsciiBox struct {
	ui.Block
	text    string
	ascii   string
	style   ui.Style
	render  *fl.AsciiRender
	options *fl.RenderOptions
}

type AsciiFont string

const (
	AsciiFontFlat AsciiFont = "flat"
	AsciiFont3D   AsciiFont = "3d"
)

func NewAsciiBox(title string, font AsciiFont, color ui.Color) *AsciiBox {

	block := *ui.NewBlock()
	block.Title = title

	render := fl.NewAsciiRender()
	err := render.LoadFont("asset/")
	if err != nil {
		panic("Can't load fonts: " + err.Error())
	}

	options := fl.NewRenderOptions()
	options.FontName = string(font)

	return &AsciiBox{
		Block:   block,
		style:   ui.NewStyle(color),
		render:  render,
		options: options,
	}
}

func (a *AsciiBox) ConsumeSample(sample data.Sample) {
	a.text = sample.Value
	a.ascii, _ = a.render.RenderOpts(sample.Value, a.options)
}

func (a *AsciiBox) Draw(buffer *ui.Buffer) {

	buffer.Fill(ui.NewCell(' ', ui.NewStyle(ui.ColorBlack)), a.GetRect())
	a.Block.Draw(buffer)

	point := a.Inner.Min
	cells := ui.ParseText(a.ascii, a.style)

	for i := 0; i < len(cells) && point.Y < a.Inner.Max.Y; i++ {
		if cells[i].Rune == '\n' {
			point = image.Pt(a.Inner.Min.X, point.Y+1)
		} else if point.In(a.Inner) {
			buffer.SetCell(cells[i], point)
			point = point.Add(image.Pt(1, 0))
		}
	}
}
