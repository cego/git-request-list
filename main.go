package main

import (
	"fmt"

	"github.com/cego/git-request-list/github"
	"github.com/cego/git-request-list/gitlab"
	"github.com/cego/git-request-list/gitrequest"
)

func main() {
	conf, err := readConfig("conf.yml")
	if err != nil {
		panic(err)
	}

	err = conf.check()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", conf)

	var requests []gitrequest.Request

	for _, sConf := range conf.Sources {
		var s interface {
			GetRequests() ([]gitrequest.Request, error)
		}

		switch sConf.API {
		case "gitlab":
			s, err = gitlab.New(sConf.Host, sConf.Token)
			break
		case "github":
			s, err = github.New(sConf.Host, sConf.User, sConf.Token)
			break
		}

		if err != nil {
			panic(err)
		}

		sRequests, err := s.GetRequests()
		if err != nil {
			panic(err)
		}

		for _, r := range sRequests {
			requests = append(requests, r)
		}
	}

	table := gitrequest.NewTable()
	for _, r := range requests {
		table.Add(r)
	}

	fmt.Print(table.String())
}
