package api

import (
	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) serviceKeysDelete() *rpc {
	return &rpc{
		Namespace: "serviceKeys",
		Action:    "delete",
		Documentation: RPCDocumentation{
			Description: "Delete a service key",
			Parameters: driplimit.ServiceKeyDeletePayload{
				SKID: "sk_uvw",
			},
			Response: nil,
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.ServiceKeyDeletePayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}
			return api.service.WithToken(token(c)).ServiceKeyDelete(c.Context(), *payload)
		},
	}
}
