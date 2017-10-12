package output

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
