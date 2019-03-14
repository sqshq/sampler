package data

import ui "github.com/gizak/termui/v3"

type Consumer struct {
	SampleChannel chan Sample
	AlertChannel  chan Alert
}

type Sample struct {
	Label string
	Value string
	Color *ui.Color
}

type Alert struct {
	Title string
	Text  string
	Color *ui.Color
}

func NewConsumer() Consumer {
	return Consumer{
		SampleChannel: make(chan Sample),
		AlertChannel:  make(chan Alert),
	}
}
