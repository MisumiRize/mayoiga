package main

import (
	"bytes"
	"flag"

	"github.com/mitchellh/cli"

	"bufio"

	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

type addCommand struct {
	ui cli.Ui
}

func (c *addCommand) Run(args []string) int {
	if err := assertConfig(); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	encrypt := flags.Bool("encrypt", false, "encrypt")
	key := flags.String("key", "", "key")
	value := flags.String("value", "", "value")
	if err := flags.Parse(args); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if len(*key) == 0 {
		c.ui.Error("key should be valid string")
		return 1
	}

	if *encrypt {
		v, err := kmsEncrypt(*value)
		if err != nil {
			c.ui.Error(err.Error())
			return 1
		}

		value = v
	}

	buf, err := s3GetObject()
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == "NoSuchKey" {
			buf = new(bytes.Buffer)
		} else {
			c.ui.Error(err.Error())
			return 1
		}
	} else if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	buf, err = updateEnv(bufio.NewScanner(buf), *key, *value)
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

func (c *addCommand) Help() string {
	helpText := `
Usage: mayoiga add -key [KEY] -value [VALUE]

  Add env value and save to S3.
`
	return strings.TrimSpace(helpText)
}

func (c *addCommand) Synopsis() string {
	return "Add env value and save to S3"
}
