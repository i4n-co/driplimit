package dashboard

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

func KeyView(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		keyspace, err := service(c).KeyspaceGet(c.Context(), driplimit.KeyspaceGetPayload{
			KSID: c.Params("ksid"),
		})
		if err != nil {
			return layouts.NewError(driplimit.HTTPCodeFromErr(err), "failed to get keyspace", err)
		}
		key, err := service(c).KeyGet(c.Context(), driplimit.KeyGetPayload{
			KSID: keyspace.KSID,
			KID:  c.Params("kid"),
		})
		if err != nil {
			return layouts.NewError(driplimit.HTTPCodeFromErr(err), "failed to get key", err)
		}

		return c.Render("dashboard/key_view", fiber.Map{
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
					Href: fmt.Sprintf("/keyspaces/%s", keyspace.KSID),
					Name: keyspace.Name,
				},
				{
					Href: c.Path(),
					Name: key.KID,
				},
			},
			"CSRFToken": c.Locals("csrf_token"),
			"Keyspace":  keyspace,
			"Key":       key,
		}, "layouts/dashboard", "layouts/page")
	}
}
