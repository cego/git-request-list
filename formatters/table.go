package formatters

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cego/git-request-list/providers"
)

// Table represents an ASCII table containing pull- and merge-requests.
type Table struct{}

// String returns the ASCII table t containing the given requests.
func (t *Table) String(requests ...providers.Request) string {
	rows := [][]string{{"Repository", "Name", "URL", "Created", "Updated"}}
	for _, r := range requests {
		rows = append(rows, []string{r.Repository, r.Name, r.URL, r.Created.Format(time.UnixDate), r.Updated.Format(time.UnixDate)})
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
