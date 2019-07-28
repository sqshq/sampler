package util

import (
	"bytes"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"image"
	"math"
	"strconv"
	"strings"
)

var AsciiLogo = []string{
	"                               __         ",
	"   _________  ____ ___  ____  / /__  _____",
	"  / ___/ __ `/ __ `__ \\/ __ \\/ / _ \\/ ___/",
	" (__  ) /_/ / / / / / / /_/ / /  __/ /    ",
	"/____/\\__,_/_/ /_/ /_/ .___/_/\\___/_/     ",
	"                    /_/                   ",
}

func FormatValue(value float64, scale int) string {
	if math.Abs(value) == math.MaxFloat64 {
		return "Inf"
	}
	return formatTrailingDigits(addRadixChars(value), scale)
}

func FormatDelta(value float64, scale int) string {

	abs := math.Abs(value)
	val := value
	scl := scale

	postfix := ""

	if abs > 1000 && abs < 1000000 {
		val = float64(value) / 1000
		postfix = "k"
	} else if abs > 1000000 && abs < 1000000000 {
		val = float64(value) / 1000000
		postfix = "M"
	} else if abs > 1000000000 {
		val = float64(value) / 1000000000
		postfix = "B"
	}

	if abs > 1000 {
		scl = 1
	}

	if val == 0 {
		return " 0"
	} else if val > 0 {
		return fmt.Sprintf("+%s%s", FormatValue(val, scl), postfix)
	} else {
		return fmt.Sprintf("%s%s", FormatValue(val, scl), postfix)
	}
}

func addRadixChars(value float64) string {
	buf := &bytes.Buffer{}
	if value < 0 {
		buf.Write([]byte{'-'})
		value = 0 - value
	}

	radix := []byte{','}

	parts := strings.Split(strconv.FormatFloat(value, 'f', -1, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(radix)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(radix)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

func formatTrailingDigits(value string, scale int) string {

	if i := strings.Index(value, "."); i >= 0 {

		formatted := value

		if scale <= 0 {
			formatted = value[:i]
		}

		i++
		if i+scale < len(value) {
			formatted = value[:i+scale]
		}

		formatted = strings.TrimRight(formatted, "0")

		return strings.TrimRight(formatted, ".")
	}

	return value
}

func GetMiddlePoint(rectangle image.Rectangle, text string, offset int) image.Point {
	return image.Pt(rectangle.Min.X+rectangle.Dx()/2-len(text)/2, rectangle.Max.Y-rectangle.Dy()/2+offset)
}

func PrintString(s string, style ui.Style, p image.Point, buffer *ui.Buffer) {
	for i, char := range s {
		buffer.SetCell(ui.Cell{Rune: char, Style: style}, image.Pt(p.X+i, p.Y))
	}
}
