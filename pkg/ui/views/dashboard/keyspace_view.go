package dashboard

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

func KeyspaceView(service driplimit.Service, logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		keyspace, err := service.KeyspaceGet(c.Context(), driplimit.KeyspaceGetPayload{
			KSID: c.Params("ksid"),
		})
		if err != nil {
			return err
		}

		keylist, err := service.KeyList(c.Context(), driplimit.KeyListPayload{
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

			"KeyList": keylist,
		}, "layouts/page", "layouts/dashboard")
	}
}
