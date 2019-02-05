package data

import (
	"time"
)

type Sampler struct {
	consumer Consumer
	item     Item
}

func NewSampler(consumer Consumer, item Item, rateMs int) Sampler {

	sampler := Sampler{consumer, item}

	go func() {
		for t := time.Tick(time.Duration(rateMs * int(time.Millisecond))); ; <-t {
			sampler.sample()
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
