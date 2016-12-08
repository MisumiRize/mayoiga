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
		return errors.New("map flag should key=value format")
	}

	key := tokens[0]
	flags[key] = value[len(key)+1:]
	*f = flags

	return nil
}
