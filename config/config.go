package config

import (
	"fmt"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"github.com/sqshq/sampler/widgets/asciibox"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Theme      *console.Theme   `yaml:"theme,omitempty"`
	RunCharts  []RunChartConfig `yaml:"runcharts,omitempty"`
	BarCharts  []BarChartConfig `yaml:"barcharts,omitempty"`
	AsciiBoxes []AsciiBoxConfig `yaml:"asciiboxes,omitempty"`
}

type ComponentConfig struct {
	Title         string   `yaml:"title"`
	RefreshRateMs *int     `yaml:"refresh-rate-ms,omitempty"`
	Position      Position `yaml:"position"`
	Size          Size     `yaml:"size"`
}

type RunChartConfig struct {
	ComponentConfig `yaml:",inline"`
	Legend          *LegendConfig `yaml:"legend,omitempty"`
	Scale           *int          `yaml:"scale,omitempty"`
	Items           []data.Item   `yaml:"items"`
}

type BarChartConfig struct {
	ComponentConfig `yaml:",inline"`
	Scale           *int        `yaml:"scale,omitempty"`
	Items           []data.Item `yaml:"items"`
}

type AsciiBoxConfig struct {
	ComponentConfig `yaml:",inline"`
	data.Item       `yaml:",inline"`
	Font            *asciibox.AsciiFont `yaml:"font,omitempty"`
}

type LegendConfig struct {
	Enabled bool `yaml:"enabled"`
	Details bool `yaml:"details"`
}

type Position struct {
	X int `yaml:"w"`
	Y int `yaml:"h"`
}

type Size struct {
	X int `yaml:"w"`
	Y int `yaml:"h"`
}

type ComponentType rune

const (
	TypeRunChart ComponentType = 0
	TypeBarChart ComponentType = 1
	TypeTextBox  ComponentType = 2
	TypeAsciiBox ComponentType = 3
)

type ComponentSettings struct {
	Type     ComponentType
	Title    string
	Size     Size
	Position Position
}

func Load() *Config {

	if len(os.Args) < 2 {
		println("Please specify config file location. See www.github.com/sqshq/sampler for the reference")
		os.Exit(0)
	}

	cfg := readFile(os.Args[1])
	cfg.validate()
	cfg.setDefaults()

	return cfg
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

	err = ioutil.WriteFile("config.yml", file, 0644)
	if err != nil {
		log.Fatalf("Can't save config file: %v", err)
	}
}
