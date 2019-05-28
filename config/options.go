package config

type Options struct {
	ConfigFile       *string  `short:"c" long:"config" required:"true" description:"Path to YAML config file"`
	LicenseKey       *string  `short:"l" long:"license" description:"License key. Visit www.sampler.dev for details"`
	Environment      []string `short:"e" long:"env" description:"Specify name=value variable to use in script placeholder as $name. This flag takes precedence over the same name variables, specified in config yml"`
	Version          bool     `short:"v" long:"version" description:"Print version"`
	DisableTelemetry bool     `long:"disable-telemetry" description:"Disable anonymous usage statistics and errors to be sent to Sampler online service for analyses"`
}
