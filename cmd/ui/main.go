package main

import (
	"log/slog"
	"os"

	"github.com/i4n-co/driplimit/pkg/client"
	"github.com/i4n-co/driplimit/pkg/ui"
)

func main() {
	httpcli := client.New("localhost:7131").WithServiceToken("sk_6JF+gSK7bH9yUT6K+ZOzUQwHI7dXyQOWJLsYF2QcbtTTuY5A8ps0gvPgxdAGjGM*")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if err := ui.New(httpcli, logger).Listen("127.0.0.1:8080"); err != nil {
		logger.Error("failed to start driplimit ui server", "err", err.Error())
		os.Exit(1)
	}
}
