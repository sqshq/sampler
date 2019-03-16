package data

import ui "github.com/gizak/termui/v3"

// TODO interface here, move fields declaration in the Component
type Consumer struct {
	SampleChannel  chan Sample
	AlertChannel   chan Alert
	CommandChannel chan Command
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

func NewConsumer() Consumer {
	return Consumer{
		SampleChannel:  make(chan Sample),
		AlertChannel:   make(chan Alert),
		CommandChannel: make(chan Command),
	}
}
