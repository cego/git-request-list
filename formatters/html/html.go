package html

import (
	"bytes"
	"html/template"

	"github.com/cego/git-request-list/formatters"
)

// Table represents a HTML table containing pull- and merge-requests.
type Table struct {
	html string
}

const htmlTemplate = `
<html>
  <head>
    <meta charset="utf-8" /> 
    <link href="style.css" rel="stylesheet">
  </head>
  <body>
    <table>
      <tr>
        <th class="header-repository">Repository</th>
        <th class="header-name">Name</th>
        <th class="header-created">Created</th>
        <th class="header-updated">Updated</th>
      </tr>
      {{range .Requests}}
      <tr>
        <td class="item-repository">{{.Repository}}</td>
        <td class="item-name"><a target="_top" href="{{.URL}}">{{.Name}}</a></td>
        <td class="item-created">{{(.Created.In $.Timezone).Format "2006-01-02 15:04"}}</td>
        <td class="item-updated">{{(.Updated.In $.Timezone).Format "2006-01-02 15:04"}}</td>
      </tr>
      {{end}}
    </table>
  </body>
</html>
`

func init() {
	tmpl := template.Must(template.New("html").Parse(htmlTemplate))

	factory := func(arguments formatters.Arguments) (formatters.Formatter, error) {
		// Compile htmlTemplate into buff
		var buffer bytes.Buffer
		err := tmpl.Execute(&buffer, arguments)
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
