package data

type Consumer struct {
	SampleChannel chan Sample
	AlertChannel  chan Alert
}

type Sample struct {
	Label string
	Value string
}

type Alert struct {
	Title string
	Text  string
}

func NewConsumer() Consumer {
	return Consumer{
		SampleChannel: make(chan Sample),
		AlertChannel:  make(chan Alert),
	}
}
