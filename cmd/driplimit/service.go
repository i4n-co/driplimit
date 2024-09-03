package main

import (
	"context"
	"fmt"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/authoritative"
	"github.com/i4n-co/driplimit/pkg/client"
	"github.com/i4n-co/driplimit/pkg/config"
	"github.com/i4n-co/driplimit/pkg/proxycache"
	"github.com/i4n-co/driplimit/pkg/store"
	"github.com/jmoiron/sqlx"
)

// initService initializes the driplimit service. If the configuration specifies
// a proxy, it will use the proxy cache. If not, it will use the authoritative service.
// If the configuration specifies async authoritative, it will wrap the authoritative
// service with the proxy cache.
func initService(ctx context.Context, cfg *config.Config) (driplimit.Service, error) {
	if cfg.IsProxy() {
		return driplimit.NewServiceValidator(
			proxycache.New(ctx, cfg,
				client.New(cfg.UpstreamURL),
			),
		), nil
	}

	store, err := initStore(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize store: %w", err)
	}

	if cfg.RootServiceKeyToken != "" {
		err = store.InitRootServiceKeyToken(ctx, cfg.RootServiceKeyToken)
		if err != nil {
			return nil, fmt.Errorf("failed to init root service key token: %w", err)
		}

		cfg.Logger().Info("root service token successfully set", "skid", "sk_root")
	}

	authoritative := authoritative.NewService(store)
	authzservice := driplimit.NewAuthorizer(authoritative)
	if cfg.IsAsyncAuthoritative() {
		return driplimit.NewServiceValidator(
			proxycache.New(ctx, cfg, authzservice),
		), nil
	}

	return driplimit.NewServiceValidator(authzservice), nil
}

// initStore initializes the database connection. If the configuration specifies
// a database path, it will use that path. If not, it will use an in-memory database.
func initStore(ctx context.Context, cfg *config.Config) (*store.Store, error) {
	if cfg.InMemoryDatabase() {
		cfg.Logger().Warn("using in-memory database")
		db, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			return nil, fmt.Errorf("failed to open in-memory database: %w", err)
		}
		return store.New(ctx, db)
	}
	db, err := sqlx.Open("sqlite3", fmt.Sprintf("file:%s", cfg.DatabasePath()))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return store.New(ctx, db)
}
