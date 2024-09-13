package ui

import (
	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

func HXRequest(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "true" {
		return c.Next()
	}
	return layouts.NewError(400, "this route only accept partial requests")
}
