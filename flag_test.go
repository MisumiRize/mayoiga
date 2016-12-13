package main

import (
	"flag"
	"reflect"
	"testing"
)

func TestMappingFlags(t *testing.T) {
	table := []struct {
		args     []string
		expected mappingFlags
	}{
		{
			[]string{"-mapping", "foo:bar", "-mapping", "foo:baz:qux"},
			mappingFlags{
				"foo": mapping{
					[]string{"bar"},
					map[string]string{"baz": "qux"},
				},
			},
		},
		{
			[]string{"-mapping", "foo:bar", "-mapping", "foo:baz"},
			mappingFlags{
				"foo": mapping{Variables: []string{"bar", "baz"}},
			},
		},
	}

	for _, entry := range table {
		var actual mappingFlags
		flags := flag.NewFlagSet("test", flag.ContinueOnError)
		flags.Var(&actual, "mapping", "mapping")

		if err := flags.Parse(entry.args); err != nil {
			t.Fatalf(err.Error())
		}

		if !reflect.DeepEqual(actual, entry.expected) {
			t.Fatalf("assertion failed, expected: %#v, actual: %#v", entry.expected, actual)
		}
	}
}
