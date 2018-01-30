package formatters

import (
	"errors"
	"time"

	"github.com/cego/git-request-list/request"
)

// Formatter is the common interface for all formatters of pull- and merge-requests
type Formatter interface {
	String() string
}

// Arguments defines what a formatter need to know to generate its output.
type Arguments struct {
	Requests []request.Request
	Timezone *time.Location
}

// FormatterFactory types a function for producing new Formatters
type FormatterFactory func(arguments Arguments) (Formatter, error)

var factories map[string]FormatterFactory

func init() {
	factories = map[string]FormatterFactory{}
}

// RegisterFormatter registers a Formatter implementation via a factory function
func RegisterFormatter(identifier string, factory FormatterFactory) {
	factories[identifier] = factory
}

// GetFormatter gets a Formatter implementation of a type previously registered with RegisterFormatter
func GetFormatter(identifier string, arguments Arguments) (Formatter, error) {
	factory, exists := factories[identifier]
	if !exists {
		return nil, errors.New("unknown provider identifier")
	}

	return factory(arguments)
}
