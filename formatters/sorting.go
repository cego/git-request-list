package formatters

import (
	"github.com/cego/git-request-list/request"
)

// ByRepository implements sort.Interface for []request.Request based on Repository.
type ByRepository []request.Request

func (rs ByRepository) Len() int           { return len(rs) }
func (rs ByRepository) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByRepository) Less(i, j int) bool { return rs[i].Repository < rs[j].Repository }

// ByName implements sort.Interface for []request.Request based on Name.
type ByName []request.Request

func (rs ByName) Len() int           { return len(rs) }
func (rs ByName) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByName) Less(i, j int) bool { return rs[i].Name < rs[j].Name }

// ByURL implements sort.Interface for []request.Request based on URL.
type ByURL []request.Request

func (rs ByURL) Len() int           { return len(rs) }
func (rs ByURL) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByURL) Less(i, j int) bool { return rs[i].URL < rs[j].URL }

// ByCreated implements sort.Interface for []request.Request based on Created.
type ByCreated []request.Request

func (rs ByCreated) Len() int           { return len(rs) }
func (rs ByCreated) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByCreated) Less(i, j int) bool { return rs[i].Created.Before(rs[j].Created) }

// ByUpdated implements sort.Interface for []request.Request based on Updated.
type ByUpdated []request.Request

func (rs ByUpdated) Len() int           { return len(rs) }
func (rs ByUpdated) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByUpdated) Less(i, j int) bool { return rs[i].Updated.Before(rs[j].Updated) }
