package request

import (
	"time"
)

// Request represents a pull- or merge-request
type Request struct {
	Repository string
	Name       string
	URL        string
	Created    time.Time
	Updated    time.Time
}
