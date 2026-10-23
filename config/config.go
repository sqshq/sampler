package config

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/sqshq/sampler/console"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

// DefaultConfigCandidates defines fallback file names searched when no config flag is passed
var DefaultConfigCandidates = []string{
	"example.yml",
	"config.yml",
	"sampler.yml",
	".sampler.yml",
}

// ConfigSearchDirectories returns directories to search for configuration files
func ConfigSearchDirectories() []string {
	dirs := []string{"."}
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		if execDir != "." && execDir != "" {
			dirs = append(dirs, execDir)
		}
	}
	dirs = append(dirs, "/app", "/sampler", "/etc/sampler")
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, ".config", "sampler"), home)
	}
	return dirs
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// ResolveConfigPath resolves a configuration file path, checking existing paths,
// standard container/system directories, and default filenames when omitted.
func ResolveConfigPath(path *string) (string, error) {
	return ResolveConfigPathWithDirs(path, ConfigSearchDirectories())
}

// ResolveConfigPathWithDirs resolves a configuration file path using specified search directories.
func ResolveConfigPathWithDirs(path *string, dirs []string) (string, error) {
	if path != nil && *path != "" {
		clean := filepath.Clean(*path)
		if fileExists(clean) {
			return clean, nil
		}
		if !filepath.IsAbs(clean) {
			for _, dir := range dirs {
				candidate := filepath.Join(dir, clean)
				if fileExists(candidate) {
					return candidate, nil
				}
			}
		}
		return "", fmt.Errorf("config file not found: %s", *path)
	}

	for _, dir := range dirs {
		for _, name := range DefaultConfigCandidates {
			candidate := filepath.Join(dir, name)
			if fileExists(candidate) {
				return candidate, nil
			}
		}
	}

	return "", errors.New("Please specify config file using --config flag. Example: sampler --config example.yml")
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

	resolvedPath, err := ResolveConfigPath(opt.ConfigFile)
	if err != nil {
		console.Exit(err.Error())
	}
	opt.ConfigFile = &resolvedPath

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

// LoadConfigFromPath reads, parses, validates, and applies defaults to a config file.
func LoadConfigFromPath(path string) (*Config, error) {
	cfg, err := readConfigFile(path)
	if err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	cfg.setDefaults()
	return cfg, nil
}

// ParseConfig unmarshals, validates, and applies defaults to YAML configuration bytes.
func ParseConfig(data []byte) (*Config, error) {
	cfg := new(Config)
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	cfg.setDefaults()
	return cfg, nil
}

func readConfigFile(location string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := new(Config)
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

func readFile(location *string) *Config {
	cfg, err := readConfigFile(*location)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", *location)
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
