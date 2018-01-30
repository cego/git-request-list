package text

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cego/git-request-list/formatters"
)

// Table represents an ASCII table containing pull- and merge-requests.
type Table struct {
	arguments formatters.Arguments
}

func init() {
	factory := func(arguments formatters.Arguments) (formatters.Formatter, error) {
		t := Table{}
		t.arguments = arguments

		return &t, nil
	}

	formatters.RegisterFormatter("text", factory)
}

// String returns the ASCII string that t represents.
func (t *Table) String() string {
	rows := [][]string{{"Repository", "Name", "URL", "Created", "Updated"}}
	for _, r := range t.arguments.Requests {
		rows = append(rows, []string{
			r.Repository,
			r.Name,
			r.URL,
			r.Created.In(t.arguments.Timezone).Format(time.UnixDate),
			r.Updated.In(t.arguments.Timezone).Format(time.UnixDate),
		})
	}

	colWidths := map[int]int{}
	for _, row := range rows {
		for i, cell := range row {
			w, exists := colWidths[i]
			if !exists || utf8.RuneCountInString(cell) > w {
				colWidths[i] = utf8.RuneCountInString(cell)
			}
		}
	}

	result := ""
	for _, row := range rows {
		for i, cell := range row {
			result = result + cell + strings.Repeat(" ", colWidths[i]-utf8.RuneCountInString(cell))

			if i < len(row)-1 {
				result = result + " "
			}
		}

		result = result + "\n"
	}

	return result
}
