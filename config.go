package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const (
	version = 1
	path    = "./.mayoiga.json"
)

type config struct {
	Version  int
	Region   *string
	S3Bucket *string
	S3Key    *string
	KMSKeyID *string
	OutFile  *string
}

type configError struct {
	err string
}

var cachedConfig *config

func assertConfig() (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New(path + " not exist. run mayoiga configure first")
	}

	config, err := readConfig()
	if err != nil {
		return err
	}

	if config.Region == nil || len(*(config.Region)) == 0 {
		return errors.New("Region is not set. run mayoiga configure first")
	}

	if config.S3Bucket == nil || len(*(config.S3Bucket)) == 0 {
		return errors.New("S3Bucket is not set. run mayoiga configure first")
	}

	if config.S3Key == nil || len(*(config.S3Key)) == 0 {
		return errors.New("S3Key is not set. run mayoiga configure first")
	}

	return nil
}

func readConfig() (*config, error) {
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config config
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	cachedConfig = &config

	return cachedConfig, nil
}

func readConfigAsMap() (map[string]interface{}, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return map[string]interface{}{
			"Version": version,
		}, nil
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if config["Version"] == nil {
		config["Version"] = version
	}

	return config, nil
}

func writeConfig(newConfig map[string]string) (err error) {
	cfg, err := readConfigAsMap()
	if err != nil {
		return
	}

	for key, value := range newConfig {
		cfg[key] = value
	}

	j, err := json.Marshal(cfg)
	if err != nil {
		return
	}

	var config config
	if err = json.Unmarshal(j, &config); err != nil {
		return
	}

	j, err = json.Marshal(config)
	if err != nil {
		return
	}

	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(j)
	return
}
