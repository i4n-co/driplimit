package ui

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/i4n-co/driplimit/pkg/client"
	"github.com/i4n-co/driplimit/pkg/ui/views/dashboard"
	"github.com/i4n-co/driplimit/pkg/ui/views/layouts"
)

type ContextKey string

type Server struct {
	service *client.HTTP
	router  *fiber.App
	logger  *slog.Logger
	events  chan string
}

func New(service *client.HTTP, logger *slog.Logger) *Server {
	server := &Server{
		service: service,
		logger:  logger,
		events:  make(chan string, 100),
	}
	server.router = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 NewTemplateEngine(),
		ErrorHandler:          server.errorHandler,
		PassLocalsToViews:     true,
	})
	server.router.Use(encryptcookie.New(encryptcookie.Config{
		Key:    "c2VjcmV0LXRoaXJ0eS0yLWNoYXJhY3Rlci1zdHJpbgo=",
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

	r := server.router.Group("/").Use(dashboard.AuthMiddleware(service))
	r.Get("/keyspaces", dashboard.KeyspaceList(logger))
	r.Get("/keyspaces/:ksid", dashboard.KeyspaceView(logger))
	r.Get("/keyspaces/:ksid/key_new", dashboard.KeyNew(logger))
	r.Post("/keyspaces/:ksid/key_new", HXRequest, dashboard.KeyNewHXPost(logger))
	r.Get("/sse", server.sse)
	r.Use("/statics", StaticsMiddleware())
	r.Get("/login", dashboard.Login(logger))
	r.Post("/login", dashboard.LoginHXPost(logger))
	r.Get("/logout", dashboard.Logout(logger))

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
