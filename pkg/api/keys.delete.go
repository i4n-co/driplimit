package api

import (
	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keysDelete() *rpc {
	return &rpc{
		Namespace: "keys",
		Action:    "delete",
		Documentation: RPCDocumentation{
			Description: "Delete a key from a keyspace",
			Parameters: driplimit.KeyDeletePayload{
				KSID: "ks_abc",
				KID:  "k_xyz",
			},
			Response: nil,
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.KeyDeletePayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}

			err = api.service.KeyDelete(c.Context(), *payload.WithServiceToken(token(c)))
			if err != nil {
				return err
			}
			return c.SendStatus(fiber.StatusNoContent)
		},
	}
}
