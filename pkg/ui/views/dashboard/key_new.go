package dashboard

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

type KeyNewForm struct {
	Limit             int    `form:"limit"`
	Expiration        string `form:"expiration"`
	Timezone          string `form:"tz"`
	OverrideRatelimit bool   `form:"override-ratelimit"`
	RefillRate        int    `form:"refill-rate"`
	RefillInterval    string `form:"refill-interval"`
}

func (form *KeyNewForm) KeyCreatePayload(ksid string) (*driplimit.KeyCreatePayload, error) {
	tz, err := time.LoadLocation(form.Timezone)
	if err != nil {
		return nil, layouts.NewError(400, "failed to load local timezone", err)
	}
	expiresAt, err := time.Parse("2006-01-02T15:04", form.Expiration)
	if err != nil {
		return nil, layouts.NewError(400, "the expiration date is invalid", err)
	}

	var ratelimit driplimit.RatelimitPayload
	if form.OverrideRatelimit {
		refillInterval, err := time.ParseDuration(form.RefillInterval)
		if err != nil {
			return nil, layouts.NewError(400, "the refill interval is invalid", err)
		}

		ratelimit = driplimit.RatelimitPayload{
			Limit:          int64(form.Limit),
			RefillRate:     int64(form.RefillRate),
			RefillInterval: driplimit.Milliseconds{Duration: refillInterval},
		}
	}

	return &driplimit.KeyCreatePayload{
		KSID:      ksid,
		ExpiresAt: expiresAt.In(tz),
		Ratelimit: ratelimit,
	}, nil
}

func KeyNewHXPost(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		keyspace, err := service(c).KeyspaceGet(c.Context(), driplimit.KeyspaceGetPayload{
			KSID: c.Params("ksid"),
		})
		if err != nil {
			return layouts.NewError(driplimit.HTTPCodeFromErr(err), "failed to get keyspace", err)
		}

		var key *driplimit.Key
		form := new(KeyNewForm)
		err = c.BodyParser(form)
		if err != nil {
			return layouts.NewError(400, "failed to parse request", err)
		}

		payload, err := form.KeyCreatePayload(c.Params("ksid"))
		if err != nil {
			return err
		}

		key, err = service(c).KeyCreate(c.Context(), *payload)
		if err != nil {
			return err
		}

		return c.Render("dashboard/key_new@created", fiber.Map{
			"Key":      key,
			"Keyspace": keyspace,
		})
	}
}

func KeyNew(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		keyspace, err := service(c).KeyspaceGet(c.Context(), driplimit.KeyspaceGetPayload{
			KSID: c.Params("ksid"),
		})
		if err != nil {
			return layouts.NewError(driplimit.HTTPCodeFromErr(err), "failed to get keyspace", err)
		}

		return c.Render("dashboard/key_new", fiber.Map{
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
					Name: "Create Key",
				},
			},
			"CSRFHiddenInput": c.Locals("csrf_template"),
			"Keyspace":        keyspace,
		}, "layouts/dashboard", "layouts/page")
	}
}
