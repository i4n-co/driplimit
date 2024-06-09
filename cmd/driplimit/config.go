package main

import (
	"context"
	"fmt"
	"os"

	"github.com/i4n-co/driplimit/pkg/config"
)

func loadConfig(ctx context.Context, configPath string) (cfg *config.Config, err error) {
	if configPath == "" {
		cfg, err = config.FromEnv(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from env: %w", err)
		}
		return cfg, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	cfg, err = config.FromEnvFile(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from file: %w", err)
	}

	return cfg, nil
}
