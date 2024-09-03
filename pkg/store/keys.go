package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/generate"
)

// KeyModel represents the database model for a key.
type KeyModel struct {
	KID       string   `db:"kid"`
	KSID      string   `db:"ksid"`
	TokenHash string   `db:"token_hash"`
	LastUsed  TimeNano `db:"last_used"`
	ExpiresAt TimeNano `db:"expires_at"`
	CreatedAt TimeNano `db:"created_at"`
	DeletedAt TimeNano `db:"deleted_at"`

	RateLimitStateLastRefilled TimeNano      `db:"rate_limit_state_last_refilled"`
	RateLimitStateRemaining    int64         `db:"rate_limit_state_remaining"`
	RateLimitLimit             int64         `db:"rate_limit_limit"`
	RateLimitRefillRate        int64         `db:"rate_limit_refill_rate"`
	RateLimitRefillInterval    time.Duration `db:"rate_limit_refill_interval"`
}

// ConfiguredRateLimit returns true if the rate limit is configured for the key.
func (model *KeyModel) ConfiguredRateLimit() bool {
	return model.RateLimitLimit > 0
}

// NewKeyModel creates a new key model from a key.
func NewKeyModel(key driplimit.Key) *KeyModel {
	model := &KeyModel{
		KID:       key.KID,
		KSID:      key.KSID,
		LastUsed:  TimeNano{Time: key.LastUsed},
		ExpiresAt: TimeNano{Time: key.ExpiresAt},
		CreatedAt: TimeNano{Time: key.CreatedAt},
	}

	if key.Ratelimit != nil {
		model.RateLimitLimit = key.Ratelimit.Limit
		model.RateLimitRefillRate = key.Ratelimit.RefillRate
		model.RateLimitRefillInterval = key.Ratelimit.RefillInterval.Duration
		if key.Ratelimit.State != nil {
			model.RateLimitStateLastRefilled = TimeNano{Time: key.Ratelimit.State.LastRefilled}
			model.RateLimitStateRemaining = key.Ratelimit.State.Remaining
		}
	}

	return model
}

// ToKey converts the key model to a key.
func (model *KeyModel) ToKey() *driplimit.Key {
	key := &driplimit.Key{
		KID:       model.KID,
		KSID:      model.KSID,
		LastUsed:  model.LastUsed.Time,
		ExpiresAt: model.ExpiresAt.Time,
		CreatedAt: model.CreatedAt.Time,
	}

	if model.RateLimitLimit > 0 {
		key.Ratelimit = &driplimit.Ratelimit{
			State: &driplimit.RatelimitState{
				LastRefilled: model.RateLimitStateLastRefilled.Time,
				Remaining:    model.RateLimitStateRemaining,
			},
			Limit:          model.RateLimitLimit,
			RefillRate:     model.RateLimitRefillRate,
			RefillInterval: driplimit.Milliseconds{Duration: model.RateLimitRefillInterval},
		}
	}

	return key
}

// ValidateToken validates the token of the key.
func (k *KeyModel) ValidateToken(token string) bool {
	return k.TokenHash == generate.Hash(token)
}

// CreateKey creates a new key.
func (sqlite *Store) CreateKey(ctx context.Context, payload driplimit.KeyCreatePayload) (*driplimit.Key, error) {
	ks, err := sqlite.GetKeyspaceByID(ctx, payload.KSID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, driplimit.ErrItemNotFound("keyspace")
		}
		return nil, fmt.Errorf("failed to get keyspace by id: %w", err)
	}

	model := new(KeyModel)
	token := ks.KeysPrefix + generate.Token()

	model.KID = "k_" + generate.ID()
	model.KSID = ks.KSID
	model.ExpiresAt = TimeNano{Time: payload.ExpiresAt}
	model.CreatedAt = TimeNano{Time: time.Now()}
	model.LastUsed = TimeNano{Time: time.Time{}}
	model.TokenHash = generate.Hash(token)
	if payload.Ratelimit.Configured() {
		model.RateLimitStateRemaining = payload.Ratelimit.Limit
		model.RateLimitStateLastRefilled = TimeNano{Time: time.Now()}
		model.RateLimitLimit = payload.Ratelimit.Limit
		model.RateLimitRefillRate = payload.Ratelimit.RefillRate
		model.RateLimitRefillInterval = payload.Ratelimit.RefillInterval.Duration
	}

	_, err = sqlite.db.NamedExecContext(ctx, `
	INSERT INTO keys
	(
		kid,
		ksid,
		token_hash,
		last_used,
		expires_at,
		created_at,
		rate_limit_state_last_refilled,
		rate_limit_state_remaining,
		rate_limit_limit,
		rate_limit_refill_rate,
		rate_limit_refill_interval
	)
	VALUES
	(
		:kid,
		:ksid,
		:token_hash,
		:last_used,
		:expires_at,
		:created_at,
		:rate_limit_state_last_refilled,
		:rate_limit_state_remaining,
		:rate_limit_limit,
		:rate_limit_refill_rate,
		:rate_limit_refill_interval
	)`, model)
	if err != nil {
		return nil, fmt.Errorf("failed to create key: %w", err)
	}

	k := model.ToKey()
	k.Token = token
	return k, nil
}

