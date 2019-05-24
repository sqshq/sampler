package metadata

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"gopkg.in/yaml.v2"
	"log"
	"runtime"
)

type Statistics struct {
	Version         string
	OS              string
	WindowWidth     int            `yaml:"ww"`
	WindowHeight    int            `yaml:"wh"`
	LaunchCount     int            `yaml:"lc"`
	ComponentsCount map[string]int `yaml:"cc"`
}

const statisticsFileName = "statistics.yml"

func PersistStatistics(config *config.Config) *Statistics {

	statistics := new(Statistics)
	w, h := ui.TerminalDimensions()

	if fileExists(statisticsFileName) {
		file := readStorageFile(getPlatformStoragePath(statisticsFileName))
		err := yaml.Unmarshal(file, statistics)

		if err != nil {
			log.Fatalf("Failed to read statistics file: %v", err)
		}

		if config != nil {
			statistics.ComponentsCount = countComponentsPerType(config)
		}

		statistics.WindowWidth = w
		statistics.WindowWidth = h
		statistics.LaunchCount += 1

	} else {
		statistics = &Statistics{
			Version:         console.AppVersion,
			OS:              runtime.GOOS,
			LaunchCount:     1,
			WindowWidth:     w,
			WindowHeight:    h,
			ComponentsCount: countComponentsPerType(config),
		}
	}

	file, err := yaml.Marshal(statistics)
	if err != nil {
		log.Fatalf("Failed to marshal statistics file: %v", err)
	}

	saveStorageFile(file, statisticsFileName)

	return statistics
}

func countComponentsPerType(config *config.Config) map[string]int {

	m := make(map[string]int)

	if config == nil {
		return m
	}

	m["runcharts"] = len(config.RunCharts)
	m["sparkLines"] = len(config.SparkLines)
	m["barcharts"] = len(config.BarCharts)
	m["gauges"] = len(config.Gauges)
	m["asciiboxes"] = len(config.AsciiBoxes)
	m["textboxes"] = len(config.TextBoxes)

	return m
}
