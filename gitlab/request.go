package gitlab

import (
	"time"
)

// Request implements github.com/cego/git-request-list/gitrequest.Request and serves as Unmarshal target type when
// reading Gitlab API responses.
type Request struct {
	RepositoryValue string
	NameValue       string    `json:"title"`
	StateValue      string    `json:"state"`
	URLValue        string    `json:"web_url"`
	CreatedValue    time.Time `json:"created_at"`
	UpdatedValue    time.Time `json:"updated_at"`
	WIP             bool      `json:"work_in_progress"`
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
