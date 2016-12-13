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

func TestReadConfigAsMap_ReturnsNoError_WhenFileIsAbsent(t *testing.T) {
	resetConfig(t)

	config, err := readConfigAsMap()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(config) != 1 {
		t.Fatalf("Empty config should have 1 entry, but there are %d", len(config))
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

	stat, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !stat.IsDir() {
		t.Fatalf(configDir + " should be directory, but it is not")
	}

	stat, err = os.Stat(configPath)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if stat.IsDir() {
		t.Fatalf(configPath + " should not be directory, but it is")
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
	if err := os.RemoveAll(configDir); err != nil {
		t.Fatalf("config file remove failed")
	}
}

func createJSON(t *testing.T, data map[string]string) {
	json, err := json.Marshal(data)
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(configDir, 0755); err != nil {
			return
		}
	} else if err != nil {
		return
	}

	file, err := os.Create(configPath)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer file.Close()

	if _, err = file.Write(json); err != nil {
		t.Fatalf(err.Error())
	}
}
