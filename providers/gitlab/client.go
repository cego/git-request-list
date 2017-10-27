package gitlab

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/cego/git-request-list/providers"
)

// Client represents a Gitlab merge-request source.
type Client struct {
	http    http.Client
	host    string
	token   string
	verbose bool
}

// repository serves as Unmarshall target type when reading Gitlab API responses
type repository struct {
	Name string `json:"path_with_namespace"`
	ID   int    `json:"id"`
}

// mergeRequest serves as Unmarshal target type when reading Gitlab API responses
type mergeRequest struct {
	Name    string    `json:"title"`
	URL     string    `json:"web_url"`
	Created time.Time `json:"created_at"`
	Updated time.Time `json:"updated_at"`
	WIP     bool      `json:"work_in_progress"`
}

func init() {
	factory := func(host, token string, verbose bool) (providers.Provider, error) {
		if token == "" {
			return nil, errors.New("a gitlab access token is required")
		}

		c := Client{}
		c.http = http.Client{}
		c.host = host
		c.token = token
		c.verbose = verbose

		return &c, nil
	}

	providers.RegisterProvider("gitlab", factory)
}

// GetRequests returns a slice of merge-requests visible to the Client c. Only merge-requests from the repositories
// whose name matches repositoryFilter are returned.
func (c *Client) GetRequests(repositoryFilter regexp.Regexp) ([]providers.Request, error) {
	var result []providers.Request

	repositories, err := c.getRepositories(repositoryFilter)
	if err != nil {
		return nil, err
	}

	for _, repository := range repositories {
		requests, err := c.getRequests(repository)
		if err != nil {
			return nil, err
		}

		result = append(result, requests...)
	}

	return result, nil
}

// getRepositories gets the repositories visible to c. Only repositories whose name matches filter are returned.
func (c *Client) getRepositories(filter regexp.Regexp) ([]repository, error) {
	var result []repository

	resp, err := c.get("HEAD", "/projects")
	if err != nil {
		return nil, err
	}
	pageCount, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
	if err != nil {
		return nil, err
	}

	for p := 1; p <= pageCount; p++ {
		resp, err = c.get("GET", "/projects?simple=1&with_merge_requests_enabled=1&page="+strconv.Itoa(p))
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var page []repository
		err = json.NewDecoder(resp.Body).Decode(&page)
		if err != nil {
			return nil, err
		}

		for _, r := range page {
			if !filter.MatchString(r.Name) {
				continue
			}
			result = append(result, r)
		}
	}

	return result, nil
}

// getRequests returns all open merge-requests of the repository with the given ID visible to c.
func (c *Client) getRequests(repos repository) ([]providers.Request, error) {

	resp, err := c.get("HEAD", "/projects/"+strconv.Itoa(repos.ID)+"/merge_requests?state=opened")
	if err != nil {
		return nil, err
	}
	pageCount, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
	if err != nil {
		return nil, err
	}

	var result []providers.Request
	for p := 1; p <= pageCount; p++ {
		resp, err := c.get("GET", "/projects/"+strconv.Itoa(repos.ID)+"/merge_requests?state=opened&page="+strconv.Itoa(p))
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var page []mergeRequest
		err = json.NewDecoder(resp.Body).Decode(&page)
		if err != nil {
			return nil, err
		}

		for _, r := range page {
			if r.WIP {
				continue
			}

			result = append(result, providers.Request{
				Repository: repos.Name,
				Name:       r.Name,
				URL:        r.URL,
				Created:    r.Created,
				Updated:    r.Updated,
			})
		}
	}

	return result, nil
}

// get completes a HTTP request to the Gitlab API represented by c.
func (c *Client) get(method string, path string) (*http.Response, error) {
	if c.verbose {
		log.Printf("%s %s/api/v4%s", method, c.host, path)
	}

	req, err := http.NewRequest(method, c.host+"/api/v4"+path, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("PRIVATE-TOKEN", c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
