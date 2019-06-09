package component

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/component/util"
	"github.com/sqshq/sampler/console"
)

type Intro struct {
	*ui.Block
	page    IntroPage
	option  introOption
	palette console.Palette
}

type IntroPage rune

const (
	IntroPageWelcome    IntroPage = 0
	IntroPageCommercial IntroPage = 1
	IntroPagePersonal   IntroPage = 2
)

type introOption rune

const (
	introOptionCommercial introOption = 0
	introOptionPersonal   introOption = 1
)

const (
	buttonCommercial string = "            COMMERCIAL USE            "
	buttonPersonal   string = "             PERSONAL USE             "
	buttonOk         string = "                 OK                   "
)

func (intro *Intro) Up() {
	intro.option = introOptionCommercial
}

func (intro *Intro) Down() {
	intro.option = introOptionPersonal
}

func (intro *Intro) NextPage() {
	if intro.option == introOptionCommercial {
		intro.page = IntroPageCommercial
	} else {
		intro.page = IntroPagePersonal
	}
}

func (intro *Intro) GetSelectedPage() IntroPage {
	return intro.page
}

func NewIntro(palette console.Palette) *Intro {
	return &Intro{
		Block:   NewBlock("", false, palette),
		palette: palette,
	}
}

func (intro *Intro) Draw(buffer *ui.Buffer) {

	introText := append(util.AsciiLogo, []string{
		"", "", "",
		"Welcome.",
		"Sampler is free of charge for personal use, but license must be purchased to use it for business purposes.",
		"By proceeding, you agree to the terms of the license agreement and privacy policy: www.sampler.dev/license",
		"", "", "",
		"How do you plan to use Sampler?",
	}...)

	commericalText := append(util.AsciiLogo, []string{
		"", "", "", "",
		"Please visit www.sampler.dev to purchase a license and then start Sampler with --license flag",
	}...)

	personalText := append(util.AsciiLogo, []string{
		"", "", "", "",
		"Sampler is always free for non-commercial use, but you can support the project and buy a personal license:",
		"www.sampler.dev",
	}...)

	text := introText

	switch intro.page {
	case IntroPageWelcome:
		text = introText
	case IntroPageCommercial:
		text = commericalText
	case IntroPagePersonal:
		text = personalText
	}

	for i, a := range text {
		util.PrintString(
			a,
			ui.NewStyle(intro.palette.BaseColor),
			util.GetMiddlePoint(intro.Block.Rectangle, a, i-15),
			buffer)
	}

	highlightedStyle := ui.NewStyle(intro.palette.ReverseColor, intro.palette.BaseColor)
	regularStyle := ui.NewStyle(intro.palette.BaseColor, intro.palette.ReverseColor)

	if intro.page == IntroPageWelcome {

		commercialButtonStyle := highlightedStyle
		if intro.option == introOptionPersonal {
			commercialButtonStyle = regularStyle
		}

		personalButtonStyle := highlightedStyle
		if intro.option == introOptionCommercial {
			personalButtonStyle = regularStyle
		}

		buffer.SetString(string(buttonCommercial), commercialButtonStyle,
			util.GetMiddlePoint(intro.Block.Rectangle, string(buttonCommercial), 6))
		buffer.SetString(string(buttonPersonal), personalButtonStyle,
			util.GetMiddlePoint(intro.Block.Rectangle, string(buttonPersonal), 8))
	} else {
		buffer.SetString(string(buttonOk), highlightedStyle,
			util.GetMiddlePoint(intro.Block.Rectangle, string(buttonOk), 4))
	}

	intro.Block.Draw(buffer)
}
