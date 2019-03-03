package config

import (
	"fmt"
	"github.com/sqshq/sampler/console"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Theme      *console.Theme   `yaml:"theme,omitempty"`
	RunCharts  []RunChartConfig `yaml:"runcharts,omitempty"`
	BarCharts  []BarChartConfig `yaml:"barcharts,omitempty"`
	Gauges     []GaugeConfig    `yaml:"gauges,omitempty"`
	AsciiBoxes []AsciiBoxConfig `yaml:"asciiboxes,omitempty"`
}

type Flags struct {
	ConfigFileName string
	Variables      map[string]string
}

func Load() (Config, Flags) {

	//if len(os.Args) < 2 {
	//	println("Please specify config file location. See www.github.com/sqshq/sampler for the reference")
	//	os.Exit(0)
	//}

	cfg := readFile("config.yml")
	cfg.validate()
	cfg.setDefaults()

	flg := Flags{ConfigFileName: "config.yml"}

	return *cfg, flg
}

func Update(settings []ComponentSettings) {
	cfg := readFile(os.Args[1])
	for _, s := range settings {
		componentConfig := cfg.findComponent(s.Type, s.Title)
		componentConfig.Size = s.Size
		componentConfig.Position = s.Position
	}
	saveFile(cfg)
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

func saveFile(config *Config) {
	file, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Can't marshal config file: %v", err)
	}

	err = ioutil.WriteFile(os.Args[1], file, 0644)
	if err != nil {
		log.Fatalf("Can't save config file: %v", err)
	}
}
