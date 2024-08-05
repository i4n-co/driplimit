package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
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
	service driplimit.ServiceWithToken
	mu      *sync.Mutex
	rpcs    rpcs
	router  *fiber.App
	logger  *slog.Logger
	cfg     *config.Config
}

// New creates an API server
func New(cfg *config.Config, service driplimit.ServiceWithToken) *Server {
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
	server.RegisterRPC(v1, server.keysCreate())
	server.RegisterRPC(v1, server.keysCheck())
	server.RegisterRPC(v1, server.keysList())
	server.RegisterRPC(v1, server.keysGet())
	server.RegisterRPC(v1, server.keysDelete())

	// Keyspaces namespace
	server.RegisterRPC(v1, server.keyspacesGet())
	server.RegisterRPC(v1, server.keyspacesList())
	server.RegisterRPC(v1, server.keyspacesCreate())
	server.RegisterRPC(v1, server.keyspacesDelete())

	// ServiceKeys namespace
	server.RegisterRPC(v1, server.serviceKeysCurrent())
	server.RegisterRPC(v1, server.serviceKeysGet())
	server.RegisterRPC(v1, server.serviceKeysList())
	server.RegisterRPC(v1, server.serviceKeysDelete())
	server.RegisterRPC(v1, server.serviceKeysCreate())
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

// errorHandler centrally handles errors returned by the service
func (api *Server) errorHandler(ctx *fiber.Ctx, err error) error {
	var jsonSyntaxErr *json.SyntaxError
	var jsonUnmarshalErr *json.UnmarshalTypeError
	var fe *fiber.Error
	var ve validator.ValidationErrors
	switch {
	case errors.Is(err, driplimit.ErrUnauthorized):
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	case errors.Is(err, driplimit.ErrInvalidPayload):
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrInvalidPayload)).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, driplimit.ErrRateLimitExceeded):
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrRateLimitExceeded)).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, driplimit.ErrKeyExpired):
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrKeyExpired)).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, driplimit.ErrCannotDeleteItself):
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, driplimit.ErrAlreadyExists):
		itemAlreadyExists := driplimit.ErrItemAlreadyExists("")
		if errors.As(err, &itemAlreadyExists) {
			return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrAlreadyExists)).JSON(fiber.Map{"error": itemAlreadyExists.Error()})
		}
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrAlreadyExists)).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, driplimit.ErrNotFound):
		itemNotFound := driplimit.ErrItemNotFound("")
		if errors.As(err, &itemNotFound) {
			return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrNotFound)).JSON(fiber.Map{"error": itemNotFound.Error()})
		}
		return ctx.Status(driplimit.HTTPCodeFromErr(driplimit.ErrNotFound)).JSON(fiber.Map{"error": driplimit.ErrNotFound.Error()})
	case errors.As(err, &ve):
		fields := make([]string, 0, len(ve))
		for _, e := range ve {
			fields = append(fields, e.Field())
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": driplimit.ErrInvalidPayload.Error(), "invalid_fields": fields})
	case errors.As(err, &jsonSyntaxErr):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	case errors.As(err, &jsonUnmarshalErr):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json", "invalid_fields": []string{jsonUnmarshalErr.Field}})
	case errors.As(err, &fe):
		return ctx.Status(fe.Code).JSON(fiber.Map{"error": fe.Message})
	default:
		api.logger.Error("internal server error", "err", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
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
