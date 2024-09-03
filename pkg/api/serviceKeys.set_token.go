package api

import (
	"net/http"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) serviceKeysSetToken() *rpc {
	return &rpc{
		Namespace: "serviceKeys",
		Action:    "set_token",
		Documentation: RPCDocumentation{
			Description: "Set a new token for a service key",
			Parameters: driplimit.ServiceKeySetTokenPayload{
				SKID:  "ks_abc",
				Token: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			Code: http.StatusNoContent,
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.ServiceKeySetTokenPayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}
			err = api.service.ServiceKeySetToken(c.Context(), *payload.WithServiceToken(token(c)))
			if err != nil {
				return err
			}
			return c.SendStatus(http.StatusNoContent)
		},
	}
}
