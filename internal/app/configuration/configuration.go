package configuration

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Configuration struct {
	StoreServiceConfig `yaml:"store_service"`
}

func NewConfig(path string) (*Configuration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("[ NewConfig ] could not open file: " + err.Error())
	}
	configData, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.New("[ NewConfig ] error while reading file: " + err.Error())
	}

	config := &Configuration{}

	if err := yaml.Unmarshal(configData, config); err != nil {
		return nil, errors.New("[ NewConfig ] error while parsing config: " + err.Error())
	}

	return config, nil
}