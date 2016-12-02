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
	S3Bucket *string
	S3Key    *string
	KMSKeyID *string
	OutFile  *string
}

type configError struct {
	err string
}

func (e *configError) Error() string {
	return (*e).err
}

var cachedConfig map[string]interface{}

func assertConfig() (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New(path + " not exist. run mayoiga configure first")
	}

	s3Bucket, err := getStringConfig("S3Bucket")
	if err != nil || len(*s3Bucket) == 0 {
		return errors.New("S3Bucket is not set. run mayoiga configure first")
	}

	s3Key, err := getStringConfig("S3Key")
	if err != nil || len(*s3Key) == 0 {
		return errors.New("S3Key is not set. run mayoiga configure first")
	}

	return nil
}

func getStringConfig(key string) (*string, error) {
	if cachedConfig == nil {
		var err error
		cachedConfig, err = readConfigAsMap()
		if err != nil {
			return nil, err
		}
	}

	if value, ok := cachedConfig[key].(string); ok {
		return &value, nil
	}

	return nil, &configError{err: "error"}
}

func getStringConfigWithDefault(key, defaultValue string) (config *string, err error) {
	config, err = getStringConfig(key)
	if err == nil {
		return
	}

	switch err.(type) {
	case *configError:
		return &defaultValue, nil
	}

	return
}

func readConfig() (*config, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return &config{
			Version: version,
		}, nil
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config config
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
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
