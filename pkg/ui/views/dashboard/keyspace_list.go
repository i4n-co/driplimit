package dashboard

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

func KeyspaceList(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		list, err := service(c).KeyspaceList(c.Context(), driplimit.KeyspaceListPayload{})
		if err != nil {
			return err
		}
		return c.Render("dashboard/keyspace_list", fiber.Map{
			"Title": "Driplimit",
			"Breadcrumbs": []layouts.Breadcrumb{
				{
					Href: "/",
					Name: "Home",
				},
			},
			"List": list,
		}, "layouts/dashboard", "layouts/page")
	}
}
