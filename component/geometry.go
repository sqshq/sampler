package component

import (
	"image"
	"math"
)

func GetRectLeftAgeCenter(rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Min.X,
		Y: rect.Min.Y + rect.Dy()/2,
	}
}

func GetRectRightAgeCenter(rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Max.X,
		Y: rect.Min.Y + rect.Dy()/2,
	}
}

func GetRectTopAgeCenter(rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Min.X + rect.Dx()/2,
		Y: rect.Min.Y,
	}
}

func GetRectBottomAgeCenter(rect image.Rectangle) image.Point {
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
