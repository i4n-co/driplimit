package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/i4n-co/driplimit"
)

func (api *Server) serviceKeysCurrent(c *fiber.Ctx) (err error) {
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
}

func (api *Server) serviceKeysGet(c *fiber.Ctx) (err error) {
	payload := new(driplimit.ServiceKeyGetPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}
	sk, err := api.service.WithToken(token(c)).ServiceKeyGet(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(sk)
}

func (api *Server) serviceKeysList(c *fiber.Ctx) (err error) {
	payload := new(driplimit.ServiceKeyListPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}
	sklist, err := api.service.WithToken(token(c)).ServiceKeyList(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(sklist)
}

func (api *Server) serviceKeysDelete(c *fiber.Ctx) (err error) {
	payload := new(driplimit.ServiceKeyDeletePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}
	return api.service.WithToken(token(c)).ServiceKeyDelete(c.Context(), *payload)
}

func (api *Server) serviceKeysCreate(c *fiber.Ctx) (err error) {
	payload := new(driplimit.ServiceKeyCreatePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}
	sk, err := api.service.WithToken(token(c)).ServiceKeyCreate(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(sk)
}
