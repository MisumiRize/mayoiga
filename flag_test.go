package main

import (
	"flag"
	"testing"
)

func TestMappingFlags(t *testing.T) {
	var mappings mappingFlags

	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	flags.Var(&mappings, "mapping", "mapping")

	args := []string{"-mapping", "foo:bar", "-mapping", "foo:baz:qux"}
	if err := flags.Parse(args); err != nil {
		t.Fatalf(err.Error())
	}

	if mapping, ok := mappings["foo"]; ok {
		if mapping.Variables[0] != "bar" {
			t.Fatalf("mapping does not contain bar variable")
		}

		if alias, ok := mapping.Aliases["baz"]; ok {
			if alias != "qux" {
				t.Fatalf("alias assertion failed. actual: %s", alias)
			}
		} else {
			t.Fatalf("mapping does not contain baz alias")
		}
	} else {
		t.Fatalf("key foo does not exist")
	}
}

func TestMappingFlags_AddingMultipleVariables(t *testing.T) {
	var mappings mappingFlags

	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	flags.Var(&mappings, "mapping", "mapping")

	args := []string{"-mapping", "foo:bar", "-mapping", "foo:baz"}
	if err := flags.Parse(args); err != nil {
		t.Fatalf(err.Error())
	}

	if mapping, ok := mappings["foo"]; ok {
		if len(mapping.Variables) != 2 {
			t.Fatalf("mapping does not contain 2 variables")
		}
	} else {
		t.Fatalf("key foo does not exist")
	}
}
