package config

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Data struct {
	Label string `yaml:"label"`
	Color string `yaml:"color"`
	Script string `yaml:"script"`
}

func (d *Data) NextValue() (float64, error) {

	output, err := exec.Command("sh", "-c", d.Script).Output()
	if err != nil {
		log.Printf("%s", err)
	}

	trimmedOutput := strings.TrimSpace(string(output))
	floatValue, err := strconv.ParseFloat(trimmedOutput, 64)

	if err != nil {
		return 0, err
	}

	return floatValue, nil
}
