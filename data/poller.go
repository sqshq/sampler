package data

import (
	"time"
)

type Poller struct {
	consumer Consumer
	item     Item
}

func NewPoller(consumer Consumer, item Item, rateMs int) Poller {

	ticker := time.NewTicker(time.Duration(rateMs * int(time.Millisecond)))
	poller := Poller{consumer, item}

	go func() {
		for {
			select {
			case <-ticker.C:
				poller.poll()
			}
		}
	}()

	return poller
}

func (self *Poller) poll() {

	value, err := self.item.nextValue()

	sample := Sample{
		Value: value,
		Error: err,
		Color: self.item.Color,
		Label: self.item.Label,
	}

	self.consumer.ConsumeSample(sample)
}
