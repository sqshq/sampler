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
	actions       Actions
	consumer      Consumer
	valuesByLabel map[string]Values
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

func NewTriggers(cfgs []config.TriggerConfig, consumer Consumer, player *asset.AudioPlayer) []Trigger {

	triggers := make([]Trigger, 0)

	for _, cfg := range cfgs {
		triggers = append(triggers, NewTrigger(cfg, consumer, player))
	}

	return triggers
}

func NewTrigger(config config.TriggerConfig, consumer Consumer, player *asset.AudioPlayer) Trigger {
	return Trigger{
		title:         config.Title,
		condition:     config.Condition,
		consumer:      consumer,
		valuesByLabel: make(map[string]Values),
		player:        player,
		digitsRegexp:  regexp.MustCompile("[^0-9]+"),
		actions: Actions{
			terminalBell: *config.Actions.TerminalBell,
			sound:        *config.Actions.Sound,
			visual:       *config.Actions.Visual,
			script:       config.Actions.Script,
		},
	}
}

func (t *Trigger) Execute(sample Sample) {
	if t.evaluate(sample) {

		if t.actions.terminalBell {
			fmt.Print(console.BellCharacter)
		}

		if t.actions.sound {
			t.player.Beep()
		}

		if t.actions.visual {
			t.consumer.AlertChannel <- Alert{
				Title: t.title, Text: fmt.Sprintf("%s value: %v", sample.Label, sample.Value),
			}
		}

		if t.actions.script != nil {
			_, _ = runScript(*t.actions.script, sample.Label, t.valuesByLabel[sample.Label])
		}
	}
}

func (t *Trigger) evaluate(sample Sample) bool {

	if values, ok := t.valuesByLabel[sample.Label]; ok {
		values.previous = values.current
		values.current = sample.Value
		t.valuesByLabel[sample.Label] = values
	} else {
		t.valuesByLabel[sample.Label] = Values{previous: InitialValue, current: sample.Value}
	}

	output, err := runScript(t.condition, sample.Label, t.valuesByLabel[sample.Label])

	if err != nil {
		//t.consumer.AlertChannel <- Alert{Title: "TRIGGER CONDITION FAILURE", Text: err.Error()}
	}

	return t.digitsRegexp.ReplaceAllString(string(output), "") == TrueIndicator
}

func runScript(script, label string, data Values) ([]byte, error) {

	cmd := exec.Command("sh", "-c", script)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("prev=%v", data.previous),
		fmt.Sprintf("cur=%v", data.current),
		fmt.Sprintf("label=%v", label))

	return cmd.Output()
}
