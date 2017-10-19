package providers

import (
	"errors"
)

var factories map[string]ProviderFactory

func init() {
	factories = map[string]ProviderFactory{}
}

// RegisterProvider registers a Provider implementation via a factory function
func RegisterProvider(identifier string, factory ProviderFactory) {
	factories[identifier] = factory
}

// GetProvider gets a Provider implementation of a type previously registered with RegisterProvider
func GetProvider(identifier, host, token string, verbose bool) (Provider, error) {
	factory, exists := factories[identifier]
	if !exists {
		return nil, errors.New("unknown provider identifier")
	}

	return factory(host, token, verbose)
}
