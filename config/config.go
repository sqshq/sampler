package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/sqshq/sampler/console"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Theme      *console.Theme    `yaml:"theme,omitempty"`
	RunCharts  []RunChartConfig  `yaml:"runcharts,omitempty"`
	BarCharts  []BarChartConfig  `yaml:"barcharts,omitempty"`
	Gauges     []GaugeConfig     `yaml:"gauges,omitempty"`
	AsciiBoxes []AsciiBoxConfig  `yaml:"asciiboxes,omitempty"`
	SparkLines []SparkLineConfig `yaml:"sparklines,omitempty"`
}

func Load() (Config, Options) {

	var opt Options
	_, err := flags.Parse(&opt)

	if err != nil {
		os.Exit(0)
	}

	cfg := readFile(opt.ConfigFile)
	cfg.validate()
	cfg.setDefaults()

	return *cfg, opt
}

func Update(settings []ComponentSettings, options Options) {
	cfg := readFile(options.ConfigFile)
	for _, s := range settings {
		componentConfig := cfg.findComponent(s.Type, s.Title)
		componentConfig.Position = getPosition(s.Location, s.Size)
	}
	saveFile(cfg, options.ConfigFile)
}

func (c *Config) findComponent(componentType ComponentType, componentTitle string) *ComponentConfig {

	switch componentType {
	case TypeRunChart:
		for i, component := range c.RunCharts {
			if component.Title == componentTitle {
				return &c.RunCharts[i].ComponentConfig
			}
		}
	case TypeBarChart:
		for i, component := range c.BarCharts {
			if component.Title == componentTitle {
				return &c.BarCharts[i].ComponentConfig
			}
		}
	case TypeGauge:
		for i, component := range c.Gauges {
			if component.Title == componentTitle {
				return &c.Gauges[i].ComponentConfig
			}
		}
	case TypeAsciiBox:
		for i, component := range c.AsciiBoxes {
			if component.Title == componentTitle {
				return &c.AsciiBoxes[i].ComponentConfig
			}
		}
	case TypeSparkLine:
		for i, component := range c.SparkLines {
			if component.Title == componentTitle {
				return &c.SparkLines[i].ComponentConfig
			}
		}
	}

	panic(fmt.Sprintf(
		"Can't find component type %v with title %v", componentType, componentTitle))
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

func saveFile(config *Config, fileName string) {
	file, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Can't marshal config file: %v", err)
	}

	err = ioutil.WriteFile(fileName, file, 0644)
	if err != nil {
		log.Fatalf("Can't save config file: %v", err)
	}
}
