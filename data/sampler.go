package data

import (
	"time"
)

type Sampler struct {
	consumer        Consumer
	items           []Item
	triggers        []Trigger
	triggersChannel chan Sample
}

func NewSampler(consumer Consumer, items []Item, triggers []Trigger, rateMs int) Sampler {

	ticker := time.NewTicker(
		time.Duration(rateMs * int(time.Millisecond)),
	)

	sampler := Sampler{
		consumer,
		items,
		triggers,
		make(chan Sample),
	}

	go func() {
		for ; true; <-ticker.C {
			for _, item := range sampler.items {
				go sampler.sample(item)
			}
		}
	}()

	go func() {
		for {
			select {
			case sample := <-sampler.triggersChannel:
				for _, t := range sampler.triggers {
					t.Execute(sample)
				}
			}
		}
	}()

	return sampler
}

func (s *Sampler) sample(item Item) {

	val, err := item.nextValue()

	if err == nil {
		sample := Sample{Label: item.Label, Value: val, Color: item.Color}
		s.consumer.SampleChannel <- sample
		s.triggersChannel <- sample
	} else {
		s.consumer.AlertChannel <- Alert{
			Title: "SAMPLING FAILURE",
			Text:  err.Error(),
			Color: item.Color,
		}
	}
}
