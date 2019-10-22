package console

import (
	"fmt"
	"math"
	"runtime"

	ui "github.com/gizak/termui/v3"
)

type Theme string

const (
	ThemeDark  Theme = "dark"
	ThemeLight Theme = "light"
)

const (
	GradientCount    ui.Color = 6  // How many colors in a gradient
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

// GetGradientColor returns a color based on an input color
// and two cur, max values. The function tries to find a range
// of 6 (GradientCount) colors that look similar but are in a
// gradient. It divids the colorspace from color 16 to 255 in
// 40 sections consisting of 6 colors each.
// If the input color is in the 0-15 range it is mapped to it's
// equivalent in the upper (16-255) range.
func GetGradientColor(color ui.Color, cur int, max int) ui.Color {
	// Remap lower 16 colors to higher values that are equivalent
	if color < 16 {
		ColorMap := [16]ui.Color {16, 88, 34, 208, 19, 53, 30, 250, 244, 196, 46, 226, 21, 201, 51, 255}
		color = ColorMap[color]
	}

	// Find the lowest id in the range for the gradient
	// selected by color
	baseColor := ui.Color(((color-16) / GradientCount) * GradientCount + 16)

	// Calculate the offset in the range starting from baseColor
	offset := int(math.Min(float64(6) / float64(max) * float64(cur), 6))

	return baseColor + ui.Color(offset)
}
