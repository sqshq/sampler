package util

import (
	"fmt"
	"math"
	"strconv"
)

func FormatValue(value float64, scale int) string {
	if math.Abs(value) == math.MaxFloat64 {
		return "Inf"
	} else {
		format := "%." + strconv.Itoa(scale) + "f"
		return fmt.Sprintf(format, value)
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
