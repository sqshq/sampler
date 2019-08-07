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
		"",
		"Sampler is an OSS project, and it needs funding to be alive and keep developing",
		"Before the first start, please explore our licensing options below. For more details, visit WWW.SAMPLER.DEV",
		"", "", "",
		"How do you plan to use Sampler?",
	}...)

	commericalText := append(util.AsciiLogo, []string{
		"", "", "", "",
		"With Sampler, you can easily save time and solve some of your business problems.",
		"That's why support of the project is in the interest of your organization.",
		"",
		"",
		"We are offering commercial licenses which provide priority support and technical assistance.",
		"After entering the licence key, your company name will appear in the status bar.",
		"",
		"",
		"To make a purchase, please visit WWW.SAMPLER.DEV",
	}...)

	personalText := append(util.AsciiLogo, []string{
		"", "", "", "",
		"Sampler is always free to use, but you can support the project and donate any amount to get a personal license.",
		"Once it is activated, your name will appear in the status bar.",
		"",
		"",
		"To become a sponsor, please visit WWW.SAMPLER.DEV",
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

		buffer.SetString(buttonCommercial, commercialButtonStyle,
			util.GetMiddlePoint(intro.Block.Rectangle, buttonCommercial, 5))
		buffer.SetString(buttonPersonal, personalButtonStyle,
			util.GetMiddlePoint(intro.Block.Rectangle, buttonPersonal, 7))
	} else {
		buffer.SetString(buttonOk, highlightedStyle,
			util.GetMiddlePoint(intro.Block.Rectangle, buttonOk, 7))
	}

	intro.Block.Draw(buffer)
}
