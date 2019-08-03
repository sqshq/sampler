package util

import (
	"strconv"
	"strings"
)

func ParseFloat(input string) (float64, error) {

	clean := strings.TrimSpace(input)
	clean = strings.Replace(clean, ",", ".", -1) // replace decimal comma with decimal point

	if strings.Contains(clean, "\n") {
		lastIndex := strings.LastIndex(clean, "\n")
		clean = clean[lastIndex+1:]
	}

	return strconv.ParseFloat(clean, 64)
}
