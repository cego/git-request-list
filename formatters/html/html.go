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
  <head>
    <link href="style.css" rel="stylesheet">
  </head>
  <body>
    <table>
      <tr>
        <th class="header-repository">Repository</th>
        <th class="header-name">Name</th>
        <th class="header-url">URL</th>
        <th class="header-created">Created</th>
        <th class="header-updated">Updated</th>
      </tr>
      {{range .Requests}}
      <tr>
        <td class="item-repository">{{.Repository}}</td>
        <td class="item-name">{{.Name}}</td>
        <td class="item-url">{{.URL}}</td>
        <td class="item-created">{{.Created}}</td>
        <td class="item-updated">{{.Updated}}</td>
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
