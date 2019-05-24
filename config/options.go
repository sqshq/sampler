package config

type Options struct {
	ConfigFile  *string  `short:"c" long:"config" required:"false" description:"set path to YAML config file"`
	License     *string  `short:"l" long:"license" required:"false" description:"provide license key. visit www.sampler.dev for details"`
	Environment []string `short:"e" long:"env" required:"false" description:"specify name=value variable to use in script placeholder as $name. This flag takes precedence over the same name variables, specified in config yml"`
	Version     bool     `short:"v" long:"version" required:"false" description:"print version"`
}
