package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const (
	version    = 1
	configDir  = "./.mayoiga"
	configPath = "./.mayoiga/mayoiga.json"
)

type config struct {
	Version      int
	Region       *string
	S3Bucket     *string
	S3Key        *string
	MappingS3Key *string
	KMSKeyID     *string
}

type configError struct {
	err string
}

var cachedConfig *config

func assertConfig() (err error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New(configPath + " not exist. run mayoiga configure first")
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

	if config.MappingS3Key == nil || len(*(config.MappingS3Key)) == 0 {
		return errors.New("MappingS3Key is not set. run mayoiga configure first")
	}

	return nil
}

func readConfig() (*config, error) {
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	file, err := ioutil.ReadFile(configPath)
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
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return map[string]interface{}{
			"Version": version,
		}, nil
	}

	file, err := ioutil.ReadFile(configPath)
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

	j, err = json.MarshalIndent(config, "", "  ")
	if err != nil {
		return
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
		return
	}
	defer file.Close()

	_, err = file.Write(j)
	return
}
