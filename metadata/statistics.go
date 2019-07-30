package metadata

import (
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"gopkg.in/yaml.v3"
	"log"
	"runtime"
	"time"
)

// Statistics represents anonymous usage data, which we collect for analyses and improvements
// User can disable it, along with crash reports, using --disable-telemetry flag
type Statistics struct {
	Version         string
	OS              string
	WindowWidth     int            `yaml:"ww"`
	WindowHeight    int            `yaml:"wh"`
	LaunchCount     int            `yaml:"lc"`
	UsageTime       int            `yaml:"ut"`
	ComponentsCount map[string]int `yaml:"cc"`
}

const statisticsFileName = "statistics.yml"

func PersistStatistics(config *config.Config, uptime time.Duration) *Statistics {

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

		statistics.Version = console.AppVersion
		statistics.WindowWidth = w
		statistics.WindowHeight = h
		statistics.LaunchCount += 1
		statistics.UsageTime += int(uptime.Seconds())

	} else {
		statistics = &Statistics{
			Version:         console.AppVersion,
			OS:              runtime.GOOS,
			WindowWidth:     w,
			WindowHeight:    h,
			LaunchCount:     1,
			UsageTime:       0,
			ComponentsCount: countComponentsPerType(config),
		}
		initStorage()
	}

	file, err := yaml.Marshal(statistics)
	if err != nil {
		log.Fatalf("Failed to marshal statistics file: %v", err)
	}

	saveStorageFile(file, statisticsFileName)

	return statistics
}

func GetStatistics(cfg *config.Config) *Statistics {

	if !fileExists(statisticsFileName) {
		return &Statistics{
			Version:         console.AppVersion,
			OS:              runtime.GOOS,
			LaunchCount:     0,
			WindowWidth:     0,
			WindowHeight:    0,
			ComponentsCount: countComponentsPerType(cfg),
		}
	}

	file := readStorageFile(getPlatformStoragePath(statisticsFileName))
	license := new(Statistics)

	err := yaml.Unmarshal(file, license)
	if err != nil {
		log.Fatalf("Failed to read statistics file: %v", err)
	}

	return license
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
