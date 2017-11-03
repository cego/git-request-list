package providers

import (
	"regexp"

	"github.com/cego/git-request-list/request"
)

// Provider is the common interface for all providers of pull- and merge-requests
type Provider interface {
	GetRequests(repositoryFilter regexp.Regexp) ([]request.Request, error)
}

// ProviderFactory types a function for producing new Providers
type ProviderFactory func(host, token string, verbose bool) (Provider, error)
