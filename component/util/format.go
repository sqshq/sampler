package util

import (
	"bytes"
	"math"
	"strconv"
	"strings"
)

func FormatValue(value float64, scale int) string {
	if math.Abs(value) == math.MaxFloat64 {
		return "Inf"
	} else {
		return formatTrailingDigits(addRadixChars(value), scale)
	}
}

func FormatValueWithSign(value float64, scale int) string {
	if value == 0 {
		return " 0"
	} else if value > 0 {
		return "+" + FormatValue(value, scale)
	} else {
		return FormatValue(value, scale)
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

		return strings.TrimRight(formatted, "0.")
	}

	return value
}
