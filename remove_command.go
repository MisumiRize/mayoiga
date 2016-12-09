package main

import (
	"bufio"
	"bytes"
	"flag"
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

	if len(args) < 1 || args[0][:1] == "-" {
		c.ui.Error(c.Help())
		return 1
	}

	flags := flag.NewFlagSet("remove", flag.ContinueOnError)
	prune := flags.Bool("prune", false, "prune")
	if err := flags.Parse(args[1:]); err != nil {
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

	if err = s3PutObject(config.S3Key, bytes.NewReader(buf.Bytes())); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if config.MappingS3Key == nil {
		return 0
	}

	mappings, err := fetchMappings()
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if *prune {
		mappings = removeVariable(mappings, args[0])
	}

	if err = writeMappingsFile(mappings); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if *prune {
		if err = putMappingsToS3(mappings); err != nil {
			c.ui.Error(err.Error())
			return 1
		}
	}

	if err = generateMappedEnvFiles(buf, mappings); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *removeCommand) Help() string {
	helpText := `
Usage: mayoiga remove <KEY>

	Remove env valiable and save to S3.

Options:

	-prune  Delete mapping associated with the variable.
`
	return strings.TrimSpace(helpText)
}

func (c *removeCommand) Synopsis() string {
	return "Remove env value and save to S3"
}
