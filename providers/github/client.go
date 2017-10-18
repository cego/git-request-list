package github

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cego/git-request-list/providers"
)

// Client represents a Github pull-request source.
type Client struct {
	http    http.Client
	host    string
	token   string
	verbose bool
}

// pullRequest serves as Unmarshal target type when reading Github API responses
type pullRequest struct {
	Name    string    `json:"title"`
	State   string    `json:"state"`
	URL     string    `json:"url"`
	Created time.Time `json:"created_at"`
	Updated time.Time `json:"updated_at"`
}

func init() {
	factory := func(host, token string, verbose bool) (providers.Provider, error) {
		if token == "" {
			return nil, errors.New("a github access token is required")
		}

		c := Client{}
		c.http = http.Client{}
		c.host = host
		c.token = token
		c.verbose = verbose

		return &c, nil
	}

	providers.RegisterProvider("github", factory)
}

// GetRequests returns a slice of pull-requests visible to the Client c. If acceptedRepositories is not empty, only
// pull-requests from the repositories whose name is included in acceptedRepositories are returned.
func (c *Client) GetRequests(acceptedRepositories []string) ([]providers.Request, error) {
	whitelist := map[string]bool{}
	for _, repository := range acceptedRepositories {
		whitelist[repository] = true
	}

	repositories, err := c.getRepositories()
	if err != nil {
		return nil, err
	}

	var result []providers.Request
	for _, repository := range repositories {
		if len(whitelist) > 0 && !whitelist[repository] {
			continue
		}

		requests, err := c.getRequests(repository)
		if err != nil {
			return nil, err
		}

		result = append(result, requests...)
	}

	return result, nil
}

// getRepositories gets the full names of repositories visible to c.
func (c *Client) getRepositories() ([]string, error) {
	resp, err := c.get("/user/repos")
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

// getRequests returns all pull-requests of the repository with the given name visible to c.
func (c *Client) getRequests(repos string) ([]providers.Request, error) {
	resp, err := c.get("/repos/" + repos + "/pulls")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var requests []pullRequest
	err = json.NewDecoder(resp.Body).Decode(&requests)
	if err != nil {
		return nil, err
	}

	var result []providers.Request
	for _, r := range requests {
		result = append(result, providers.Request{
			Repository: repos,
			Name:       r.Name,
			State:      r.State,
			URL:        r.URL,
			Created:    r.Created,
			Updated:    r.Updated,
		})
	}

	return result, nil
}

// get completes a HTTP request to the Github API represented by c.
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
