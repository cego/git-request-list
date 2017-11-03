package html

import (
	"bytes"
	"text/template"

	"github.com/cego/git-request-list/formatters"
	"github.com/cego/git-request-list/request"
)

// Table represents a HTML table containing pull- and merge-requests.
type Table struct {
	html string
}

const htmlTemplate = `
<html>
  <head></head>
  <body>
    <table>
      <tr>
        <th>Repository</th>
        <th>Name</th>
        <th>URL</th>
        <th>Created</th>
        <th>Updated</th>
      </tr>
      {{range .Requests}}
      <tr>
        <td>{{.Repository}}</td>
        <td>{{.Name}}</td>
        <td>{{.URL}}</td>
        <td>{{.Created}}</td>
        <td>{{.Updated}}</td>
      </tr>
      {{end}}
    </table>
  </body>
</html>
`

func init() {
	tmpl := template.Must(template.New("html").Parse(htmlTemplate))

	factory := func(requests []request.Request) (formatters.Formatter, error) {
		// Compile htmlTemplate into buff
		data := struct{ Requests []request.Request }{Requests: requests}
		var buffer bytes.Buffer
		err := tmpl.Execute(&buffer, data)
		if err != nil {
			return nil, err
		}

		t := Table{}
		t.html = buffer.String()

		return &t, nil
	}

	formatters.RegisterFormatter("html", factory)
}

// String returns the HTML string that t represents.
func (t *Table) String() string {
	return t.html
}
