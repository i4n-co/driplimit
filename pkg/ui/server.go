package ui

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/ui/views/dashboard"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

type ContextKey string

type Server struct {
	service driplimit.Service
	router  *fiber.App
	logger  *slog.Logger
	events  chan string
}

func New(service driplimit.Service, logger *slog.Logger) *Server {
	server := &Server{
		service: service,
		logger:  logger,
		events:  make(chan string, 100),
	}
	server.router = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 NewTemplateEngine(),
		ErrorHandler:          server.errorHandler,
	})
	server.router.Use(encryptcookie.New(encryptcookie.Config{
		Key:    "secret-thirty-2-character-string",
		Except: []string{csrf.ConfigDefault.CookieName}, // exclude CSRF cookie
	}))
	server.router.Use(csrf.New(csrf.Config{
		ContextKey: CSRFTokenContextKey,
		Extractor: func(c *fiber.Ctx) (string, error) {
			type form struct {
				Token string `form:"csrf_token"`
			}
			body := new(form)
			err := c.BodyParser(body)
			if err != nil {
				server.logger.Error("body parser", "err", err)

				return "", fmt.Errorf("failed to parse body: %w", err)
			}

			return body.Token, nil
		},
	}))
	server.router.Use(csrfTokenTemplateMiddleware())

	server.router.Get("/keyspaces", dashboard.KeyspaceList(service, logger))
	server.router.Get("/keyspaces/:ksid", dashboard.KeyspaceView(service, logger))
	server.router.Get("/keyspaces/:ksid/key_new", dashboard.KeyNew(service, logger))
	server.router.Post("/keyspaces/:ksid/key_new", HXRequest, dashboard.KeyNewHXPost(service, logger))
	server.router.Get("/sse", server.sse)
	server.router.Use("/statics", StaticsMiddleware())
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

func (ui *Server) errorHandler(c *fiber.Ctx, err error) error {
	var layouterr *layouts.Error
	var fiberr *fiber.Error
	log := ui.logger.Error

	switch {
	case errors.As(err, &layouterr):
		break
	case errors.As(err, &fiberr):
		layouterr = layouts.NewError(fiberr.Code, fiberr.Message)
	default:
		layouterr = layouts.NewError(500, "Something went wrong", err)
	}

	if layouterr.Code < 500 {
		log = ui.logger.Warn
	}
	if c.Get("HX-Request") == "true" {
		c.Set("HX-Reswap", "innerHTML")
		c.Set("HX-Retarget", "#notification")
		return c.Status(200).Render("layouts/error@partial", layouterr.Message)
	}
	log(layouterr.Message, "code", layouterr.Code, "err", layouterr.Error())
	err = c.Status(layouterr.Code).Render("layouts/error", fiber.Map{
		"Title": "Driplimit - Error",
		"Breadcrumbs": []layouts.Breadcrumb{
			{
				Href: "/",
				Name: "Home",
			},
			{
				Href: c.Path(),
				Name: "Error",
			},
		},
		"Message": layouterr.Message,
	}, "layouts/dashboard", "layouts/page")
	if err != nil {
		log(err.Error())
	}
	return nil
}
