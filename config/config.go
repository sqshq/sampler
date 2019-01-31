package config

import (
	"github.com/sqshq/vcmd/data"
	. "github.com/sqshq/vcmd/layout"
	"github.com/sqshq/vcmd/settings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Theme     settings.Theme `yaml:"theme"`
	RunCharts []RunChart     `yaml:"run-charts"`
}

type RunChart struct {
	Title         string      `yaml:"title"`
	Items         []data.Item `yaml:"data"`
	Position      Position    `yaml:"position"`
	Size          Size        `yaml:"size"`
	RefreshRateMs int         `yaml:"refresh-rate-ms"`
	TimeScaleSec  int         `yaml:"time-scale-sec"`
}

func Load(location string) *Config {

	cfg := readFile(location)
	validate(cfg)
	setColors(cfg)

	return cfg
}

func readFile(location string) *Config {

	yamlFile, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatalf("Can't read config file: %s", location)
	}

	cfg := new(Config)
	err = yaml.Unmarshal(yamlFile, cfg)

	if err != nil {
		log.Fatalf("Can't read config file: %v", err)
	}

	return cfg
}

/*
 TODO validation
 - title uniquness and mandatory within a single type of widget
 - label uniqueness and mandatory (if > 1 data bullets)
*/
func validate(config *Config) {

}

func setColors(config *Config) {

	palette := settings.GetPalette(config.Theme)

	for i, chart := range config.RunCharts {
		for j, item := range chart.Items {
			if item.Color == 0 {
				item.Color = palette.Colors[i+j]
				chart.Items[j] = item
			}
		}
	}
}
