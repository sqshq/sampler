package data

import ui "github.com/gizak/termui/v3"

type Consumer struct {
	SampleChannel  chan *Sample
	AlertChannel   chan *Alert
	CommandChannel chan *Command
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

type Command struct {
	Type  string
	Value interface{}
}

func NewConsumer() *Consumer {
	return &Consumer{
		SampleChannel:  make(chan *Sample, 10),
		AlertChannel:   make(chan *Alert, 10),
		CommandChannel: make(chan *Command, 10),
	}
}
