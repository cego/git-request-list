package github

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"regexp"
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

// links represents the Link header used for pagination in the Github API
type links struct {
	Next string
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
	var result []string

	for next := "/user/repos"; next != ""; {
		resp, err := c.get(next)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var page []struct {
			Name string `json:"full_name"`
		}

		err = json.NewDecoder(resp.Body).Decode(&page)
		if err != nil {
			return nil, err
		}

		for _, r := range page {
			result = append(result, r.Name)
		}

		links, err := readPaginationLinks(resp.Header)
		if err != nil {
			return nil, err
		}

		next = links.Next
	}

	return result, nil
}

// getRequests returns all pull-requests of the repository with the given name visible to c.
func (c *Client) getRequests(repos string) ([]providers.Request, error) {
	var result []providers.Request

	for next := "/repos/" + repos + "/pulls"; next != ""; {
		resp, err := c.get(next)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var page []pullRequest

		err = json.NewDecoder(resp.Body).Decode(&page)
		if err != nil {
			return nil, err
		}

		for _, r := range page {
			result = append(result, providers.Request{
				Repository: repos,
				Name:       r.Name,
				State:      r.State,
				URL:        r.URL,
				Created:    r.Created,
				Updated:    r.Updated,
			})
		}

		links, err := readPaginationLinks(resp.Header)
		if err != nil {
			return nil, err
		}

		next = links.Next
	}

	return result, nil
}

// get completes a HTTP request to the Github API represented by c.
func (c *Client) get(path string) (*http.Response, error) {

	// Construct an url

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	if !u.IsAbs() {
		if c.host != "" {
			h, err := url.Parse(c.host)
			if err != nil {
				return nil, err
			}

			u.Scheme = h.Scheme
			u.Host = h.Host
		} else {
			u.Scheme = "https"
			u.Host = "api.github.com"
		}
	}

	// Build a HTTP request

	if c.verbose {
		log.Printf("GET %s", u.String())
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	// Get the response, failing if we hit the Gitlab API rate limit

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

// readPaginationLinks extracts links from headers of Gitlab API responses
// See https://developer.github.com/v3/#pagination
func readPaginationLinks(header http.Header) (links, error) {
	raw := header.Get("Link")
	if raw == "" {
		return links{}, nil
	}

	re := regexp.MustCompile("<([^>]+)>; rel=\"next\"")
	match := re.FindStringSubmatch(raw)

	if match == nil {
		return links{}, nil
	}

	return links{Next: match[1]}, nil
}