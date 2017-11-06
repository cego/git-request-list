package formatters

import (
	"sort"

	"github.com/cego/git-request-list/request"
)

// byRepository implements sort.Interface for []request.Request based on Repository.
type byRepository []request.Request

func (rs byRepository) Len() int           { return len(rs) }
func (rs byRepository) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs byRepository) Less(i, j int) bool { return rs[i].Repository < rs[j].Repository }

// byName implements sort.Interface for []request.Request based on Name.
type byName []request.Request

func (rs byName) Len() int           { return len(rs) }
func (rs byName) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs byName) Less(i, j int) bool { return rs[i].Name < rs[j].Name }

// byURL implements sort.Interface for []request.Request based on URL.
type byURL []request.Request

func (rs byURL) Len() int           { return len(rs) }
func (rs byURL) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs byURL) Less(i, j int) bool { return rs[i].URL < rs[j].URL }

// byCreated implements sort.Interface for []request.Request based on Created.
type byCreated []request.Request

func (rs byCreated) Len() int           { return len(rs) }
func (rs byCreated) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs byCreated) Less(i, j int) bool { return rs[i].Created.Before(rs[j].Created) }

// byUpdated implements sort.Interface for []request.Request based on Updated.
type byUpdated []request.Request

func (rs byUpdated) Len() int           { return len(rs) }
func (rs byUpdated) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs byUpdated) Less(i, j int) bool { return rs[i].Updated.Before(rs[j].Updated) }

// Sort sorts requests by the property named by
func Sort(requests []request.Request, by string) {
	switch by {
	case "name":
		sort.Sort(byName(requests))
		break
	case "url":
		sort.Sort(byURL(requests))
		break
	case "created":
		sort.Sort(byCreated(requests))
		break
	case "updated":
		sort.Sort(byUpdated(requests))
		break
	case "repository":
		sort.Sort(byRepository(requests))
	default:
		break
	}
}
