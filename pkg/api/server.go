package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
)

// Server is the API server
type Server struct {
	service driplimit.Service
	mu      *sync.Mutex
	rpcs    rpcs
	router  *fiber.App
	logger  *slog.Logger
	cfg     *config.Config
}

// New creates an API server
func New(cfg *config.Config, service driplimit.Service) *Server {
	server := new(Server)
	server.mu = new(sync.Mutex)
	server.cfg = cfg
	server.service = service
	server.logger = cfg.Logger().With("component", "api")
	network := fiber.NetworkTCP4
	if cfg.UseIPv6Addr() {
		network = fiber.NetworkTCP6
	}
	server.router = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          server.errorHandler,
		Network:               network,
	})
	server.router.Use(slogfiber.NewWithFilters(server.logger, slogfiber.IgnorePath("/healthz")))
	server.router.Use(compress.New(compress.Config{Level: compressionMode(cfg)}))
	server.router.Use(recover.New())

	server.router.Get("/healthz", healthz)

	v1 := server.router.Group("/v1")
	v1.Use(authenticate())

	// Keys namespace
	server.registerRPC(v1, server.keysCreate())
	server.registerRPC(v1, server.keysCheck())
	server.registerRPC(v1, server.keysList())
	server.registerRPC(v1, server.keysGet())
	server.registerRPC(v1, server.keysDelete())

	// Keyspaces namespace
	server.registerRPC(v1, server.keyspacesGet())
	server.registerRPC(v1, server.keyspacesList())
	server.registerRPC(v1, server.keyspacesCreate())
	server.registerRPC(v1, server.keyspacesDelete())

	// ServiceKeys namespace
	server.registerRPC(v1, server.serviceKeysCurrent())
	server.registerRPC(v1, server.serviceKeysGet())
	server.registerRPC(v1, server.serviceKeysList())
	server.registerRPC(v1, server.serviceKeysDelete())
	server.registerRPC(v1, server.serviceKeysCreate())
	server.registerRPC(v1, server.serviceKeysSetToken())
	return server
}

// Listen starts the API server on the given address
func (api *Server) Listen(addr string) error {
	api.logger.Info("starting driplimit...", "addr", addr)
	return api.router.Listen(addr)
}

// ShutdownWithContext shuts down the API server gracefully
func (api *Server) ShutdownWithContext(ctx context.Context) error {
	api.logger.Info("shutting down driplimit...")
	return api.router.ShutdownWithContext(ctx)
}

// Test is used for internal debugging by passing a *http.Request.
func (api *Server) Test(req *http.Request) (*http.Response, error) {
	return api.router.Test(req, 100)
}

type Err struct {
	Message       string   `json:"error"`
	InvalidFields []string `json:"invalid_fields,omitempty"`
}

func (e *Err) Error() string {
	return e.Message
}

// errorHandler centrally handles errors returned by the service
func (api *Server) errorHandler(ctx *fiber.Ctx, err error) error {
	var jsonSyntaxErr *json.SyntaxError
	var jsonUnmarshalErr *json.UnmarshalTypeError
	var fe *fiber.Error
	var ve validator.ValidationErrors
	switch {
	case errors.Is(err, driplimit.ErrUnauthorized):
		return ctx.Status(fiber.StatusUnauthorized).JSON(Err{Message: "unauthorized"})
	case errors.Is(err, driplimit.ErrInvalidPayload):
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrInvalidPayload)).JSON(Err{Message: err.Error()})
	case errors.Is(err, driplimit.ErrInvalidExpiration):
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrInvalidExpiration)).JSON(Err{Message: err.Error()})
	case errors.Is(err, driplimit.ErrRateLimitExceeded):
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrRateLimitExceeded)).JSON(Err{Message: err.Error()})
	case errors.Is(err, driplimit.ErrKeyExpired):
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrKeyExpired)).JSON(Err{Message: err.Error()})
	case errors.Is(err, driplimit.ErrCannotDeleteItself):
		return ctx.Status(fiber.StatusForbidden).JSON(Err{Message: err.Error()})
	case errors.Is(err, driplimit.ErrAlreadyExists):
		itemAlreadyExists := driplimit.ErrItemAlreadyExists("")
		if errors.As(err, &itemAlreadyExists) {
			return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrAlreadyExists)).JSON(Err{Message: itemAlreadyExists.Error()})
		}
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrAlreadyExists)).JSON(Err{Message: err.Error()})
	case errors.Is(err, driplimit.ErrNotFound):
		itemNotFound := driplimit.ErrItemNotFound("")
		if errors.As(err, &itemNotFound) {
			return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrNotFound)).JSON(Err{Message: itemNotFound.Error()})
		}
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrNotFound)).JSON(Err{Message: driplimit.ErrNotFound.Error()})
	case errors.As(err, &ve):
		fields := make([]string, 0, len(ve))
		for _, e := range ve {
			fields = append(fields, e.Field())
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(Err{Message: driplimit.ErrInvalidPayload.Error(), InvalidFields: fields})
	case errors.As(err, &jsonSyntaxErr):
		return ctx.Status(fiber.StatusBadRequest).JSON(Err{Message: "invalid json"})
	case errors.As(err, &jsonUnmarshalErr):
		return ctx.Status(fiber.StatusBadRequest).JSON(Err{Message: "invalid json", InvalidFields: []string{jsonUnmarshalErr.Field}})
	case errors.As(err, &fe):
		return ctx.Status(fe.Code).JSON(Err{Message: fe.Message})
	default:
		api.logger.Error("internal server error", "err", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(Err{Message: "internal server error"})
	}
}

// compressionMode returns the compression level based on the configuration
func compressionMode(cfg *config.Config) compress.Level {
	if cfg.GzipCompression {
		return compress.LevelBestSpeed
	}
	return compress.LevelDisabled
}

// healthz always returns true. This handler is used for monitoring purpose.
func healthz(c *fiber.Ctx) error { return c.JSON(map[string]bool{"healthy": true}) }
