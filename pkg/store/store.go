package store

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

// Store is a store implementation using SQLite.
type Store struct {
	db        *sqlx.DB
	validator *validator.Validate
}

// New creates a new SQLite store.
func New(ctx context.Context, db *sqlx.DB) (*Store, error) {
	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")
	db.Exec("PRAGMA synchronous = NORMAL")
	sqlite := &Store{
		db:        db,
		validator: validator.New(),
	}

	if err := sqlite.migrate(ctx); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return sqlite, nil
}

// Close closes the SQLite store.
func (s *Store) Close() error {
	return s.db.Close()
}
