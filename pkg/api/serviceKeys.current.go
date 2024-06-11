package api

import (
	"fmt"
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) serviceKeysCurrent() *rpc {
	return &rpc{
		Namespace: "serviceKeys",
		Action:    "current",
		Documentation: RPCDocumentation{
			Description: "Get the current authenticated service key",
			Parameters:  nil,
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
			sk, err := api.service.WithToken(token(c)).ServiceKeyGet(c.Context(), driplimit.ServiceKeyGetPayload{
				Token: token(c),
			})
			if err != nil {
				if err == driplimit.ErrNotFound {
					return driplimit.ErrUnauthorized
				}
				return err
			}
			return c.JSON(sk)
		},
	}
}
