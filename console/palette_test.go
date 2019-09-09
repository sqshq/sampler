package console

import (
	"testing"

	ui "github.com/gizak/termui/v3"
)

func TestGetPalette(t *testing.T) {
	var (
		darkPalette = Palette{
			BaseColor:    ColorWhite,
			ReverseColor: ColorBlack,
		}
		lightPalette = Palette{
			BaseColor:    ColorBlack,
			ReverseColor: ColorWhite,
		}
	)

	tests := []struct {
		name  string
		input Theme
		Palette
	}{
		{"should return dark theme with base color white", ThemeDark, darkPalette},
		{"should return light theme with base color black", ThemeLight, lightPalette},
	}

	for _, test := range tests {
		palette := GetPalette(test.input)
		if got := palette.BaseColor; got != test.BaseColor {
			t.Errorf("GetPalette(%q) = %d, want %d", test.input, got, test.BaseColor)
		}
	}
}

func TestGetPaletteInvalidTheme(t *testing.T) {
	const invalid Theme = "invalid"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetPalette(%q) should have panicked", invalid)
		}
	}()

	GetPalette(invalid)
}

func TestGetGradientColor(t *testing.T) {
	type args struct {
		gradient []ui.Color
		cur      int
		max      int
	}

	var (
		lightThemeGradientInput = args{
			gradient: []ui.Color{
				250, 248, 246, 244, 242, 240, 238, 236, 234, 232, 16,
			},
			cur: 200,
			max: 250,
		}

		darkThemeGradientInput = args{
			gradient: []ui.Color{
				39, 33, 62, 93, 164, 161,
			},
			cur: 40,
			max: 180,
		}

		grey ui.Color = 234

		blue ui.Color = 33
	)

	tests := []struct {
		name string
		args
		want ui.Color
	}{
		{"should return color grey", lightThemeGradientInput, grey},
		{"should return color blue", darkThemeGradientInput, blue},
	}

	for _, test := range tests {
		gradientColor := GetGradientColor(
			test.gradient,
			test.cur,
			test.max,
		)

		if got := gradientColor; got != test.want {
			t.Errorf("GetGradientColor(%v) = %d, want %d", test.args, got, test.want)
		}
	}
}
