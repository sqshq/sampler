package util

import "testing"
import "image"

func TestGetRectLeftSideCenter(t *testing.T) {
	rect := image.Rect(10, 10, 20, 20)
	expected := image.Point{
		X: 10,
		Y: 15,
	}
	result := GetRectLeftSideCenter(rect)
	if result != expected {
		t.Errorf("GetRectLeftSideCenter was incorrect. Expected %v, got %v.", expected, result)
	}
}

func TestGetRectRightSideCenter(t *testing.T) {
	rect := image.Rect(10, 10, 20, 21)
	expected := image.Point{
		X: 20,
		Y: 15,
	}
	result := GetRectRightSideCenter(rect)
	if result != expected {
		t.Errorf("GetRectRightSideCenter was incorrect. Expected %v, got %v.", expected, result)
	}
}

func TestGetRectTopSideCenter(t *testing.T) {
	rect := image.Rect(10, 10, 20, 20)
	expected := image.Point{
		X: 15,
		Y: 10,
	}
	result := GetRectTopSideCenter(rect)
	if result != expected {
		t.Errorf("GetRectTopSideCenter was incorrect. Expected %v, got %v.", expected, result)
	}
}

func TestGetRectBottomSideCenter(t *testing.T) {
	rect := image.Rect(10, 9, 20, 20)
	expected := image.Point{
		X: 15,
		Y: 20,
	}
	result := GetRectBottomSideCenter(rect)
	if result != expected {
		t.Errorf("GetRectBottomSideCenter was incorrect. Expected %v, got %v.", expected, result)
	}
}
