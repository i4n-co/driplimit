package dashboard

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/client"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

func Login(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("dashboard/login", fiber.Map{
			"CSRFHiddenInput": c.Locals("csrf_template"),
			"Title":           "Driplimit",
		}, "layouts/dashboard", "layouts/page")
	}
}

func Logout(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.ClearCookie("service_key_token")
		c.Set("HX-Redirect", "/login")
		return c.Redirect("/login")
	}
}

type LoginForm struct {
	ServiceKeyToken string `form:"service_key_token"`
}

func LoginHXPost(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		form := new(LoginForm)
		err := c.BodyParser(form)
		if err != nil {
			return layouts.NewError(400, "failed to parse request", err)
		}
		cookie := new(fiber.Cookie)
		cookie.Name = "service_key_token"
		cookie.Value = form.ServiceKeyToken
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.Cookie(cookie)
		c.Set("HX-Redirect", "/keyspaces")
		return nil
	}
}

func AuthMiddleware(httpcli *client.HTTP) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		for _, allowedPrefix := range []string{"/sse", "/statics", "/login"} {
			if strings.HasPrefix(c.Path(), allowedPrefix) {
				return c.Next()
			}
		}

		token := c.Cookies("service_key_token")
		if token == "" {
			c.Set("HX-Redirect", "/login")
			return c.Redirect("/login")
		}
		httpcli = httpcli.WithServiceToken(token)
		sk, err := httpcli.ServiceKeyCurrent(c.Context())
		if err != nil {
			if errors.Is(err, driplimit.ErrUnauthorized) {
				c.ClearCookie("service_key_token")
				return layouts.NewError(403, "Unauthorized", err)
			}
			return layouts.NewError(500, "Something went wrong", err)
		}
		c.Locals("sk", sk)
		c.Locals("service", httpcli)
		return c.Next()
	}
}

func service(c *fiber.Ctx) *client.HTTP {
	return c.Locals("service").(*client.HTTP)
}

func ServiceKey(c *fiber.Ctx) *driplimit.ServiceKey {
	return c.Locals("sk").(*driplimit.ServiceKey)
}
