package widgets

import (
	"image"
)

const (
	xBrailleMultiplier = 2
	yBrailleMultiplier = 4
)

func braille(point image.Point) image.Point {
	return image.Point{X: point.X * xBrailleMultiplier, Y: point.Y * yBrailleMultiplier}
}

func deBraille(point image.Point) image.Point {
	return image.Point{X: point.X / xBrailleMultiplier, Y: point.Y / yBrailleMultiplier}
}
