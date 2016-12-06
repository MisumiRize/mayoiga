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

func updateEnv(scanner *bufio.Scanner, key string, value string) (*bytes.Buffer, error) {
	found := false
	buf := new(bytes.Buffer)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "=")
		if len(tokens) >= 2 {
			k := tokens[0]
			v := strings.Join(tokens[1:], "=")
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

func compareWithEnvFile(ui cli.Ui, compareTo string) (err error) {
	env, err := readEnvFile()
	if err != nil {
		return
	}

	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(string(env), compareTo)
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

func readEnvFile() ([]byte, error) {
	config, err := readConfig()
	if err != nil {
		return nil, err
	}

	outFile := config.OutFile
	if outFile == nil {
		outFile = config.S3Key
	}

	bs, err := ioutil.ReadFile(*outFile)
	if os.IsNotExist(err) {
		return make([]byte, 0, 100), nil
	} else if err != nil {
		return nil, err
	}

	return bs, nil
}

func writeEnvFile(env []byte) (err error) {
	config, err := readConfig()
	if err != nil {
		return
	}

	outFile := config.OutFile
	if outFile == nil {
		outFile = config.S3Key
	}

	file, err := os.Create(*outFile)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(env)
	return
}
