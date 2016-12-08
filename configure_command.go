package main

import (
	"flag"
	"strings"

	"github.com/mitchellh/cli"
)

type configureCommand struct {
	ui cli.Ui
}

func (c *configureCommand) Run(args []string) int {
	var cfg mapFlags

	flags := flag.NewFlagSet("configure", flag.ContinueOnError)
	flags.Var(&cfg, "config", "config")
	if err := flags.Parse(args); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	config := make(map[string]string)
	for k, v := range cfg {
		config[k] = v
	}

	if err := writeConfig(config); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	c.ui.Output("mayoiga is successful configured")
	c.ui.Output("do not forget to add '.mayoiga' to .gitignore")
	return 0
}

func (c *configureCommand) Help() string {
	helpText := `
Usage: mayoiga configure -config [KEY=VALUE]

	Configure mayoiga.

Options:

  -config 'foo=bar'  Set a variable in the Mayoiga configuration.
`
	return strings.TrimSpace(helpText)
}

func (c *configureCommand) Synopsis() string {
	return "Configure mayoiga"
}
