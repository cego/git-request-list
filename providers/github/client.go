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
	"github.com/cego/git-request-list/request"
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

// GetRequests returns a slice of pull-requests visible to the Client c. Only pull-requests from the repositories whose
// name is matched by repositoryFilter are returned.
func (c *Client) GetRequests(repositoryFilter regexp.Regexp) ([]request.Request, error) {
	repositories, err := c.getRepositories(repositoryFilter)
	if err != nil {
		return nil, err
	}

	var result []request.Request
	for _, repository := range repositories {
		requests, err := c.getRequests(repository)
		if err != nil {
			return nil, err
		}

		result = append(result, requests...)
	}

	return result, nil
}

// getRepositories gets the full names of repositories visible to c. Only repository names matching filter are returned.
func (c *Client) getRepositories(filter regexp.Regexp) ([]string, error) {
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
			if !filter.MatchString(r.Name) {
				continue
			}
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
func (c *Client) getRequests(repos string) ([]request.Request, error) {
	var result []request.Request

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
			result = append(result, request.Request{
				Repository: repos,
				Name:       r.Name,
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
