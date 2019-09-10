package console

import (
	"fmt"
	"runtime"

	ui "github.com/gizak/termui/v3"
)

type Theme string

const (
	ThemeDark  Theme = "dark"
	ThemeLight Theme = "light"
)

const (
	ColorOlive       ui.Color = 178
	ColorDeepSkyBlue ui.Color = 39
	ColorDeepPink    ui.Color = 198
	ColorCian        ui.Color = 43
	ColorOrange      ui.Color = 166
	ColorPurple      ui.Color = 129
	ColorGreen       ui.Color = 64
	ColorDarkRed     ui.Color = 88
	ColorBlueViolet  ui.Color = 57
	ColorDarkGrey    ui.Color = 238
	ColorLightGrey   ui.Color = 254
	ColorGrey        ui.Color = 242
	ColorWhite       ui.Color = 15
	ColorBlack       ui.Color = 0
	ColorClear       ui.Color = -1
)

const (
	menuColorNix            ui.Color = 255
	menuColorReverseNix     ui.Color = 235
	menuColorWindows        ui.Color = 255
	menuColorReverseWindows ui.Color = 0
)

type Palette struct {
	ContentColors  []ui.Color
	GradientColors [][]ui.Color
	BaseColor      ui.Color
	MediumColor    ui.Color
	ReverseColor   ui.Color
}

// GetPalette returns a color palette based on specified theme
func GetPalette(theme Theme) Palette {
	switch theme {
	case ThemeDark:
		return Palette{
			ContentColors:  []ui.Color{ColorOlive, ColorDeepSkyBlue, ColorDeepPink, ColorWhite, ColorGrey, ColorGreen, ColorOrange, ColorCian, ColorPurple},
			GradientColors: [][]ui.Color{{39, 33, 62, 93, 164, 161}, {95, 138, 180, 179, 178, 178}},
			BaseColor:      ColorWhite,
			MediumColor:    ColorDarkGrey,
			ReverseColor:   ColorBlack,
		}
	case ThemeLight:
		return Palette{
			ContentColors:  []ui.Color{ColorBlack, ColorDarkRed, ColorBlueViolet, ColorGrey, ColorGreen},
			GradientColors: [][]ui.Color{{250, 248, 246, 244, 242, 240, 238, 236, 234, 232, 16}},
			BaseColor:      ColorBlack,
			MediumColor:    ColorLightGrey,
			ReverseColor:   ColorWhite,
		}
	default:
		panic(fmt.Sprintf("Specified theme is not supported: %v", theme))
	}
}

// GetMenuColor returns a color based on the
// operating system target
func GetMenuColor() ui.Color {
	switch runtime.GOOS {
	case "windows":
		return menuColorWindows
	default:
		return menuColorNix
	}
}

// GetMenuColorReverse returns a color based on the
// operating system target
func GetMenuColorReverse() ui.Color {
	switch runtime.GOOS {
	case "windows":
		return menuColorReverseWindows
	default:
		return menuColorReverseNix
	}
}

func GetGradientColor(gradient []ui.Color, cur int, max int) ui.Color {
	ratio := float64(len(gradient)) / float64(max)
	index := int(ratio * float64(cur))
	if index > len(gradient)-1 {
		index = len(gradient) - 1
	}
	return gradient[index]
}
