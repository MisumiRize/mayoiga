package main

import (
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

	buf, err := s3GetObject()
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if err = compareWithEnvFile(c.ui, buf.String()); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if err = writeEnvFile(buf.Bytes()); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *pullCommand) Help() string {
	helpText := `
Usage: mayoiga pull

  Pulls env file from S3.
`
	return strings.TrimSpace(helpText)
}

func (c *pullCommand) Synopsis() string {
	return "Pull env file from S3"
}
