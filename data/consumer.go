package data

import (
	ui "github.com/gizak/termui/v3"
	"strings"
	"sync"
)

type Consumer struct {
	SampleChannel  chan *Sample
	AlertChannel   chan *Alert
	CommandChannel chan *Command
	Alert          *Alert
	colorOverrides map[string]map[int]ui.Color
	colorMutex     *sync.Mutex
}

func (c *Consumer) SetColorOverride(label string, trigger int, color *ui.Color) {
	c.colorMutex.Lock()
	defer c.colorMutex.Unlock()
	if color == nil {
		delete(c.colorOverrides[label], trigger)
		return
	}
	if c.colorOverrides[label] == nil {
		c.colorOverrides[label] = make(map[int]ui.Color)
	}
	c.colorOverrides[label][trigger] = *color
}

func (c *Consumer) ColorOverride(label string) (ui.Color, bool) {
	c.colorMutex.Lock()
	defer c.colorMutex.Unlock()
	last := -1
	var color ui.Color
	for trigger, col := range c.colorOverrides[label] {
		if trigger > last {
			last = trigger
			color = col
		}
	}
	return color, last >= 0
}

func (c *Consumer) ColorFor(label string, defaultColor ui.Color) ui.Color {
	if color, ok := c.ColorOverride(label); ok {
		return color
	}
	return defaultColor
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
		colorOverrides: make(map[string]map[int]ui.Color),
		colorMutex:     &sync.Mutex{},
	}
}
