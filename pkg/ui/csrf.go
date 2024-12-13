package ui

import (
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

// csrfTokenTemplateMiddleware injects an html template in fiber context to
// be used in html forms
func csrfTokenTemplateMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals("csrf_template", template.HTML(
			fmt.Sprintf(`<input type="hidden" name="csrf_token" value="%s" />`, c.Locals("csrf_token")),
		))
		return c.Next()
	}
}
