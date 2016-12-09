package main

import (
	"errors"
	"strings"
)

type stringFlags []string

func (f *stringFlags) String() string {
	return strings.Join(*f, "\n")
}

func (f *stringFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

type mapFlags map[string]string

func (f *mapFlags) String() string {
	return ""
}

func (f *mapFlags) Set(value string) error {
	flags := *f
	if flags == nil {
		flags = map[string]string{}
	}

	tokens := strings.Split(value, "=")
	if len(tokens) < 2 {
		return errors.New("this flag should be key=value format")
	}

	key := tokens[0]
	flags[key] = value[len(key)+1:]
	*f = flags

	return nil
}

type mappingFlags map[string]mapping

func (f *mappingFlags) String() string {
	return ""
}

func (f *mappingFlags) Set(value string) error {
	flags := *f
	if flags == nil {
		flags = map[string]mapping{}
	}

	tokens := strings.Split(value, ":")
	if len(tokens) < 2 {
		return errors.New("this flag should be file:variable or file:variable:alias format")
	}

	file := tokens[0]
	variable := tokens[1]
	if len(tokens) > 2 {
		alias := strings.Join(tokens[2:], ":")

		if m, ok := flags[file]; ok {
			m.Aliases[variable] = alias
			flags[file] = m
		} else {
			flags[file] = mapping{
				Variables: []string{},
				Aliases:   map[string]string{variable: alias},
			}
		}
	} else {
		if m, ok := flags[file]; ok {
			m.Variables = addIfNotExist(m.Variables, variable)
			flags[file] = m
		} else {
			flags[file] = mapping{
				Variables: []string{variable},
				Aliases:   map[string]string{},
			}
		}
	}

	*f = flags

	return nil
}
