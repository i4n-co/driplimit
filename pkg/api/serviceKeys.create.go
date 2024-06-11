package api

import (
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) serviceKeysCreate() *rpc {
	return &rpc{
		Namespace: "serviceKeys",
		Action:    "create",
		Documentation: RPCDocumentation{
			Description: "Get the service key by ID or by token",
			Parameters: driplimit.ServiceKeyCreatePayload{
				Description: "api generated non admin service key",
				Admin:       false,
				KeyspacesPolicies: map[string]driplimit.Policy{
					"ks_abc": {
						Read:  true,
						Write: false,
					},
				},
			},
			Response: driplimit.ServiceKey{
				SKID:        "sk_uvw",
				Description: "api generated non admin service key",
				Admin:       true,
				KeyspacesPolicies: map[string]driplimit.Policy{
					"ks_abc": {
						Read:  true,
						Write: false,
					},
				},
				CreatedAt: time.Now(),
			},
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.ServiceKeyCreatePayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}
			sk, err := api.service.WithToken(token(c)).ServiceKeyCreate(c.Context(), *payload)
			if err != nil {
				return err
			}
			return c.JSON(sk)
		},
	}
}
