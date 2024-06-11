package api

import (
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keyspacesList() *rpc {
	return &rpc{
		Namespace: "keyspaces",
		Action:    "list",
		Documentation: RPCDocumentation{
			Description: "Get keyspace by ID",
			Parameters: driplimit.KeyspaceListPayload{
				List: driplimit.ListPayload{
					Page:  1,
					Limit: 10,
				},
			},
			Response: driplimit.KeyspaceList{
				List: driplimit.ListMetadata{
					Page:     1,
					Limit:    10,
					LastPage: 1,
				},
				Keyspaces: []*driplimit.Keyspace{
					{
						KSID:       "ks_abc",
						Name:       "demo.yourapi.com (env: production)",
						KeysPrefix: "demo_",
						Ratelimit: &driplimit.Ratelimit{
							Limit:          100,
							RefillRate:     1,
							RefillInterval: driplimit.Milliseconds{Duration: time.Second},
						},
					},
				},
			},
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.KeyspaceListPayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}
			keyspace, err := api.service.WithToken(token(c)).KeyspaceList(c.Context(), *payload)
			if err != nil {
				return err
			}
			return c.JSON(keyspace)
		},
	}
}
