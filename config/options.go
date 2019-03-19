package config

type Options struct {
	ConfigFile string   `short:"c" long:"config" required:"true" description:"path to YAML config file"`
	Variables  []string `short:"v" long:"variable" required:"false" description:"variable for script ${var-name} placeholder" long-description:"one or more variables can be specified as flags, in order to replace repeated patterns in the scripts, which can be replaced with {$variable-name} placeholder" `
	Examples   []string `short:"e" long:"example" required:"false" choice:"runchart" choice:"barchart" choice:"asciibox" choice:"textbox" choice:"gauge" choice:"sparkline" description:"add an example component to the specified config file" long-description:"one or more example component types can be added to the specified config file, in order to jump-start the configuration"`
}
