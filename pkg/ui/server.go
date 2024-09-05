package ui

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/i4n-co/driplimit"
)

type Server struct {
	service driplimit.Service
	router  *fiber.App
	engine  *html.Engine
	logger  *slog.Logger
	events  chan string
}

func New(service driplimit.Service, logger *slog.Logger) *Server {
	server := &Server{
		service: service,
		logger:  logger,
		engine:  loadEngine(),
		events:  make(chan string, 100),
	}
	server.router = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 server.engine,
	})

	server.router.Use("/statics", loadStaticsMiddleware())
	server.router.Get("/keyspaces", server.keyspaceList)
	server.router.Get("/keyspaces/:id", server.keyspaceView)
	server.router.Get("/keyspaces/:id/key_new", server.keyCreate)
	server.router.Get("/sse", server.sse)

	return server
}

// Listen starts the UI server on the given address
func (ui *Server) Listen(addr string) error {
	ui.events <- "start"
	ui.logger.Info("starting driplimit ui...", "addr", addr)
	return ui.router.Listen(addr)
}

// ShutdownWithContext shuts down the UI server gracefully
func (ui *Server) ShutdownWithContext(ctx context.Context) error {
	ui.logger.Info("shutting down driplimit ui...")
	return ui.router.ShutdownWithContext(ctx)
}

type Breadcrumb struct {
	Name string
	Href string
}

func (ui *Server) keyspaceList(c *fiber.Ctx) error {
	return c.Render("dashboard/keyspace_list", fiber.Map{
		"Title": "Driplimit",
		"Breadcrumbs": []Breadcrumb{
			{
				Href: "/",
				Name: "Home",
			},
			{
				Href: c.Path(),
				Name: "Keyspaces",
			},
		},
		"Keyspaces": []driplimit.Keyspace{
			{
				KSID: "id1",
				Name: "data.bordeaux-port.fr",
			},
			{
				KSID: "id2",
				Name: "api.weather.com",
			},
			{
				KSID: "id3",
				Name: "dev-data.bordeaux-port.fr",
			},
			{
				KSID: "id4",
				Name: "websocket.pokemon.co",
			},
		},
	}, "layouts/page", "layouts/dashboard")
}

