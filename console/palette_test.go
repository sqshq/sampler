package console

import "testing"

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
		want  Palette
	}{
		{"should return dark theme with base color white", ThemeDark, darkPalette},
		{"should return light theme with base color black", ThemeLight, lightPalette},
	}

	for _, test := range tests {
		palette := GetPalette(test.input)
		if got := palette.BaseColor; got != test.want.BaseColor {
			t.Errorf("GetPalette(%q) = %d, want %d", test.input, got, test.want.BaseColor)
		}
	}
}
