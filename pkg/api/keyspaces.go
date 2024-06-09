package api

import (
	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) keyspacesGet(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeyspaceGetPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}
	keyspace, err := api.service.WithToken(token(c)).KeyspaceGet(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(keyspace)
}

func (api *Server) keyspacesCreate(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeyspaceCreatePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}
	keyspace, err := api.service.WithToken(token(c)).KeyspaceCreate(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(keyspace)
}

func (api *Server) keyspacesList(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeyspaceListPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}
	keyspace, err := api.service.WithToken(token(c)).KeyspaceList(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.JSON(keyspace)
}

func (api *Server) keyspacesDelete(c *fiber.Ctx) (err error) {
	payload := new(driplimit.KeyspaceDeletePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}
	err = api.service.WithToken(token(c)).KeyspaceDelete(c.Context(), *payload)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
