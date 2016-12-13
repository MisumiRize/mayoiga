package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestMapping_IsEmpty(t *testing.T) {
	table := []struct {
		mapping  mapping
		expected bool
	}{
		{
			mapping{},
			true,
		},
		{
			mapping{Variables: []string{}},
			true,
		},
	}

	for _, entry := range table {
		if actual := entry.mapping.isEmpty(); actual != entry.expected {
			t.Fatalf("isEmpty() assertion failured, mapping: %#v", entry.mapping)
		}
	}
}

func TestMapping_Merge(t *testing.T) {
	table := []struct {
		from, to, expected mapping
	}{
		{
			mapping{
				Variables: []string{"FOO"},
				Aliases:   map[string]string{"BAR": "baz"},
			},
			mapping{},
			mapping{
				Variables: []string{"FOO"},
				Aliases:   map[string]string{"BAR": "baz"},
			},
		},
		{
			mapping{
				Variables: []string{"FOO"},
			},
			mapping{},
			mapping{
				Variables: []string{"FOO"},
				Aliases:   map[string]string{},
			},
		},
	}

	for _, entry := range table {
		if actual := entry.to.merge(entry.from); !reflect.DeepEqual(actual, entry.expected) {
			t.Fatalf("assertion failed, expected: %#v, actual %#v", entry.expected, actual)
		}
	}
}

func TestMapping_BuildMappedEnvBody(t *testing.T) {
	table := []struct {
		mapping  mapping
		env      map[string]string
		expected string
	}{
		{
			mapping{
				Variables: []string{"FOO"},
				Aliases:   map[string]string{"FOO": "BAR"},
			},
			map[string]string{"FOO": "foo"},
			`
FOO=foo
BAR=foo
`,
		},
	}

	for _, entry := range table {
		buf, err := entry.mapping.buildMappedEnvBody(entry.env)
		if err != nil {
			t.Fatalf(err.Error())
		}

		actual := buf.String()
		if strings.TrimSpace(entry.expected) != strings.TrimSpace(actual) {
			t.Fatalf("assertion failed, expected: %s, actual: %s", entry.expected, actual)
		}
	}
}
