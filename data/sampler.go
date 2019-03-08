package data

import (
	"github.com/sqshq/sampler/trigger"
	"time"
)

type Sampler struct {
	consumer Consumer
	item     Item
	triggers []trigger.Trigger
}

func NewSampler(consumer Consumer, item Item, triggers []trigger.Trigger, rateMs int) Sampler {

	ticker := time.NewTicker(time.Duration(rateMs * int(time.Millisecond)))
	sampler := Sampler{consumer, item, triggers}

	go func() {
		sampler.sample()
		for ; true; <-ticker.C {
			sampler.sample()
		}
	}()

	return sampler
}

func (s *Sampler) sample() {
	value, err := s.item.nextValue()
	if err == nil {
		sample := Sample{Value: value, Label: s.item.Label}
		s.consumer.SampleChannel <- sample
	} else {
		s.consumer.AlertChannel <- Alert{Text: err.Error()}
	}
}
