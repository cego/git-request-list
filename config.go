package main

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	Sources []struct {
		API   string `yaml:"api"`
		Host  string `yaml:"host"`
		User  string `yaml:"user"`
		Token string `yaml:"token"`
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
	for _, s := range c.Sources {
		switch s.API {
		case "gitlab":
			break
		case "github":
			if s.User == "" {
				return errors.New("all `github` api sources must specify a `user`")
			}
			break
		default:
			return errors.New("the `api` of each source must be either `gitlab` or `github`")
		}
	}

	return nil
}
