package providers

import (
	"time"
)

// Request is the common interface for pull- and merge-requests
type Request interface {
	Repository() string
	Name() string
	State() string
	URL() string
	Created() time.Time
	Updated() time.Time
}

