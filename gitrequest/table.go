package gitrequest

import (
	"strings"
	"time"
)

// Table represents an ASCII table containing pull- and merge-requests.
type Table struct{}

// String returns the ASCII table t containing the given requests.
func (t *Table) String(requests ...Request) string {
	rows := [][]string{{"Repository", "Name", "State", "URL", "Created", "Updated"}}
	for _, r := range requests {
		rows = append(rows, []string{r.Repository(), r.Name(), r.State(), r.URL(), r.Created().Format(time.UnixDate), r.Updated().Format(time.UnixDate)})
	}

	colWidths := map[int]int{}
	for _, row := range rows {
		for i, cell := range row {
			w, exists := colWidths[i]
			if !exists || strLen(cell) > w {
				colWidths[i] = strLen(cell)
			}
		}
	}

	result := ""
	for _, row := range rows {
		for i, cell := range row {
			result = result + cell + strings.Repeat(" ", colWidths[i]-strLen(cell))

			if i < len(row)-1 {
				result = result + " "
			}
		}

		result = result + "\n"
	}

	return result
}

// strLen returns the rune-length of s.
func strLen(s string) int {
	return len([]rune(s))
}
