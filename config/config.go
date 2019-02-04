package config

import (
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	. "github.com/sqshq/sampler/widgets"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Theme     console.Theme    `yaml:"theme"`
	RunCharts []RunChartConfig `yaml:"runcharts"`
}

type RunChartConfig struct {
	Title         string      `yaml:"title"`
	Items         []data.Item `yaml:"items"`
	Position      Position    `yaml:"position"`
	Size          Size        `yaml:"size"`
	RefreshRateMs int         `yaml:"refresh-rate-ms"`
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
