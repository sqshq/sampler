package config

// Options with cli flags
type Options struct {
	ConfigFile  *string  `short:"c" long:"config" description:"Path to YAML config file"`
	HideStatusBar bool `long:"hide-statusbar" description:"Hide the bottom status bar (version, select, pause, quit)"`
    Environment []string `short:"e" long:"env" description:"Specify name=value variable to use in script placeholder as $name. This flag takes precedence over the same name variables, specified in config yml"`
	Version     bool     `short:"v" long:"version" description:"Print version"`
}
