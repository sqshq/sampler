package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/sqshq/sampler/console"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme      *console.Theme    `yaml:"theme,omitempty"`
	Variables  map[string]string `yaml:"variables,omitempty"`
	RunCharts  []RunChartConfig  `yaml:"runcharts,omitempty"`
	BarCharts  []BarChartConfig  `yaml:"barcharts,omitempty"`
	Gauges     []GaugeConfig     `yaml:"gauges,omitempty"`
	SparkLines []SparkLineConfig `yaml:"sparklines,omitempty"`
	TextBoxes  []TextBoxConfig   `yaml:"textboxes,omitempty"`
	AsciiBoxes []AsciiBoxConfig  `yaml:"asciiboxes,omitempty"`
}

/** It's same to metadata/storage
	We should refact it
 **/
const (
	macOSDir   = "/Library/Application Support/Sampler"
	linuxDir   = "/.config/Sampler"
	windowsDir = "Sampler"
)

func getPlatformStoragePath(filename string) string {
	switch runtime.GOOS {
	case "darwin":
		home, _ := homedir.Dir()
		return filepath.Join(home, macOSDir, filename)
	case "windows":
		cache, _ := os.UserCacheDir()
		return filepath.Join(cache, windowsDir, filename)
	default:
		home, _ := homedir.Dir()
		return filepath.Join(home, linuxDir, filename)
	}
}

func LoadConfig() (*Config, Options) {

	var opt Options
	_, err := flags.Parse(&opt)

	if err != nil {
		console.Exit("")
	}

	if opt.Version == true {
		console.Exit(console.AppVersion)
	}

	if opt.ConfigFile == nil && opt.LicenseKey == nil {
		defaultConfigFile := getPlatformStoragePath("config.yml")

		if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
			console.Exit("Default config file is not existing! Please specify config file using --config flag. Example: sampler --config example.yml")
		} else {
			opt.ConfigFile = &defaultConfigFile
		}
	}

	if opt.LicenseKey != nil {
		return nil, opt
	}

	cfg := readFile(opt.ConfigFile)
	cfg.validate()
	cfg.setDefaults()

	return cfg, opt
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
	case TypeSparkLine:
		for i, component := range c.SparkLines {
			if component.Title == componentTitle {
				return &c.SparkLines[i].ComponentConfig
			}
		}
	case TypeAsciiBox:
		for i, component := range c.AsciiBoxes {
			if component.Title == componentTitle {
				return &c.AsciiBoxes[i].ComponentConfig
			}
		}
	case TypeTextBox:
		for i, component := range c.TextBoxes {
			if component.Title == componentTitle {
				return &c.TextBoxes[i].ComponentConfig
			}
		}
	}

	panic(fmt.Sprintf(
		"Failed to find component type %v with title %v", componentType, componentTitle))
}

func readFile(location *string) *Config {

	yamlFile, err := ioutil.ReadFile(*location)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", *location)
	}

	cfg := new(Config)
	err = yaml.Unmarshal(yamlFile, cfg)

	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	return cfg
}

func saveFile(config *Config, fileName *string) {
	file, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Failed to marshal config file: %v", err)
	}
	err = ioutil.WriteFile(*fileName, file, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to save config file: %v", err)
	}
}
