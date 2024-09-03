package api

import (
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keysCreate() *rpc {
	return &rpc{
		Namespace: "keys",
		Action:    "create",
		Documentation: RPCDocumentation{
			Description: "Create a key",
			Parameters: driplimit.KeyCreatePayload{
				KSID:      "ks_abc",
				ExpiresIn: driplimit.Milliseconds{Duration: time.Minute * 5},
				Ratelimit: driplimit.RatelimitPayload{
					Limit:          5,
					RefillRate:     1,
					RefillInterval: driplimit.Milliseconds{Duration: time.Second},
				},
			},
			Response: driplimit.Key{
				KID:       "k_xyz",
				KSID:      "ks_abc",
				Token:     "demo_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Minute * 5),
				Ratelimit: &driplimit.Ratelimit{
					State: &driplimit.RatelimitState{
						LastRefilled: time.Now(),
						Remaining:    5,
					},
					Limit:          5,
					RefillRate:     1,
					RefillInterval: driplimit.Milliseconds{Duration: time.Second},
				},
			},
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.KeyCreatePayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}

			key, token, err := api.service.KeyCreate(c.Context(), *payload.WithServiceToken(token(c)))
			if err != nil {
				return err
			}

			key.Token = *token
			return c.JSON(key)
		},
	}
}
