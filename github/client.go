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
	http  http.Client
	host  string
	user  string
	token string
}

func New(host, user, token string) (*Client, error) {
	c := Client{}

	c.http = http.Client{}
	c.host = host
	c.user = user
	c.token = token

	return &c, nil
}

func (c *Client) SetUser(u string) {
	c.user = u
}

func (c *Client) SetToken(t string) {
	c.token = t
}

func (c *Client) SetHost(h string) {
	c.host = h
}

func (c *Client) GetRequests() ([]gitrequest.Request, error) {
	var result []gitrequest.Request

	repositories, err := c.getRepositories()
	if err != nil {
		return nil, err
	}

	for _, repository := range repositories {
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

func (c *Client) getRepositories() ([]string, error) {
	if c.user == "" {
		return nil, errors.New("No github user set.")
	}

	resp, err := c.get("/users/" + c.user + "/repos?type=all")
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

	log.Printf("GET %s%s", host, path)

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
