package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/cego/git-request-list/providers/github"
	"github.com/cego/git-request-list/providers/gitlab"
	"github.com/cego/git-request-list/providers"
	"github.com/cego/git-request-list/formatters"
)

func main() {
	// Read flags and configuration file

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

	// Gather requests from configured sources

	var requests []providers.Request
	for _, sConf := range conf.Sources {
		var source interface {
			GetRequests(repositories []string) ([]providers.Request, error)
		}

		switch sConf.API {
		case "gitlab":
			source, err = gitlab.New(sConf.Host, sConf.Token, *verbose)
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

	// Output the requests

	switch conf.SortBy {
	case "name":
		sort.Sort(formatters.ByName(requests))
		break
	case "state":
		sort.Sort(formatters.ByState(requests))
		break
	case "url":
		sort.Sort(formatters.ByURL(requests))
		break
	case "created":
		sort.Sort(formatters.ByCreated(requests))
		break
	case "updated":
		sort.Sort(formatters.ByUpdated(requests))
		break
	case "repository":
	default:
		sort.Sort(formatters.ByRepository(requests))
		break
	}

	table := formatters.Table{}
	fmt.Print(table.String(requests...))
}
