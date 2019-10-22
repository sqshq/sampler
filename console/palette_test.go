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
	type testArgs struct {
		color    ui.Color
		cur      int
		max      int
		want     ui.Color
	}

	var tests = []testArgs{
		{
			color: 2,
			cur: 55,
			max: 60,
			want: 39,
		},
		{
			color: 6,
			cur: 12,
			max: 60,
			want: 29,
		},
		{
			color: 233,
			cur: 19,
			max: 60,
			want: 233,
		},
		{
			color: 109,
			cur: 8,
			max: 60,
			want: 106,
		},
		{
			color: 14,
			cur: 34,
			max: 60,
			want: 49,
		},
		{
			color: 201,
			cur: 10,
			max: 20,
			want: 199,
		},
		{
			color: 201,
			cur: 15,
			max: 30,
			want: 199,
		},
		{
			color: 201,
			cur: 82,
			max: 260,
			want: 197,
		},
		{
			color: 201,
			cur: 9388,
			max: 100032,
			want: 196,
		},
		{
			color: 201,
			cur: 2,
			max: 302,
			want: 196,
		},
	}

	for i := 0; i<len(tests); i++ {
		gradientColor := GetGradientColor(
			tests[i].color,
			tests[i].cur,
			tests[i].max,
		)
		if got := gradientColor; got != tests[i].want {
			t.Errorf("GetGradientColor(%d, %d, %d) = %d, want %d", tests[i].color, tests[i].cur, tests[i].max, got, tests[i].want)
		}
	}
}
