package api

import (
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keysList() *rpc {
	return &rpc{
		Namespace: "keys",
		Action:    "list",
		Documentation: RPCDocumentation{
			Description: "List keys",
			Parameters: driplimit.KeyListPayload{
				KSID: "ks_abc",
				List: driplimit.ListPayload{
					Page:  1,
					Limit: 10,
				},
			},
			Response: driplimit.KeyList{
				List: driplimit.ListMetadata{
					Page:  1,
					Limit: 10,
					LastPage: 1,
				},
				Keys: []*driplimit.Key{
					{
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
			},
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.KeyListPayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}

			keylist, err := api.service.WithToken(token(c)).KeyList(c.Context(), *payload)
			if err != nil {
				return err
			}
			return c.JSON(keylist)
		},
	}
}
