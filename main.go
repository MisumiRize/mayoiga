package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mitchellh/cli"
)

func main() {
	ui := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorGreen,
		WarnColor:   cli.UiColorYellow,
		ErrorColor:  cli.UiColorRed,
		Ui:          &cli.BasicUi{Writer: os.Stdout},
	}

	s, err := session.NewSession()
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}
	sess = s

	commands := map[string]cli.CommandFactory{
		"add": func() (cli.Command, error) {
			return &addCommand{ui: ui}, nil
		},
		"remove": func() (cli.Command, error) {
			return &removeCommand{ui: ui}, nil
		},
		"pull": func() (cli.Command, error) {
			return &pullCommand{ui: ui}, nil
		},
		"configure": func() (cli.Command, error) {
			return &configureCommand{ui: ui}, nil
		},
	}

	cli := &cli.CLI{
		Name:     "mayoiga",
		Version:  "0.1.0",
		Args:     os.Args[1:],
		Commands: commands,
	}

	status, err := cli.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(status)
}
