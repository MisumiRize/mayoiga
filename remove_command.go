package main

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/mitchellh/cli"
)

type removeCommand struct {
	ui cli.Ui
}

func (c *removeCommand) Run(args []string) int {
	if err := assertConfig(); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if len(args) < 1 {
		c.ui.Error(c.Help())
		return 1
	}

	buf, err := s3GetObject()
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	buf, err = deleteEnv(bufio.NewScanner(buf), args[0])
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

	if err = s3PutObject(bytes.NewReader(buf.Bytes())); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *removeCommand) Help() string {
	helpText := `
Usage: mayoiga remove <KEY>

	Remove env value and save to S3.
`
	return strings.TrimSpace(helpText)
}

func (c *removeCommand) Synopsis() string {
	return "Remove env value and save to S3"
}
