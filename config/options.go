package config

type Options struct {
	ConfigFile string   `short:"c" long:"config" required:"true" description:"path to YAML config file"`
	Variables  []string `short:"v" long:"variable" required:"false" description:"specify name=value variable to use in script placeholder as $name. This flag takes precedence over the same name variables, specified in config yml" long-description:"one or more variables can be specified as flags, in order to replace repeated patterns in the scripts, which can be replaced with {$variable-name} placeholder" `
	License    []string `short:"l" long:"license" required:"false" description:"provide license key. see www.sampler.dev for details"`
}
