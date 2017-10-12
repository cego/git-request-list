package gitlab

import (
    "log"
    "strconv"

    "net/http"
    "encoding/json"
)

type Client struct {
    http http.Client
    token string
}

type repository struct {
    Name string `json:"path_with_namespace"`
    ID int `json:"id"`
}

func New () (*Client, error) {
    c := Client{}

    c.http = http.Client{}

    return &c, nil
}

func (c *Client) SetToken(t string) {
    c.token = t
}

func (c *Client) GetRequests() ([]Request, error) {
    var result []Request

    repositories, err := c.getRepositories()
    if err != nil {
        return nil, err
    }

    for _, repository := range(repositories) {
        requests, err := c.getRequests(repository.ID)
        if err != nil {
            return nil, err
        }

        for _, request := range(requests) {
            request.RepositoryValue = repository.Name
            result = append(result, request)
        }
    }

    return result, nil
}

func (c *Client) getRepositories () ([]repository, error) {
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
        resp, err = c.get("GET", "/projects?simple=1&with_merge_requests_enabled=1&page=" + strconv.Itoa(p))
        if err != nil {
            return nil, err
        }

        defer resp.Body.Close()

        var page []repository
        err = json.NewDecoder(resp.Body).Decode(&page)
        if err != nil {
            return nil, err
        }

        result = append(result, page...)
    }

    return result, nil
}

func (c *Client) getRequests(repos int) ([]Request, error) {
    var result []Request

    resp, err := c.get("HEAD", "/projects/" + strconv.Itoa(repos) + "/merge_requests?state=opened")
    if err != nil {
        return nil, err
    }
    pageCount, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
    if err != nil {
        return nil, err
    }

    for p := 1; p <= pageCount; p++ {
        resp, err := c.get("GET", "/projects/" + strconv.Itoa(repos) + "/merge_requests?state=opened&page=" + strconv.Itoa(p))
        if err != nil {
            return nil, err
        }

        defer resp.Body.Close()

        var page []Request

        err = json.NewDecoder(resp.Body).Decode(&page)
        if err != nil {
            return nil, err
        }

        result = append(result, page...)
    }

    return result, nil
}

func (c *Client) get(method string, path string) (*http.Response, error) {
    log.Printf("%s https://gitlab.cego.dk/api/v4%s", method, path)

    req, err := http.NewRequest(method, "https://gitlab.cego.dk/api/v4" + path, nil)
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

