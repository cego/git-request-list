package gitrequest

import (
	"strings"
	"time"
)

type Table struct {
	rows [][]string
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) Add(r Request) {
	t.rows = append(t.rows, []string{r.Repository(), r.Name(), r.State(), r.URL(), r.Created().Format(time.UnixDate), r.Updated().Format(time.UnixDate)})
}

func (t *Table) String() string {

	rows := append(
		[][]string{{"Repository", "Name", "State", "URL", "Created", "Updated"}},
		t.rows...,
	)

	result := ""

	colWidths := map[int]int{}
	for _, row := range rows {
		for i, cell := range row {
			w, exists := colWidths[i]
			if !exists || strLen(cell) > w {
				colWidths[i] = strLen(cell)
			}
		}
	}

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

func strLen(s string) int {
	return len([]rune(s))
}
