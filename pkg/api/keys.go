package api

import (
	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keysCreate(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeyCreatePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	key, token, err := api.service.WithToken(token(c)).KeyCreate(c.Context(), *payload)
	if err != nil {
		return err
	}

	key.Token = *token
	return c.JSON(key)
}

func (api *Server) keysCheck(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeysCheckPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	keyinfo, err := api.service.WithToken(token(c)).KeyCheck(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(keyinfo)
}

func (api *Server) keysGet(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeyGetPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	keyinfo, err := api.service.WithToken(token(c)).KeyGet(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(keyinfo)
}

func (api *Server) keysList(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeyListPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	keylist, err := api.service.WithToken(token(c)).KeyList(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(keylist)
}

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
