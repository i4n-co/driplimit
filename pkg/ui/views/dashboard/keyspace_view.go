package dashboard

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

func KeyspaceView(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		keyspace, err := service(c).KeyspaceGet(c.Context(), driplimit.KeyspaceGetPayload{
			KSID: c.Params("ksid"),
		})
		if err != nil {
			return err
		}

		listpayload := new(driplimit.ListPayload)
		if err := c.QueryParser(listpayload); err != nil {
			return layouts.NewError(400, "Failed to parse url parameters")
		}

		keylist, err := service(c).KeyList(c.Context(), driplimit.KeyListPayload{
			List: *listpayload,
			KSID: c.Params("ksid"),
		})
		if err != nil {
			return err
		}
		return c.Render("dashboard/keyspace_view", fiber.Map{
			"Title": "Driplimit",
			"Breadcrumbs": []layouts.Breadcrumb{
				{
					Href: "/",
					Name: "Home",
				},
				{
					Href: "/keyspaces",
					Name: "Keyspaces",
				},
				{
					Href: c.Path(),
					Name: keyspace.Name,
				},
			},
			"Keyspace": keyspace,

			"KeyList":    keylist,
			"Pagination": layouts.Paginate(c, "dashboard/keys_pagination", keylist.List),
		}, "layouts/dashboard", "layouts/page")
	}
}
