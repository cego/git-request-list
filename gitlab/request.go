package gitlab

import (
	"time"
)

type Request struct {
	RepositoryValue string
	NameValue       string    `json:"title"`
	StateValue      string    `json:"state"`
	URLValue        string    `json:"web_url"`
	CreatedValue    time.Time `json:"created_at"`
	UpdatedValue    time.Time `json:"updated_at"`
}

func (r *Request) Repository() string {
	return r.RepositoryValue
}

func (r *Request) Name() string {
	return r.NameValue
}

func (r *Request) State() string {
	return r.StateValue
}

func (r *Request) URL() string {
	return r.URLValue
}

func (r *Request) Created() time.Time {
	return r.CreatedValue
}

func (r *Request) Updated() time.Time {
	return r.UpdatedValue
}
