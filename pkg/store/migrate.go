package store

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/maragudk/migrate"
)

//go:embed migrations/*.sql
var migrations embed.FS

// migrate applies the migrations to the database.
func (s *Store) migrate(ctx context.Context) error {
	sub, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get sub filesystem: %w", err)
	}

	err = migrate.Up(ctx, s.db.DB, sub)
	if err != nil {
		return fmt.Errorf("failed to apply migrate up: %w", err)
	}

	return nil
}
