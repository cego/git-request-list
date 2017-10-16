package github

import (
	"errors"
	"log"
	"strconv"

	"encoding/json"
	"net/http"

	"github.com/cego/git-request-list/gitrequest"
)

type Client struct {
	http    http.Client
	host    string
	token   string
	verbose bool
}

func New(host, token string, verbose bool) (*Client, error) {
	c := Client{}

	c.http = http.Client{}
	c.host = host
	c.token = token
	c.verbose = verbose

	return &c, nil
}

func (c *Client) GetRequests(acceptedRepositories []string) ([]gitrequest.Request, error) {
	whitelist := map[string]bool{}
	for _, repository := range acceptedRepositories {
		whitelist[repository] = true
	}

	var result []gitrequest.Request

    user, err := c.getUser()
    if err != nil {
        return nil, err
    }

	repositories, err := c.getRepositories(user)
	if err != nil {
		return nil, err
	}

	for _, repository := range repositories {
		if len(whitelist) > 0 && !whitelist[repository] {
			continue
		}

		requests, err := c.getRequests(repository)
		if err != nil {
			return nil, err
		}

		for i := range requests {
			requests[i].RepositoryValue = repository
			result = append(result, &requests[i])
		}
	}

	return result, nil
}

func (c *Client) getUser() (string, error) {
	resp, err := c.get("/user")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var user struct {
		Login string `json:"login"`
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return "", err
	}

	return user.Login, nil
}

func (c *Client) getRepositories(user string) ([]string, error) {
	resp, err := c.get("/users/" + user + "/repos?type=all")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var repos []struct {
		Name string `json:"full_name"`
	}

	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, r := range repos {
		names = append(names, r.Name)
	}

	return names, nil
}

func (c *Client) getRequests(repos string) ([]Request, error) {
	resp, err := c.get("/repos/" + repos + "/pulls")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var requests []Request

	err = json.NewDecoder(resp.Body).Decode(&requests)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (c *Client) get(path string) (*http.Response, error) {
	host := "https://api.github.com"
	if c.host != "" {
		host = c.host
	}

	if c.verbose {
		log.Printf("GET %s%s", host, path)
	}

	req, err := http.NewRequest("GET", host+path, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	remainder, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	if err != nil {
		return nil, err
	}
	if remainder <= 0 {
		return nil, errors.New("Github API rate limit exceeded.")
	}

	return resp, nil
}
