package data

import (
	"os/exec"
	"strings"
	"time"
)

type Poller struct {
	consumer Consumer
	script   string
	label    string
	pause    bool
}

func NewPoller(consumer Consumer, script string, label string, rateMs int) Poller {

	ticker := time.NewTicker(time.Duration(rateMs * int(time.Millisecond)))
	poller := Poller{consumer, script, label, false}

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

	output, err := exec.Command("sh", "-c", self.script).Output()

	if err != nil {
		self.consumer.ConsumeError(err)
	}

	value := strings.TrimSpace(string(output))
	self.consumer.ConsumeValue(value, self.label)
}
