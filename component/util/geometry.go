package util

import (
	"image"
	"math"
)

func GetRectLeftSideCenter(rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Min.X,
		Y: rect.Min.Y + rect.Dy()/2,
	}
}

func GetRectRightSideCenter(rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Max.X,
		Y: rect.Min.Y + rect.Dy()/2,
	}
}

func GetRectTopSideCenter(rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Min.X + rect.Dx()/2,
		Y: rect.Min.Y,
	}
}

func GetRectBottomSideCenter(rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Min.X + rect.Dx()/2,
		Y: rect.Max.Y,
	}
}

func GetDistance(p1 image.Point, p2 image.Point) float64 {
	x := math.Abs(float64(p1.X - p2.X))
	y := math.Abs(float64(p1.Y - p2.Y))
	return math.Sqrt(x*x + y*y)
}

func GetRectCoordinates(area image.Rectangle, width int, height int) (int, int, int, int) {
	x1 := area.Min.X + area.Dx()/2 - width/2
	y1 := area.Min.Y + area.Dy()/2 - height
	return x1, y1, x1 + width, y1 + height + 2
}
