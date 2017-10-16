package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

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
			source, err = github.New(sConf.Host, sConf.Token, *verbose)
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

	switch conf.SortBy {
	case "name":
		sort.Sort(gitrequest.ByName(requests))
		break
	case "state":
		sort.Sort(gitrequest.ByState(requests))
		break
	case "url":
		sort.Sort(gitrequest.ByURL(requests))
		break
	case "created":
		sort.Sort(gitrequest.ByCreated(requests))
		break
	case "updated":
		sort.Sort(gitrequest.ByUpdated(requests))
		break
	case "repository":
	default:
		sort.Sort(gitrequest.ByRepository(requests))
		break
	}

	table := gitrequest.Table{}
	fmt.Print(table.String(requests...))
}
