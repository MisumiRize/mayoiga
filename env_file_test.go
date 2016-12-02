package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestUpdateEnv(t *testing.T) {
	envText := `
FOO=foo
BAR=bar
BAZ=baz
`
	reader := strings.NewReader(envText)
	env, err := updateEnv(bufio.NewScanner(reader), "BAR", "baz")
	if err != nil {
		t.Fatalf(err.Error())
	}

	entries := strings.Split(env.String(), "\n")
	for _, entry := range entries {
		if entry == "BAR=baz" {
			return
		}
	}
	t.Fatalf("BAR=baz not found")
}

func TestUpdateEnv_AppendKey(t *testing.T) {
	envText := `
FOO=foo
BAR=bar
BAZ=baz
`
	reader := strings.NewReader(envText)
	env, err := updateEnv(bufio.NewScanner(reader), "QUX", "qux")
	if err != nil {
		t.Fatalf(err.Error())
	}

	entries := strings.Split(env.String(), "\n")
	for _, entry := range entries {
		if entry == "QUX=qux" {
			return
		}
	}
	t.Fatalf("QUX=qux not found")
}

func TestUpdateEnv_UpdateMalformedKey(t *testing.T) {
	envText := `
FOO=foo
BAR=bar=bar
BAZ=baz
`
	reader := strings.NewReader(envText)
	env, err := updateEnv(bufio.NewScanner(reader), "BAR", "newbar")
	if err != nil {
		t.Fatalf(err.Error())
	}

	entries := strings.Split(env.String(), "\n")
	for _, entry := range entries {
		if entry == "BAR=bar=bar" {
			t.Fatalf("malformed entry exists")
		}
	}
}
