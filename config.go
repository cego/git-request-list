package main

import (
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// config represents a fully parsed configuration file
type config struct {
	SortBy  string
	Format  string
	Sources []sourceConfig
}

// sourceConfig holds configuration parameters of a single source
type sourceConfig struct {
	API          string
	Host         string
	Token        string
	Repositories regexp.Regexp
}

// readConfig reads the configuration file at the given path
func readConfig(path string) (*config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var unmarshallTarget struct {
		SortBy  string `yaml:"sort_by"`
		Format  string `yaml:"format"`
		Sources []struct {
			API          string   `yaml:"api"`
			Host         string   `yaml:"host"`
			Token        string   `yaml:"token"`
			Repositories []string `yaml:"repositories"`
		} `yaml:"sources"`
	}

	err = yaml.Unmarshal(bytes, &unmarshallTarget)
	if err != nil {
		return nil, err
	}

	c := config{
		SortBy:  unmarshallTarget.SortBy,
		Format:  unmarshallTarget.Format,
		Sources: make([]sourceConfig, len(unmarshallTarget.Sources)),
	}

	for i, s := range unmarshallTarget.Sources {
		re, err := regexp.Compile(strings.Join(s.Repositories, "|"))
		if err != nil {
			return nil, err
		}

		c.Sources[i] = sourceConfig{
			API:          s.API,
			Host:         s.Host,
			Token:        s.Token,
			Repositories: *re,
		}
	}

	return &c, nil
}
