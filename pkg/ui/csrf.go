package ui

import (
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

const CSRFTokenContextKey ContextKey = "csrf_context_key"

// csrfTokenTemplateMiddleware injects an html template in fiber context to
// be used in html forms
func csrfTokenTemplateMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals("csrf_template", template.HTML(
			fmt.Sprintf(`<input type="hidden" name="csrf_token" value="%s" />`, c.Locals(CSRFTokenContextKey)),
		))
		return c.Next()
	}
}
