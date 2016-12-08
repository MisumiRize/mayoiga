package main

import (
	"flag"
	"strings"

	"github.com/mitchellh/cli"
)

type pullCommand struct {
	ui cli.Ui
}

func (c *pullCommand) Run(args []string) int {
	if err := assertConfig(); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	flags := flag.NewFlagSet("pull", flag.ContinueOnError)
	silent := flags.Bool("silent", false, "silent")
	if err := flags.Parse(args); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	config, err := readConfig()
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	buf, err := s3GetObject(config.S3Key)
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if !*silent {
		if err = compareWithEnvFile(c.ui, buf.String()); err != nil {
			c.ui.Error(err.Error())
			return 1
		}
	}

	if err = writeEnvFile(buf.Bytes()); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	mappings, err := fetchMappings()
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if err = writeMappingsFile(mappings); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if err = generateMappedEnvFiles(buf, mappings); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *pullCommand) Help() string {
	helpText := `
Usage: mayoiga pull

  Pulls env file from S3.

Options:

	-silent  Suppress diff output.
`
	return strings.TrimSpace(helpText)
}

func (c *pullCommand) Synopsis() string {
	return "Pull env file from S3"
}
