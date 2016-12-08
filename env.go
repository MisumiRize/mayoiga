package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"bytes"

	"github.com/mitchellh/cli"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const envPath = "./.mayoiga/env"

func updateEnv(scanner *bufio.Scanner, key string, value string) (*bytes.Buffer, error) {
	found := false
	buf := new(bytes.Buffer)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "=")
		if len(tokens) >= 2 {
			k := tokens[0]
			v := line[len(k)+1:]
			if k == key {
				found = true
				v = value
			}
			if _, err := buf.WriteString(fmt.Sprintln(fmt.Sprintf("%s=%s", k, v))); err != nil {
				return nil, err
			}
		}
	}

	if !found {
		if _, err := buf.WriteString(fmt.Sprintln(fmt.Sprintf("%s=%s", key, value))); err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func deleteEnv(scanner *bufio.Scanner, key string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "=")
		if len(tokens) >= 2 && tokens[0] != key {
			if _, err := buf.WriteString(fmt.Sprintln(line)); err != nil {
				return nil, err
			}
		}
	}

	return buf, nil
}

func readEnvFile() ([]byte, error) {
	bs, err := ioutil.ReadFile(envPath)
	if os.IsNotExist(err) {
		return make([]byte, 0, 100), nil
	} else if err != nil {
		return nil, err
	}

	return bs, nil
}

func compareWithEnvFile(ui cli.Ui, compareTo string) (err error) {
	bs, err := readEnvFile()
	if err != nil {
		return
	}

	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(string(bs), compareTo)
	diffs := dmp.DiffMain(a, b, false)
	result := dmp.DiffCharsToLines(diffs, c)

	for _, v := range result {
		switch v.Type {
		case diffmatchpatch.DiffInsert:
			ui.Info(fmt.Sprintf("+ %s", v.Text))
		case diffmatchpatch.DiffEqual:
			ui.Info(v.Text)
		case diffmatchpatch.DiffDelete:
			ui.Info(fmt.Sprintf("- %s", v.Text))
		}
	}
	return
}

func parseEnv(scanner *bufio.Scanner) map[string]string {
	env := map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "=")
		if len(tokens) >= 2 {
			k := tokens[0]
			v := line[len(k)+1:]
			env[k] = v
		}
	}

	return env
}

func writeEnvFile(env []byte) (err error) {
	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(configDir, 0755); err != nil {
			return
		}
	} else if err != nil {
		return
	}

	file, err := os.Create(envPath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(env)
	return
}
