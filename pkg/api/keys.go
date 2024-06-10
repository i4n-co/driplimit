package api

import (
	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keysDelete(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeyDeletePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	err = api.service.WithToken(token(c)).KeyDelete(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
