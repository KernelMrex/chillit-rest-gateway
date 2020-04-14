package configuration

import (
	"chillit-rest-gateway/internal/app/apiserver"
	"chillit-rest-gateway/internal/app/places"
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Configuration configures application
type Configuration struct {
	StoreService *places.Config    `yaml:"store_service"`
	APIServer    *apiserver.Config `yaml:"api_server"`
}

// ParseConfig parses from file
func ParseConfig(path string) (*Configuration, error) {
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
