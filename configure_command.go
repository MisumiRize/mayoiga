package main

import (
	"flag"
	"strings"

	"fmt"

	"github.com/mitchellh/cli"
)

type configureCommand struct {
	ui cli.Ui
}

type stringFlags []string

func (f *stringFlags) String() string {
	return strings.Join(*f, "\n")
}

func (f *stringFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (c *configureCommand) Run(args []string) int {
	var cfg stringFlags

	flags := flag.NewFlagSet("configure", flag.ContinueOnError)
	flags.Var(&cfg, "config", "config")
	if err := flags.Parse(args); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	config := make(map[string]string)
	for _, f := range cfg {
		v := strings.Split(f, "=")
		if len(v) < 2 {
			c.ui.Warn(fmt.Sprintf("invalid config %s, config should be key=value", f))
			continue
		}
		config[v[0]] = v[1]
	}

	if err := writeConfig(config); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	c.ui.Output("mayoiga is successful configured")
	c.ui.Output("do not forget to add '.mayoiga.json' to .gitignore")
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
