package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestAssertConfig_RequiresRegion(t *testing.T) {
	resetConfig(t)

	data := map[string]string{
		"S3Bucket": "bucket",
		"S3Key":    "key",
	}
	if err := writeConfig(data); err != nil {
		t.Fatalf(err.Error())
	}

	err := assertConfig()
	if err == nil {
		t.Fatalf("error expected, but there is no error")
	}
}

func TestWriteConfig(t *testing.T) {
	resetConfig(t)

	data := map[string]string{
		"S3Bucket":   "test_key",
		"InvalidKey": "another_key",
	}
	if err := writeConfig(data); err != nil {
		t.Fatalf(err.Error())
	}

	config, err := readConfig()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if config.Version != version {
		t.Fatalf("Version assertion failed. actual: %s", config.Version)
	}

	if *config.S3Bucket != "test_key" {
		t.Fatalf("S3Bucket assertion failed. actual: %s", *config.S3Bucket)
	}
}

func resetConfig(t *testing.T) {
	cachedConfig = nil
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			t.Fatalf("config file remove failed")
		}
	}
}

func createJSON(t *testing.T, data map[string]string) {
	json, err := json.Marshal(data)
	if err != nil {
		t.Fatalf(err.Error())
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer file.Close()

	if _, err = file.Write(json); err != nil {
		t.Fatalf(err.Error())
	}
}
