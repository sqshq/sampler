package data

import . "github.com/sqshq/termui"

type Consumer interface {
	ConsumeSample(sample Sample)
}

type Sample struct {
	Label string
	Color Color
	Value string
	Error error
}
