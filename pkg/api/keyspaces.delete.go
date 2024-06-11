package api

import (
	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keyspacesDelete() *rpc {
	return &rpc{
		Namespace: "keyspaces",
		Action:    "delete",
		Documentation: RPCDocumentation{
			Description: "Delete a keyspace",
			Parameters: driplimit.KeyspaceDeletePayload{
				KSID: "ks_abc",
			},
			Response: nil,
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.KeyspaceDeletePayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}
			err = api.service.WithToken(token(c)).KeyspaceDelete(c.Context(), *payload)
			if err != nil {
				return err
			}
			return c.SendStatus(fiber.StatusNoContent)
		},
	}
}
