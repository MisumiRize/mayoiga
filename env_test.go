package main

import (
	"bufio"
	"os"
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

func TestUpdateEnv_DoesNotChangeExistingKey(t *testing.T) {
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
		if entry == "FOO=foo" {
			return
		}
	}
	t.Fatalf("FOO=foo not found")
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

func TestDeleteEnv(t *testing.T) {
	envText := `
FOO=foo
BAR=bar
BAZ=baz
`
	reader := strings.NewReader(envText)
	env, err := deleteEnv(bufio.NewScanner(reader), "BAR")
	if err != nil {
		t.Fatalf(err.Error())
	}

	entries := strings.Split(env.String(), "\n")
	for _, entry := range entries {
		if entry == "BAR=bar" {
			t.Fatalf("BAR=bar remains")
		}
	}
}

func TestDeleteEnv_KeyNotFound(t *testing.T) {
	envText := `
FOO=foo
BAR=bar
BAZ=baz
`
	reader := strings.NewReader(envText)
	env, err := deleteEnv(bufio.NewScanner(reader), "QUX")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if strings.TrimSpace(env.String()) != strings.TrimSpace(envText) {
		t.Fatalf("assertion failed actual: %s", env.String())
	}
}

func TestWriteEnvFile(t *testing.T) {
	resetConfig(t)

	if err := writeEnvFile([]byte("PAYLOAD")); err != nil {
		t.Fatalf(err.Error())
	}

	stat, err := os.Stat(envPath)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if stat.IsDir() {
		t.Fatalf(envPath + " should not be directory")
	}
}
