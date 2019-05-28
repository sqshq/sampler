package main

import (
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
	"time"
)

type Starter struct {
	player  *asset.AudioPlayer
	lout    *layout.Layout
	palette console.Palette
	opt     config.Options
	cfg     config.Config
}

func (s *Starter) startAll() {
	for _, c := range s.cfg.RunCharts {
		cpt := runchart.NewRunChart(c, s.palette)
		s.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}
	for _, c := range s.cfg.SparkLines {
		cpt := sparkline.NewSparkLine(c, s.palette)
		s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers)
	}
	for _, c := range s.cfg.BarCharts {
		cpt := barchart.NewBarChart(c, s.palette)
		s.start(cpt, cpt.Consumer, c.ComponentConfig, c.Items, c.Triggers)
	}
	for _, c := range s.cfg.Gauges {
		cpt := gauge.NewGauge(c, s.palette)
		s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Cur, c.Min, c.Max}, c.Triggers)
	}
	for _, c := range s.cfg.AsciiBoxes {
		cpt := asciibox.NewAsciiBox(c, s.palette)
		s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers)
	}
	for _, c := range s.cfg.TextBoxes {
		cpt := textbox.NewTextBox(c, s.palette)
		s.start(cpt, cpt.Consumer, c.ComponentConfig, []config.Item{c.Item}, c.Triggers)
	}
}

func (s *Starter) start(drawable ui.Drawable, consumer *data.Consumer, componentConfig config.ComponentConfig, itemsConfig []config.Item, triggersConfig []config.TriggerConfig) {
	cpt := component.NewComponent(drawable, consumer, componentConfig)
	triggers := data.NewTriggers(triggersConfig, consumer, s.opt, s.player)
	items := data.NewItems(itemsConfig, *componentConfig.RateMs)
	data.NewSampler(consumer, items, triggers, s.opt, s.cfg.Variables, *componentConfig.RateMs)
	s.lout.AddComponent(cpt)
	time.Sleep(10 * time.Millisecond) // desync coroutines
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
	defer player.Close()

	palette := console.GetPalette(*cfg.Theme)
	lout := layout.NewLayout(
		component.NewStatusLine(*opt.ConfigFile, palette, license), component.NewMenu(palette), component.NewIntro(palette))

	if statistics.LaunchCount == 0 {
		if !opt.DisableTelemetry {
			go bc.ReportInstallation(statistics)
		}
		lout.RunIntro()
	} else /* with random */ {
		// TODO if license == nil lout.showNagWindow() with timeout and OK button
		// TODO if license != nil, verify license
		// TODO report statistics
	}

	metadata.PersistStatistics(cfg)
	starter := &Starter{player, lout, palette, opt, *cfg}
	starter.startAll()

	handler := event.NewHandler(lout, opt)
	handler.HandleEvents()
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
