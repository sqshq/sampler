package data

import (
	"fmt"
	"github.com/sqshq/sampler/asset"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"os"
	"os/exec"
	"regexp"
)

const (
	TrueIndicator = "1"
	InitialValue  = "0"
)

type Trigger struct {
	title         string
	condition     string
	actions       *Actions
	consumer      *Consumer
	valuesByLabel map[string]Values
	options       config.Options
	player        *asset.AudioPlayer
	digitsRegexp  *regexp.Regexp
}

type Actions struct {
	terminalBell bool
	sound        bool
	visual       bool
	script       *string
}

type Values struct {
	current  string
	previous string
}

func NewTriggers(cfgs []config.TriggerConfig, consumer *Consumer, options config.Options, player *asset.AudioPlayer) []*Trigger {

	triggers := make([]*Trigger, 0)

	for _, cfg := range cfgs {
		triggers = append(triggers, NewTrigger(cfg, consumer, options, player))
	}

	return triggers
}

func NewTrigger(config config.TriggerConfig, consumer *Consumer, options config.Options, player *asset.AudioPlayer) *Trigger {
	return &Trigger{
		title:         config.Title,
		condition:     config.Condition,
		consumer:      consumer,
		valuesByLabel: make(map[string]Values),
		options:       options,
		player:        player,
		digitsRegexp:  regexp.MustCompile("[^0-9]+"),
		actions: &Actions{
			terminalBell: *config.Actions.TerminalBell,
			sound:        *config.Actions.Sound,
			visual:       *config.Actions.Visual,
			script:       config.Actions.Script,
		},
	}
}

func (t *Trigger) Execute(sample *Sample) {
	if t.evaluate(sample) {

		if t.actions.terminalBell {
			fmt.Print(console.BellCharacter)
		}

		if t.actions.sound && t.player != nil {
			t.player.Beep()
		}

		if t.actions.visual {
			t.consumer.AlertChannel <- &Alert{
				Title:       t.title,
				Text:        fmt.Sprintf("%s: %v", sample.Label, sample.Value),
				Color:       sample.Color,
				Recoverable: false,
			}
		}

		if t.actions.script != nil {
			_, _ = t.runScript(*t.actions.script, sample.Label, t.valuesByLabel[sample.Label])
		}
	}
}

func (t *Trigger) evaluate(sample *Sample) bool {

	if values, ok := t.valuesByLabel[sample.Label]; ok {
		values.previous = values.current
		values.current = sample.Value
		t.valuesByLabel[sample.Label] = values
	} else {
		t.valuesByLabel[sample.Label] = Values{previous: InitialValue, current: sample.Value}
	}

	output, err := t.runScript(t.condition, sample.Label, t.valuesByLabel[sample.Label])

	if err != nil {
		t.consumer.AlertChannel <- &Alert{
			Title:       "Trigger condition failure",
			Text:        getErrorMessage(err),
			Color:       sample.Color,
			Recoverable: true,
		}
	}

	return t.digitsRegexp.ReplaceAllString(string(output), "") == TrueIndicator
}

func (t *Trigger) runScript(script, label string, data Values) ([]byte, error) {

	cmd := exec.Command("sh", "-c", script)
	cmd.Env = os.Environ()

	for _, variable := range t.options.Environment {
		cmd.Env = append(cmd.Env, variable)
	}

	cmd.Env = append(cmd.Env,
		fmt.Sprintf("prev=%v", data.previous),
		fmt.Sprintf("cur=%v", data.current),
		fmt.Sprintf("label=%v", label))

	return cmd.Output()
}
