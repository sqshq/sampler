package config

import (
	. "github.com/sqshq/vcmd/layout"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	LineChartConfigs []LineChartConfig `yaml:"line-charts"`
}

type DataConfig struct {
	Script string `yaml:"script"`
	Label  string `yaml:"label"`
}

type LineChartConfig struct {
	Title         string       `yaml:"title"`
	DataConfig    []DataConfig `yaml:"data"`
	Position      Position     `yaml:"position"`
	Size          Size         `yaml:"size"`
	RefreshRateMs int          `yaml:"refresh-rate-ms"`
	TimeScaleSec  int          `yaml:"time-scale-sec"`
}

func Load(location string) *Config {

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
