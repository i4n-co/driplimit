package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
)

// authenticate is a middleware that checks the presence of an Bearer token
func authenticate() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return driplimit.ErrUnauthorized
		}

		c.Locals("token", strings.TrimPrefix(auth, "Bearer "))
		return c.Next()
	}
}

// token returns the Bearer token from the context
func token(c *fiber.Ctx) string {
	return c.Locals("token").(string)
}
