package api

import (
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keyspacesCreate() *rpc {
	return &rpc{
		Namespace: "keyspaces",
		Action:    "create",
		Documentation: RPCDocumentation{
			Description: "Create a new keyspace",
			Parameters: driplimit.KeyspaceCreatePayload{
				Name:       "demo.yourapi.com (env: production)",
				KeysPrefix: "demo_",
				Ratelimit: driplimit.RatelimitPayload{
					Limit:          100,
					RefillRate:     1,
					RefillInterval: driplimit.Milliseconds{Duration: time.Second},
				},
			},
			Response: driplimit.Keyspace{
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
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.KeyspaceCreatePayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}
			keyspace, err := api.service.WithToken(token(c)).KeyspaceCreate(c.Context(), *payload)
			if err != nil {
				return err
			}
			return c.JSON(keyspace)
		},
	}
}
