package main

import (
	"flag"
	"sort"
	"strings"

	"github.com/mitchellh/cli"
)

type mapCommand struct {
	ui cli.Ui
}

func (c *mapCommand) Run(args []string) int {
	if len(args) < 1 || args[0][:1] == "-" {
		c.ui.Error(c.Help())
		return 1
	}

	var variables stringFlags
	var aliases mapFlags

	flags := flag.NewFlagSet("var", flag.ContinueOnError)
	flags.Var(&variables, "variable", "variable")
	flags.Var(&aliases, "alias", "alias")
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

	if err = compareWithEnvFile(c.ui, buf.String()); err != nil {
		c.ui.Error(err.Error())
		return 1
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

	mapping := mappings[args[0]]
	newVariables := mapping.Variables
	newAliases := mapping.Aliases

	for _, v := range variables {
		newVariables = addVariable(newVariables, v)
	}

	for k, v := range aliases {
		newAliases[k] = v
	}

	sort.Strings(newVariables)

	mapping.Variables = newVariables
	mapping.Aliases = newAliases
	mappings[args[0]] = mapping

	if err = writeMappingsFile(mappings); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if err = putMappingsToS3(mappings); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	if err = generateMappedEnvFiles(buf, mappings); err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *mapCommand) Help() string {
	helpText := `
Usage: mayoiga map <FILE> -variable <KEY> -alias <KEY=ALIAS>

  Add mapping to file and save to S3.

Options:

	-variable  Define variable usage.
	-alias     Define variable alias.
`
	return strings.TrimSpace(helpText)
}

func (c *mapCommand) Synopsis() string {
	return "Add mapping to file and save to S3"
}

func addVariable(variables []string, variable string) []string {
	for _, v := range variables {
		if v == variable {
			return variables
		}
	}
	return append(variables, variable)
}