package data

import (
	"time"
)

type Sampler struct {
	consumer Consumer
	item     Item
}

func NewSampler(consumer Consumer, item Item, rateMs int) Sampler {

	ticker := time.NewTicker(time.Duration(rateMs * int(time.Millisecond)))
	sampler := Sampler{consumer, item}

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
	sample := Sample{Value: value, Error: err, Label: *s.item.Label}
	s.consumer.ConsumeSample(sample)
}
