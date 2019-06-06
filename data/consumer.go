package data

import (
	ui "github.com/gizak/termui/v3"
	"strings"
)

type Consumer struct {
	SampleChannel  chan *Sample
	AlertChannel   chan *Alert
	CommandChannel chan *Command
	Alert          *Alert
}

func (c *Consumer) HandleConsumeSuccess() {
	if c.Alert != nil && c.Alert.Recoverable {
		c.Alert = nil
	}
}

func (c *Consumer) HandleConsumeFailure(title string, err error, sample *Sample) {
	c.AlertChannel <- &Alert{
		Title:       strings.ToUpper(title),
		Text:        err.Error(),
		Color:       sample.Color,
		Recoverable: true,
	}
}

type Sample struct {
	Label string
	Value string
	Color *ui.Color
}

type Alert struct {
	Title       string
	Text        string
	Color       *ui.Color
	Recoverable bool
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
