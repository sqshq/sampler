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
}

func NewSampler(consumer *Consumer, items []*Item, triggers []*Trigger, options config.Options, fileVariables map[string]string, rateMs int) Sampler {

	ticker := time.NewTicker(time.Duration(rateMs * int(time.Millisecond)))

	sampler := Sampler{
		consumer,
		items,
		triggers,
		make(chan *Sample),
		mergeVariables(fileVariables, options.Environment),
	}

	go func() {
		for ; true; <-ticker.C {
			for _, item := range sampler.items {
				sampler.sample(item, options)
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

	val, err := item.nextValue(s.variables)

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

// option variables takes precedence over the file variables with the same name
func mergeVariables(fileVariables map[string]string, optionsVariables []string) []string {

	result := optionsVariables

	for key, value := range fileVariables {
		result = append([]string{fmt.Sprintf("%s=%s", key, value)}, result...)
	}

	return result
}
