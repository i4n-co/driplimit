package main

import (
	"context"
	"fmt"
	"time"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/authoritative"
	"github.com/i4n-co/driplimit/pkg/config"
)

func initAdmin(ctx context.Context, cfg *config.Config) error {
	store, err := initStore(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	service := authoritative.NewService(store)

	payload := driplimit.ServiceKeyCreatePayload{
		Description: fmt.Sprintf("cli generated admin service key at %s", time.Now().Format(time.RFC3339)),
		Admin:       true,
	}
	sk, err := service.ServiceKeyCreate(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to create admin service key: %w", err)
	}

	fmt.Printf("\nAdmin service key created successfully.\n\n")
	fmt.Printf("Please store the key in a safe place. It will not be shown again.\n")
	fmt.Printf("Service Key ID:         %s\n", sk.SKID)
	fmt.Printf("Service Key Token:      %s\n", sk.Token)
	fmt.Printf("Description:            %s\n", sk.Description)
	fmt.Printf("Creation time:          %s\n\n", sk.CreatedAt.Format(time.RFC3339))

	return nil
}
