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
		for {
			select {
			case <-ticker.C:
				sampler.sample()
			}
		}
	}()

	return sampler
}

func (self *Sampler) sample() {

	value, err := self.item.nextValue()

	sample := Sample{
		Value: value,
		Error: err,
		Color: self.item.Color,
		Label: self.item.Label,
	}

	self.consumer.ConsumeSample(sample)
}
