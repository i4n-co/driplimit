package layouts

import (
	"html/template"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
)

type Pagination struct {
	Path    string
	Page    int
	Limit   int
	Current bool
}

func Paginate(c *fiber.Ctx, tpl string, list driplimit.ListMetadata) template.HTML {
	return template.HTML(render(10, c.Path(), list, func(w io.Writer, args Pagination) {
		if err := c.App().Config().Views.Render(w, tpl, args); err != nil {
			_, _ = io.Copy(w, strings.NewReader(err.Error()))
		}
	}))
}

func render(maxlinks int, path string, list driplimit.ListMetadata, render func(io.Writer, Pagination)) string {
	w := new(strings.Builder)
	numlinks := list.LastPage
	if numlinks > maxlinks {
		numlinks = maxlinks
	}
	for i := 1; i <= numlinks; i++ {
		if numlinks == maxlinks {
			pageNearLastPage := list.Page >= list.LastPage-maxlinks/2
			if pageNearLastPage {
				render(w, Pagination{
					Path:    path,
					Page:    list.LastPage - maxlinks + i,
					Limit:   list.Limit,
					Current: list.Page == list.LastPage-maxlinks+i,
				})
				continue
			}
			if list.Page >= maxlinks/2 {
				render(w, Pagination{
					Path:    path,
					Page:    i + list.Page - maxlinks/2,
					Limit:   list.Limit,
					Current: list.Page == i+list.Page-maxlinks/2,
				})
				continue
			}
		}
		render(w, Pagination{
			Page:    i,
			Limit:   list.Limit,
			Current: list.Page == i,
		})
	}
	return w.String()
}
