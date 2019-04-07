package data

import (
	"github.com/sqshq/sampler/config"
	"time"
)

type Sampler struct {
	consumer        *Consumer
	items           []*Item
	triggers        []*Trigger
	triggersChannel chan *Sample
}

func NewSampler(consumer *Consumer, items []*Item, triggers []*Trigger, options config.Options, rateMs int) Sampler {

	ticker := time.NewTicker(time.Duration(rateMs * int(time.Millisecond)))

	sampler := Sampler{
		consumer,
		items,
		triggers,
		make(chan *Sample),
	}

	go func() {
		for ; true; <-ticker.C {
			for _, item := range sampler.items {
				go sampler.sample(item, options)
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

func (s *Sampler) sample(item *Item, options config.Options) {

	val, err := item.nextValue(options.Variables)

	if len(val) > 0 {
		sample := &Sample{Label: item.Label, Value: val, Color: item.Color}
		s.consumer.SampleChannel <- sample
		s.triggersChannel <- sample
	} else if err != nil {
		s.consumer.AlertChannel <- &Alert{
			Title: "SAMPLING FAILURE",
			Text:  getErrorMessage(err),
			Color: item.Color,
		}
	}
}
