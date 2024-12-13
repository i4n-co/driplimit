package dashboard

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

func KeyDelete(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		keyspace, err := service(c).KeyspaceGet(c.Context(), driplimit.KeyspaceGetPayload{
			KSID: c.Params("ksid"),
		})
		if err != nil {
			return layouts.NewError(driplimit.HTTPCodeFromErr(err), "failed to get keyspace", err)
		}
		err = service(c).KeyDelete(c.Context(), driplimit.KeyDeletePayload{
			KSID: keyspace.KSID,
			KID:  c.Params("kid"),
		})
		if err != nil {
			return layouts.NewError(driplimit.HTTPCodeFromErr(err), "failed to delete key", err)
		}

		keyspaceURL := fmt.Sprintf("/keyspaces/%s", keyspace.KSID)
		if c.Get("HX-Request") == "true" {
			c.Set("HX-Redirect", keyspaceURL)
			return c.SendStatus(200)
		}
		return c.Redirect(keyspaceURL)
	}
}
