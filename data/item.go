package data

import (
	. "github.com/sqshq/termui"
	"os/exec"
	"strings"
)

type Item struct {
	Script string `yaml:"script"`
	Label  string `yaml:"label"`
	Color  Color  `yaml:"color"`
}

func (self *Item) nextValue() (value string, err error) {

	output, err := exec.Command("sh", "-c", self.Script).Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
