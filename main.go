package main

import (
	"flag"
	"fmt"

	"github.com/cego/git-request-list/github"
	"github.com/cego/git-request-list/gitlab"
	"github.com/cego/git-request-list/gitrequest"
)

func main() {
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()

	conf, err := readConfig("conf.yml")
	if err != nil {
		panic(err)
	}

	err = conf.check()
	if err != nil {
		panic(err)
	}

	var requests []gitrequest.Request

	for _, sConf := range conf.Sources {
		var source interface {
			GetRequests() ([]gitrequest.Request, error)
			SetVerbose(bool)
		}

		switch sConf.API {
		case "gitlab":
			source, err = gitlab.New(sConf.Host, sConf.Token)
			break
		case "github":
			source, err = github.New(sConf.Host, sConf.User, sConf.Token)
			break
		}

		if err != nil {
			panic(err)
		}

		source.SetVerbose(*verbose)

		sRequests, err := source.GetRequests()
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
