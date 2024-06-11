package api

import (
	"fmt"
	"time"

	"github.com/i4n-co/driplimit"

	"github.com/gofiber/fiber/v2"
)

func (api *Server) serviceKeysList() *rpc {
	return &rpc{
		Namespace: "serviceKeys",
		Action:    "list",
		Documentation: RPCDocumentation{
			Description: "List all service keys",
			Parameters: driplimit.ServiceKeyListPayload{
				List: driplimit.ListPayload{
					Page:  1,
					Limit: 10,
				},
			},
			Response: driplimit.ServiceKeyList{
				List: driplimit.ListMetadata{
					Page:     1,
					Limit:    10,
					LastPage: 1,
				},
				ServiceKeys: []*driplimit.ServiceKey{
					{
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
			},
		},
		Handler: func(c *fiber.Ctx) (err error) {
			payload := new(driplimit.ServiceKeyListPayload)
			if err := c.BodyParser(payload); err != nil {
				return err
			}
			sklist, err := api.service.WithToken(token(c)).ServiceKeyList(c.Context(), *payload)
			if err != nil {
				return err
			}
			return c.JSON(sklist)
		},
	}
}
