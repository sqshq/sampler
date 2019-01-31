package data

import (
	"time"
)

type Poller struct {
	consumer Consumer
	item     Item
	pause    bool
}

func NewPoller(consumer Consumer, item Item, rateMs int) Poller {

	ticker := time.NewTicker(time.Duration(rateMs * int(time.Millisecond)))
	poller := Poller{consumer, item, false}

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

func (self *Poller) TogglePause() {
	self.pause = !self.pause
}

func (self *Poller) poll() {

	if self.pause {
		return
	}

	value, err := self.item.nextValue()

	if err != nil {
		self.consumer.ConsumeError(self.item, err)
	}

	self.consumer.ConsumeValue(self.item, value)
}
