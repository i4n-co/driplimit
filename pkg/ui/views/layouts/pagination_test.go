package layouts

import (
	"html/template"
	"io"
	"testing"

	"github.com/i4n-co/driplimit"
	"github.com/stretchr/testify/assert"
)

func TestPaginationTemplate(t *testing.T) {

	linkgen := func(w io.Writer, arg Pagination) {
		tpl := `{{.Page}}{{ if .Current }}*{{ end }}-`
		tmpl, err := template.New("").Parse(tpl)
		assert.NoError(t, err)
		tmpl.Execute(w, arg)
	}

	assert.Equal(t, "", render(10, "/", driplimit.ListMetadata{}, linkgen))
	assert.Equal(t, "1*-", render(10, "/", driplimit.ListMetadata{Page: 1, Limit: 10, LastPage: 1}, linkgen))
	assert.Equal(t, "1*-2-3-4-5-6-7-8-", render(10, "/", driplimit.ListMetadata{Page: 1, Limit: 10, LastPage: 8}, linkgen))
	assert.Equal(t, "1*-2-3-4-5-6-7-8-9-10-", render(10, "/", driplimit.ListMetadata{Page: 1, Limit: 10, LastPage: 11}, linkgen))
	assert.Equal(t, "6-7-8-9-10*-11-12-13-14-15-", render(10, "/", driplimit.ListMetadata{Page: 10, Limit: 10, LastPage: 20}, linkgen))
	assert.Equal(t, "11-12-13-14-15-16-17-18*-19-20-", render(10, "/", driplimit.ListMetadata{Page: 18, Limit: 10, LastPage: 20}, linkgen))
}