// GetKey returns a key by the given payload. Ratelimit is set with keyspace ratelimit if not configured
// on the key itself.
func (sqlite *Store) GetKey(ctx context.Context, payload driplimit.KeyGetPayload) (key *driplimit.Key, err error) {
	field, value, err := payload.GetKeyBy()
	if err != nil {
		return nil, err
	}

	model, err := sqlite.getKeyBy(ctx, payload.KSID, field, value)
	if err != nil {
		return nil, err
	}
	if !model.ConfiguredRateLimit() {
		ks, err := sqlite.GetKeyspaceByID(ctx, model.KSID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, driplimit.ErrItemNotFound("keyspace")
			}
			return nil, fmt.Errorf("failed to get keyspace by id: %w", err)
		}
		if !ks.ConfiguredRateLimit() {
			return model.ToKey(), nil
		}
		model.RateLimitLimit = ks.Ratelimit.Limit
		model.RateLimitRefillRate = ks.Ratelimit.RefillRate
		model.RateLimitRefillInterval = ks.Ratelimit.RefillInterval.Duration
	}

	return model.ToKey(), nil
}

// getKeyBy returns a key by the given field.
func (sqlite *Store) getKeyBy(ctx context.Context, ksid string, field string, value string) (*KeyModel, error) {
	key := new(KeyModel)
	err := sqlite.db.GetContext(ctx, key, fmt.Sprintf("SELECT * FROM v_keys WHERE %s = $1 AND ksid = $2", field), value, ksid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, driplimit.ErrItemNotFound("key")
		}
		return nil, fmt.Errorf("failed to get key by %s: %w", field, err)
	}

	return key, nil
}

// UpdateLastUsed updates the last used field of the key.
func (sqlite *Store) UpdateLastUsed(ctx context.Context, key *driplimit.Key) error {
	model := NewKeyModel(*key)
	_, err := sqlite.db.ExecContext(ctx, "UPDATE keys SET last_used = $1 WHERE kid = $2", model.LastUsed, model.KID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return driplimit.ErrItemNotFound("key")
		}
		return fmt.Errorf("failed to update key last used: %w", err)
	}
	return nil
}

// DecrementKeyRemaining decrements the remaining field of the key.
func (sqlite *Store) DecrementKeyRemaining(ctx context.Context, key *driplimit.Key) error {
	model := NewKeyModel(*key)
	row := sqlite.db.QueryRow("UPDATE keys SET rate_limit_state_remaining = rate_limit_state_remaining - 1, last_used = $1 WHERE kid = $2 RETURNING rate_limit_state_remaining",
		model.LastUsed,
		model.KID,
	)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return driplimit.ErrItemNotFound("key")
		}
		return fmt.Errorf("failed to decrement remaining from key: %w", row.Err())
	}
	return row.Scan(&key.Ratelimit.State.Remaining)
}

// SetKeyRemaining sets the remaining field of the key.
func (sqlite *Store) SetKeyRemaining(ctx context.Context, key *driplimit.Key) error {
	model := NewKeyModel(*key)
	_, err := sqlite.db.Exec("UPDATE keys SET rate_limit_state_remaining = $1, rate_limit_state_last_refilled = $2 WHERE kid = $3",
		model.RateLimitStateRemaining,
		model.RateLimitStateLastRefilled,
		model.KID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return driplimit.ErrItemNotFound("key")
		}
		return fmt.Errorf("failed to set remaining to key: %w", err)
	}
	return nil
}

// ListKeys returns a list of keys based on the given payload.
func (sqlite *Store) ListKeys(ctx context.Context, payload driplimit.KeyListPayload) (klist *driplimit.KeyList, err error) {
	totalCount := 0
	keys := make([]*KeyModel, 0)

	conn, err := sqlite.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}
	defer conn.Close()

	ks, err := sqlite.GetKeyspaceByID(ctx, payload.KSID)
	if err != nil {
		return nil, fmt.Errorf("failed to get keyspace by id: %w", err)
	}
	err = conn.SelectContext(ctx, &keys, "SELECT * FROM v_keys WHERE ksid = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", ks.KSID, payload.List.Limit, payload.List.Offset())
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}
	err = conn.GetContext(ctx, &totalCount, "SELECT COUNT(*) FROM v_keys WHERE ksid = $1", ks.KSID)
	if err != nil {
		return nil, fmt.Errorf("failed to count keys: %w", err)
	}

	klist = &driplimit.KeyList{
		List: driplimit.NewListMetadata(payload.List, totalCount),
		Keys: make([]*driplimit.Key, 0),
	}
	for _, k := range keys {
		klist.Keys = append(klist.Keys, k.ToKey())
	}
	return klist, nil
}

// DeleteKey deletes a key based on the given payload.
func (sqlite *Store) DeleteKey(ctx context.Context, payload driplimit.KeyDeletePayload) error {
	res, err := sqlite.db.ExecContext(ctx, "UPDATE keys SET deleted_at = $1 WHERE kid = $2 AND ksid = $3 AND deleted_at = 0", TimeNano{Time: time.Now()}, payload.KID, payload.KSID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return driplimit.ErrItemNotFound("key")
		}
		return fmt.Errorf("failed to delete key: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return driplimit.ErrItemNotFound("key")
	}
	return nil
}