func (ui *Server) keyspaceView(c *fiber.Ctx) error {
	return c.Render("dashboard/keyspace_view", fiber.Map{
		"Title": "Driplimit",
		"Breadcrumbs": []Breadcrumb{
			{
				Href: "/",
				Name: "Home",
			},
			{
				Href: "/keyspaces",
				Name: "Keyspaces",
			},
			{
				Href: c.Path(),
				Name: "data.jumeaux-numeriques.fr",
			},
		},
		"Keyspace": driplimit.Keyspace{
			KSID:       "id1",
			Name:       "data.jumeaux-numeriques.fr",
			KeysPrefix: "bdx_",
			Ratelimit: &driplimit.Ratelimit{
				Limit:      10,
				RefillRate: 1,
				RefillInterval: driplimit.Milliseconds{
					Duration: 1000 * time.Millisecond,
				},
			},
		},

		"KeyList": driplimit.KeyList{
			List: driplimit.ListMetadata{
				Page:     1,
				Limit:    10,
				LastPage: 5,
			},
			Keys: []*driplimit.Key{
				{
					KID:       "skjsksjd",
					KSID:      "id1",
					LastUsed:  time.Now().Add(-time.Hour),
					ExpiresAt: time.Now().Add(time.Hour * 10000),
					CreatedAt: time.Now().Add(time.Second * -19863723),
					Ratelimit: &driplimit.Ratelimit{
						State: &driplimit.RatelimitState{
							Remaining: 3,
						},
						Limit:      10,
						RefillRate: 1,
						RefillInterval: driplimit.Milliseconds{
							Duration: 1000 * time.Millisecond,
						},
					},
				},
				{
					KID:       "pajdcucx",
					KSID:      "id1",
					LastUsed:  time.Now().Add(-time.Hour),
					ExpiresAt: time.Now().Add(time.Hour * 10000),
					CreatedAt: time.Now().Add(time.Second * -19863723),
					Ratelimit: &driplimit.Ratelimit{
						State: &driplimit.RatelimitState{
							Remaining: 3,
						},
						Limit:      10,
						RefillRate: 1,
						RefillInterval: driplimit.Milliseconds{
							Duration: 1000 * time.Millisecond,
						},
					},
				},
				{
					KID:       "spandcnc",
					KSID:      "id1",
					LastUsed:  time.Now().Add(-time.Hour),
					ExpiresAt: time.Now().Add(time.Hour * 10000),
					CreatedAt: time.Now().Add(time.Second * -19863723),
					Ratelimit: &driplimit.Ratelimit{
						State: &driplimit.RatelimitState{
							Remaining: 3,
						},
						Limit:      10,
						RefillRate: 1,
						RefillInterval: driplimit.Milliseconds{
							Duration: 1000 * time.Millisecond,
						},
					},
				},
				{
					KID:       "oxlnzksh",
					KSID:      "id1",
					LastUsed:  time.Now().Add(-time.Hour),
					ExpiresAt: time.Now().Add(time.Hour * 10000),
					CreatedAt: time.Now().Add(time.Second * -19863723),
					Ratelimit: &driplimit.Ratelimit{
						State: &driplimit.RatelimitState{
							Remaining: 3,
						},
						Limit:      10,
						RefillRate: 1,
						RefillInterval: driplimit.Milliseconds{
							Duration: 1000 * time.Millisecond,
						},
					},
				},
				{
					KID:       "vparxfud",
					KSID:      "id1",
					LastUsed:  time.Now().Add(-time.Hour),
					ExpiresAt: time.Now().Add(time.Hour * 10000),
					CreatedAt: time.Now().Add(time.Second * -19863723),
					Ratelimit: &driplimit.Ratelimit{
						State: &driplimit.RatelimitState{
							Remaining: 3,
						},
						Limit:      10,
						RefillRate: 1,
						RefillInterval: driplimit.Milliseconds{
							Duration: 1000 * time.Millisecond,
						},
					},
				},
				{
					KID:       "skpqnwgd",
					KSID:      "id1",
					LastUsed:  time.Now().Add(-time.Hour),
					ExpiresAt: time.Now().Add(time.Hour * 10000),
					CreatedAt: time.Now().Add(time.Second * -19863723),
					Ratelimit: &driplimit.Ratelimit{
						State: &driplimit.RatelimitState{
							Remaining: 3,
						},
						Limit:      10,
						RefillRate: 1,
						RefillInterval: driplimit.Milliseconds{
							Duration: 1000 * time.Millisecond,
						},
					},
				},
			},
		},
	}, "layouts/page", "layouts/dashboard")
}

func (ui *Server) keyCreate(c *fiber.Ctx) error {
	return c.Render("dashboard/key_new", fiber.Map{
		"Title": "Driplimit",
		"Breadcrumbs": []Breadcrumb{
			{
				Href: "/",
				Name: "Home",
			},
			{
				Href: "/keyspaces",
				Name: "Keyspaces",
			},
			{
				Href: "/keyspaces/id1",
				Name: "data.jumeaux-numeriques.fr",
			},
			{
				Href: c.Path(),
				Name: "Create Key",
			},
		},
		"Keyspace": driplimit.Keyspace{
			KSID:       "id1",
			Name:       "data.jumeaux-numeriques.fr",
			KeysPrefix: "bdx_",
			Ratelimit: &driplimit.Ratelimit{
				Limit:      10,
				RefillRate: 1,
				RefillInterval: driplimit.Milliseconds{
					Duration: 1000 * time.Millisecond,
				},
			},
		},
	}, "layouts/page", "layouts/dashboard")
}
