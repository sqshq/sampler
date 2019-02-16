package data

type Consumer interface {
	ConsumeSample(sample Sample)
}

type Sample struct {
	Label string
	Value string
	Error error
}
