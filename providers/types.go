package providers

import (
	"time"
)

// Request represents a pull- or merge-request
type Request struct {
	Repository string
	Name       string
	State      string
	URL        string
	Created    time.Time
	Updated    time.Time
}

// Provider is the common interface for all providers of pull- and merge-requests
type Provider interface {
	GetRequests(repositories []string) ([]Request, error)
}

// ProviderFactory types a function for producing new Providers
type ProviderFactory func(host, token string, verbose bool) (Provider, error)
