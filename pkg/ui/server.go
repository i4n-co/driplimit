package ui

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/i4n-co/driplimit"
)

type Server struct {
	service driplimit.Service
	router  *fiber.App
	engine  *html.Engine
	logger  *slog.Logger
}

func New(service driplimit.Service, logger *slog.Logger) *Server {
	server := &Server{
		service: service,
		logger:  logger,
		engine:  loadEngine(),
	}
	server.router = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 server.engine,
	})

	server.router.Use("/statics", loadStaticsMiddleware())
	server.router.Get("/", server.home)

	return server
}

// Listen starts the UI server on the given address
func (ui *Server) Listen(addr string) error {
	ui.logger.Info("starting driplimit ui...", "addr", addr)
	return ui.router.Listen(addr)
}

// ShutdownWithContext shuts down the UI server gracefully
func (ui *Server) ShutdownWithContext(ctx context.Context) error {
	ui.logger.Info("shutting down driplimit ui...")
	return ui.router.ShutdownWithContext(ctx)
}

func (ui *Server) home(c *fiber.Ctx) error {
	return c.Render("dashboard", fiber.Map{
		"Title": "Driplimit Dashboard",
	}, "layouts/page")
}
