package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cego/git-request-list/github"
	"github.com/cego/git-request-list/gitlab"
	"github.com/cego/git-request-list/gitrequest"
)

func main() {
	verbose := flag.Bool("v", false, "verbose")
	configPath := flag.String("c", "/etc/git-request-list.yml", "config path")
	flag.Parse()

	conf, err := readConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = conf.check()
	if err != nil {
		log.Fatal(err)
	}

	var requests []gitrequest.Request

	for _, sConf := range conf.Sources {
		var source interface {
			GetRequests(repositories []string) ([]gitrequest.Request, error)
		}

		switch sConf.API {
		case "gitlab":
			source, err = gitlab.New(sConf.Host, sConf.Token, sConf.SkipWIP, *verbose)
			break
		case "github":
			source, err = github.New(sConf.Host, sConf.User, sConf.Token, *verbose)
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		sRequests, err := source.GetRequests(sConf.Repositories)
		if err != nil {
			log.Fatal(err)
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
