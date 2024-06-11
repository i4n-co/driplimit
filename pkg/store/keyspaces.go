package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/generate"
	"github.com/mattn/go-sqlite3"
)

// KeyspaceModel represents the database model for a keyspace.
type KeyspaceModel struct {
	KSID                    string        `db:"ksid"`
	Name                    string        `db:"name"`
	KeysPrefix              string        `db:"keys_prefix"`
	DeletedAt               TimeNano      `db:"deleted_at"`
	RateLimitLimit          int64         `db:"rate_limit_limit"`
	RateLimitRefillRate     int64         `db:"rate_limit_refill_rate"`
	RateLimitRefillInterval time.Duration `db:"rate_limit_refill_interval"`
}

// ToKeyspace converts the keyspace model to a keyspace.
func (k *KeyspaceModel) ToKeyspace() *driplimit.Keyspace {
	ks := &driplimit.Keyspace{
		KSID:       k.KSID,
		Name:       k.Name,
		KeysPrefix: k.KeysPrefix,
	}
	if k.RateLimitLimit > 0 {
		ks.Ratelimit = &driplimit.Ratelimit{
			Limit:          k.RateLimitLimit,
			RefillRate:     k.RateLimitRefillRate,
			RefillInterval: driplimit.Milliseconds{Duration: k.RateLimitRefillInterval},
		}
	}
	return ks
}

// CreateKeyspace creates a new keyspace with the given payload.
func (s *Store) CreateKeyspace(ctx context.Context, payload driplimit.KeyspaceCreatePayload) (*driplimit.Keyspace, error) {
	ks := new(KeyspaceModel)
	ks.KSID = generate.IDWithPrefix("ks_")
	ks.Name = payload.Name
	ks.KeysPrefix = payload.KeysPrefix
	if payload.Ratelimit.Configured() {
		ks.RateLimitLimit = payload.Ratelimit.Limit
		ks.RateLimitRefillRate = payload.Ratelimit.RefillRate
		ks.RateLimitRefillInterval = payload.Ratelimit.RefillInterval.Duration
	}
	_, err := s.db.NamedExecContext(ctx, `
		INSERT INTO keyspaces (
			ksid, 
			name,
			keys_prefix,
			rate_limit_limit,
			rate_limit_refill_rate,
			rate_limit_refill_interval
		) 
		VALUES (
			:ksid, 
			:name,
			:keys_prefix,
			:rate_limit_limit,
			:rate_limit_refill_rate,
			:rate_limit_refill_interval
		)`, ks)
	if err != nil {
		// unique constraint violation
		sqliteConstraintErr := new(sqlite3.Error)
		if errors.As(err, sqliteConstraintErr) {
			if sqliteConstraintErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, driplimit.ErrItemAlreadyExists("keyspace")
			}
		}
		return nil, fmt.Errorf("failed to create keyspace: %w", err)
	}
	return ks.ToKeyspace(), nil
}

// GetKeyspaceByID returns a keyspace based on the given ID.
func (s *Store) GetKeyspaceByID(ctx context.Context, id string) (*driplimit.Keyspace, error) {
	ks := new(KeyspaceModel)
	err := s.db.GetContext(ctx, ks, "SELECT * FROM v_keyspaces WHERE ksid = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, driplimit.ErrItemNotFound("keyspace")
		}
		return nil, fmt.Errorf("failed to get keyspace by id: %w", err)
	}
	return ks.ToKeyspace(), nil
}

// ListKeyspaces returns a list of keyspaces based on the given payload.
func (s *Store) ListKeyspaces(ctx context.Context, payload driplimit.KeyspaceListPayload) (*driplimit.KeyspaceList, error) {
	totalCount := 0
	ks := make([]*KeyspaceModel, 0)

	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}
	defer conn.Close()

	from := "v_keyspaces"
	if payload.FilterBySKIDKeyspacesPolicies != "" {
		from = fmt.Sprintf(`(
			select v_keyspaces.* from v_keyspaces
			join service_keys_keyspaces_policies  ON v_keyspaces.ksid = service_keys_keyspaces_policies.ksid
			where service_keys_keyspaces_policies.skid = '%s'
			and service_keys_keyspaces_policies.read = 1
		)`, payload.FilterBySKIDKeyspacesPolicies)
	}
	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY name LIMIT $1 OFFSET $2", from)

	err = conn.SelectContext(ctx, &ks, sql, payload.List.Limit, payload.List.Offset())
	if err != nil {
		return nil, fmt.Errorf("failed to list keyspaces: %w", err)
	}

	sql = fmt.Sprintf("SELECT COUNT(*) FROM %s", from)
	err = conn.GetContext(ctx, &totalCount, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to count keyspaces: %w", err)
	}

	kslist := &driplimit.KeyspaceList{
		List:      driplimit.NewListMetadata(payload.List, totalCount),
		Keyspaces: make([]*driplimit.Keyspace, 0),
	}
	for _, k := range ks {
		kslist.Keyspaces = append(kslist.Keyspaces, k.ToKeyspace())
	}
	return kslist, nil
}

// DeleteKeyspace deletes a keyspace based on the given payload.
func (s *Store) DeleteKeyspace(ctx context.Context, payload driplimit.KeyspaceDeletePayload) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, "UPDATE keyspaces SET deleted_at = $1 WHERE ksid = $2 AND deleted_at = 0", TimeNano{Time: time.Now()}, payload.KSID)
	if err != nil {
		return fmt.Errorf("failed to delete keyspace: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return driplimit.ErrItemNotFound("keyspace")
	}

	_, err = tx.ExecContext(ctx, "UPDATE keys SET deleted_at = $1 WHERE ksid = $2 AND deleted_at = 0", TimeNano{Time: time.Now()}, payload.KSID)
	if err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM keyspaces_policies WHERE ksid = ?", payload.KSID)
	if err != nil {
		return fmt.Errorf("failed to delete keyspace policies: %w", err)
	}

	return tx.Commit()
}
