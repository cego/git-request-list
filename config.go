package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// config serves as Unmarshall target of configuration files
type config struct {
	SortBy  string `yaml:"sort_by"`
	Sources []struct {
		API          string   `yaml:"api"`
		Host         string   `yaml:"host"`
		Token        string   `yaml:"token"`
		Repositories []string `yaml:"repositories"`
	} `yaml:"sources"`
}

// readConfig reads the configuration file at the given path
func readConfig(path string) (*config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c config
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
