package util

import (
	"strconv"
	"strings"
)

func ParseFloat(input string) (float64, error) {

	clean := strings.TrimSpace(input)

	if strings.Contains(clean, "\n") {
		lastIndex := strings.LastIndex(clean, "\n")
		clean = clean[lastIndex+1:]
	}

	return strconv.ParseFloat(clean, 64)
}
