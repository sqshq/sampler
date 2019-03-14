package data

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/config"
	"os/exec"
	"strings"
)

type Item struct {
	Label  string
	Script string
	Color  *ui.Color
}

func NewItems(cfgs []config.Item) []Item {

	items := make([]Item, 0)

	for _, i := range cfgs {
		item := Item{Label: *i.Label, Script: i.Script, Color: i.Color}
		items = append(items, item)
	}

	return items
}

func (i *Item) nextValue() (value string, err error) {

	output, err := exec.Command("sh", "-c", i.Script).Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
