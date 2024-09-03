package api

import (
	"fmt"
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) serviceKeysGet() *rpc {
	return &rpc{
		Namespace: "serviceKeys",
		Action:    "get",
		Documentation: RPCDocumentation{
			Description: "Get the service key by ID or by token",
			Parameters: driplimit.ServiceKeyGetPayload{
				SKID:  "sk_uvw",
				Token: "",
			},
			Response: driplimit.ServiceKey{
				SKID:        "sk_uvw",
				Description: fmt.Sprintf("cli generated admin service key at %s", time.Now().Format(time.RFC3339)),
				Admin:       true,
				KeyspacesPolicies: map[string]driplimit.Policy{
					"ks_abc": {
						Read:  true,
						Write: true,
					},
				},
				CreatedAt: time.Now(),
			},
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.ServiceKeyGetPayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}
			sk, err := api.service.ServiceKeyGet(c.Context(), *payload.WithServiceToken(token(c)))
			if err != nil {
				return err
			}
			return c.JSON(sk)
		},
	}
}
