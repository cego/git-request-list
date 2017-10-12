package main

import (
    "os"
    "fmt"

    "github.com/cego/git-request-list/github"
    "github.com/cego/git-request-list/gitlab"
    "github.com/cego/git-request-list/output"
)

func main () {

    var requests []output.Request

    // hub

    github, err := github.New()
    if err != nil {
        panic(err)
    }

    github.SetUser(os.Args[1])
    github.SetToken(os.Args[2])

    githubRequests, err := github.GetRequests()
    if err != nil {
        panic(err)
    }

    for i, _ := range(githubRequests) {
        requests = append(requests, &githubRequests[i])
    }

    // lab

    gitlab, err := gitlab.New()
    if err != nil {
        panic(err)
    }

    gitlab.SetToken(os.Args[3])

    gitlabRequests, err := gitlab.GetRequests()
    if err != nil {
        panic(err)
    }

    for i, _ := range(gitlabRequests) {
        requests = append(requests, &gitlabRequests[i])
    }

    // out

    table := output.NewTable()
    for _, r := range(requests) {
        table.Add(r)
    }

    fmt.Print(table.String())
}
