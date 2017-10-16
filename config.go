package main

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	Sources []struct {
		API          string   `yaml:"api"`
		Host         string   `yaml:"host"`
		Token        string   `yaml:"token"`
		SkipWIP      bool     `yaml:"skip_wip"`
		Repositories []string `yaml:"repositories"`
	} `yaml:"sources"`
}

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

func (c *config) check() error {
	if len(c.Sources) == 0 {
		return errors.New("no sources defined")
	}

	for _, s := range c.Sources {
		switch s.API {
		case "gitlab":
		case "github":
			break
		default:
			return errors.New("the `api` of each source must be either `gitlab` or `github`")
		}
	}

	return nil
}
