package data

import (
	"fmt"
	"github.com/sqshq/sampler/config"
	"time"
)

type Sampler struct {
	consumer        *Consumer
	items           []*Item
	triggers        []*Trigger
	triggersChannel chan *Sample
	variables       []string
	pause           bool
}

func NewSampler(consumer *Consumer, items []*Item, triggers []*Trigger, options config.Options, fileVariables map[string]string, rateMs int) *Sampler {

	ticker := time.NewTicker(time.Duration(uint32(rateMs) * uint32(time.Millisecond)))

	sampler := &Sampler{
		consumer,
		items,
		triggers,
		make(chan *Sample),
		mergeVariables(fileVariables, options.Environment),
		false,
	}

	go func() {
		for ; true; <-ticker.C {
			for _, item := range sampler.items {
				if !sampler.pause {
					go sampler.sample(item, options)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case sample := <-sampler.triggersChannel:
				for _, t := range sampler.triggers {
					if !sampler.pause {
						t.Execute(sample)
					}
				}
			}
		}
	}()

	return sampler
}

func (s *Sampler) sample(item *Item, options config.Options) {

	val, err := item.nextValue(s.variables)

	if len(val) > 0 {
		sample := &Sample{Label: item.label, Value: val, Color: item.color}
		s.consumer.SampleChannel <- sample
		s.triggersChannel <- sample
	} else if err != nil {
		s.consumer.AlertChannel <- &Alert{
			Title:       "Sampling failure",
			Text:        getErrorMessage(err),
			Color:       item.color,
			Recoverable: true,
		}
	}
}

// option variables takes precedence over the file variables with the same name
func mergeVariables(fileVariables map[string]string, optionsVariables []string) []string {

	result := optionsVariables

	for key, value := range fileVariables {
		result = append([]string{fmt.Sprintf("%s=%s", key, value)}, result...)
	}

	return result
}

func (s *Sampler) Pause(pause bool) {
	s.pause = pause
}
