package util

import (
	"image"
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

func GetRectCoordinates(area image.Rectangle, width int, height int) (int, int, int, int) {
	x1 := area.Min.X + area.Dx()/2 - width/2
	y1 := area.Min.Y + area.Dy()/2 - height
	return x1, y1, x1 + width, y1 + height + 2
}
