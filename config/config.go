package config

import (
	"github.com/sqshq/vcmd/data"
	"github.com/sqshq/vcmd/settings"
	. "github.com/sqshq/vcmd/widgets"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Theme     settings.Theme   `yaml:"theme"`
	RunCharts []RunChartConfig `yaml:"run-charts"`
}

type RunChartConfig struct {
	Title         string      `yaml:"title"`
	Items         []data.Item `yaml:"data"`
	Position      Position    `yaml:"position"`
	Size          Size        `yaml:"size"`
	RefreshRateMs int         `yaml:"refresh-rate-ms"`
	TimeScaleSec  int         `yaml:"time-scale-sec"`
}

func Load(location string) *Config {

	cfg := readFile(location)
	cfg.validate()
	cfg.setDefaultValues()
	cfg.setDefaultColors()
	cfg.setDefaultLayout()

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
