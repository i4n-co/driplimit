package api

import (
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keysGet() *rpc {
	return &rpc{
		Namespace: "keys",
		Action:    "get",
		Documentation: RPCDocumentation{
			Description: "Get a key",
			Parameters: driplimit.KeyGetPayload{
				KSID: "ks_abc",
				KID:  "k_xyz",
			},
			Response: driplimit.Key{
				KID:       "k_xyz",
				KSID:      "ks_abc",
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Minute * 5),
				Ratelimit: &driplimit.Ratelimit{
					State: &driplimit.RatelimitState{
						LastRefilled: time.Now(),
						Remaining:    4,
					},
					Limit:          5,
					RefillRate:     1,
					RefillInterval: driplimit.Milliseconds{Duration: time.Second},
				},
			},
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.KeyGetPayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}

			keyinfo, err := api.service.KeyGet(c.Context(), *payload.WithServiceToken(token(c)))
			if err != nil {
				return err
			}
			return c.JSON(keyinfo)
		},
	}
}
