package gitrequest

import (
	"time"
)

type Request interface {
	Repository() string
	Name() string
	State() string
	URL() string
	Created() time.Time
	Updated() time.Time
}

// ByRepository implements sort.Interface for []Request based on Repository.
type ByRepository []Request

func (rs ByRepository) Len() int           { return len(rs) }
func (rs ByRepository) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByRepository) Less(i, j int) bool { return rs[i].Repository() < rs[j].Repository() }

// ByName implements sort.Interface for []Request based on Name.
type ByName []Request

func (rs ByName) Len() int           { return len(rs) }
func (rs ByName) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByName) Less(i, j int) bool { return rs[i].Name() < rs[j].Name() }

// ByState implements sort.Interface for []Request based on State.
type ByState []Request

func (rs ByState) Len() int           { return len(rs) }
func (rs ByState) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByState) Less(i, j int) bool { return rs[i].State() < rs[j].State() }

// ByURL implements sort.Interface for []Request based on URL.
type ByURL []Request

func (rs ByURL) Len() int           { return len(rs) }
func (rs ByURL) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByURL) Less(i, j int) bool { return rs[i].URL() < rs[j].URL() }

// ByCreated implements sort.Interface for []Request based on Created.
type ByCreated []Request

func (rs ByCreated) Len() int           { return len(rs) }
func (rs ByCreated) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByCreated) Less(i, j int) bool { return rs[i].Created().Before(rs[j].Created()) }

// ByUpdated implements sort.Interface for []Request based on Updated.
type ByUpdated []Request

func (rs ByUpdated) Len() int           { return len(rs) }
func (rs ByUpdated) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs ByUpdated) Less(i, j int) bool { return rs[i].Updated().Before(rs[j].Updated()) }
