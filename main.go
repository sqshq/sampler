package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/asset"
	"github.com/sqshq/sampler/client"
	"github.com/sqshq/sampler/component"
	"github.com/sqshq/sampler/component/asciibox"
	"github.com/sqshq/sampler/component/barchart"
	"github.com/sqshq/sampler/component/gauge"
	"github.com/sqshq/sampler/component/layout"
	"github.com/sqshq/sampler/component/runchart"
	"github.com/sqshq/sampler/component/sparkline"
	"github.com/sqshq/sampler/component/textbox"
	"github.com/sqshq/sampler/config"
	"github.com/sqshq/sampler/console"
	"github.com/sqshq/sampler/data"
	"github.com/sqshq/sampler/event"
	"github.com/sqshq/sampler/metadata"
	"runtime/debug"
	"time"
)

type Starter struct {
	player  *asset.AudioPlayer
	lout    *layout.Layout
	palette console.Palette
	opt     config.Options
	cfg     config.Config
}

func (s *Starter) startAll() []*data.Sampler {
	samplers := make([]*data.Sampler, 0)
	for _, c := range s.cfg.RunCharts {
		cpt := runchart.NewRunChart(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers))
	}
	for _, c := range s.cfg.SparkLines {
		cpt := sparkline.NewSparkLine(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers))
	}
	for _, c := range s.cfg.BarCharts {
		cpt := barchart.NewBarChart(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers))
	}
	for _, c := range s.cfg.Gauges {
		cpt := gauge.NewGauge(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Cur, c.Min, c.Max}, c.Triggers))
	}
	for _, c := range s.cfg.AsciiBoxes {
		cpt := asciibox.NewAsciiBox(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers))
	}
	for _, c := range s.cfg.TextBoxes {
		cpt := textbox.NewTextBox(c, s.palette)
		samplers = append(samplers, s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers))
	}
	return samplers
}

func (s *Starter) start(drawable ui.Drawable, consumer *data.Consumer, componentConfig config.ComponentConfig, itemsConfig []config.Item, triggersConfig []config.TriggerConfig) *data.Sampler {
	cpt := component.NewComponent(drawable, consumer, componentConfig)
	triggers := data.NewTriggers(triggersConfig, consumer, s.opt, s.player)
	items := data.NewItems(itemsConfig, *componentConfig.RateMs)
	s.lout.AddComponent(cpt)
	time.Sleep(10 * time.Millisecond) // desync coroutines
	return data.NewSampler(consumer, items, triggers, s.opt, s.cfg.Variables, *componentConfig.RateMs)
}

func main() {

	cfg, opt := config.LoadConfig()
	bc := client.NewBackendClient()

	statistics := metadata.GetStatistics(cfg)
	license := metadata.GetLicense()

	if opt.LicenseKey != nil {
		registerLicense(statistics, opt, bc)
	}

	console.Init()
	defer console.Close()

	player := asset.NewAudioPlayer()
	if player != nil {
		defer player.Close()
	}

	defer handleCrash(statistics, opt, bc)
	defer updateStatistics(cfg, time.Now())

	palette := console.GetPalette(*cfg.Theme)
	lout := layout.NewLayout(component.NewStatusBar(*opt.ConfigFile, palette, license),
		component.NewMenu(palette), component.NewIntro(palette), component.NewNagWindow(palette))

	if statistics.LaunchCount == 0 {
		if !opt.DisableTelemetry {
			go bc.ReportInstallation(statistics)
		}
		lout.StartWithIntro()
	} else if statistics.LaunchCount%20 == 0 { // once in a while
		if license == nil || !license.Valid {
			lout.StartWithNagWindow()
		} else {
			go verifyLicense(license, bc)
		}
		if !opt.DisableTelemetry {
			go bc.ReportStatistics(statistics)
		}
	}

	starter := &Starter{player, lout, palette, opt, *cfg}
	samplers := starter.startAll()

	handler := event.NewHandler(samplers, opt, lout)
	handler.HandleEvents()
}

func handleCrash(statistics *metadata.Statistics, opt config.Options, bc *client.BackendClient) {
	if rec := recover(); rec != nil {
		err := rec.(error)
		if !opt.DisableTelemetry {
			bc.ReportCrash(fmt.Sprintf("%s\n%s", err.Error(), string(debug.Stack())), statistics)
		}
		panic(err)
	}
}

func updateStatistics(cfg *config.Config, startTime time.Time) {
	metadata.PersistStatistics(cfg, time.Since(startTime))
}

func registerLicense(statistics *metadata.Statistics, opt config.Options, bc *client.BackendClient) {
	lc, err := bc.RegisterLicenseKey(*opt.LicenseKey, statistics)
	if err != nil {
		console.Exit("License registration failed: " + err.Error())
	} else {
		metadata.SaveLicense(*lc)
		console.Exit("License successfully verified, Sampler can be restarted without --license flag now. Thank you.")
	}
}

func verifyLicense(license *metadata.License, bc *client.BackendClient) {
	verifiedLicense, _ := bc.VerifyLicenseKey(*license.Key)
	if verifiedLicense != nil {
		metadata.SaveLicense(*verifiedLicense)
	}
}
