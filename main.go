package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/cego/git-request-list/request"

	"github.com/cego/git-request-list/formatters"
	_ "github.com/cego/git-request-list/formatters/text"

	"github.com/cego/git-request-list/providers"
	_ "github.com/cego/git-request-list/providers/github"
	_ "github.com/cego/git-request-list/providers/gitlab"
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

	// Gather requests from configured sources

	var requests []request.Request
	for _, sConf := range conf.Sources {
		source, err := providers.GetProvider(sConf.API, sConf.Host, sConf.Token, *verbose)
		if err != nil {
			log.Fatal(err)
		}

		sRequests, err := source.GetRequests(sConf.Repositories)
		if err != nil {
			log.Fatal(err)
		}

		requests = append(requests, sRequests...)
	}

	// Output the requests

	switch conf.SortBy {
	case "name":
		sort.Sort(formatters.ByName(requests))
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

	formatter, err := formatters.GetFormatter("text", requests)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(formatter.String())
}
