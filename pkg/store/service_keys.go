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

// ServiceKeyModel represents the database model for a service key.
type ServiceKeyModel struct {
	SKID        string   `db:"skid"`
	TokenHash   string   `db:"token_hash"`
	Admin       bool     `db:"admin"`
	Description string   `db:"description"`
	CreatedAt   TimeNano `db:"created_at"`
	DeletedAt   TimeNano `db:"deleted_at"`
}

// ServiceKey returns the service key from the model.
func (r *ServiceKeyModel) ServiceKey() *driplimit.ServiceKey {
	return &driplimit.ServiceKey{
		SKID:        r.SKID,
		Admin:       r.Admin,
		Description: r.Description,
		CreatedAt:   r.CreatedAt.Time,
	}
}

// NewServiceKey creates a new service key model from a service key.
func NewServiceKeyModel(sk driplimit.ServiceKey) *ServiceKeyModel {
	if sk.CreatedAt.IsZero() {
		sk.CreatedAt = time.Now()
	}
	return &ServiceKeyModel{
		SKID:      sk.SKID,
		TokenHash: generate.Hash(sk.Token),
		Admin:     sk.Admin,
		CreatedAt: TimeNano{Time: sk.CreatedAt},
		DeletedAt: TimeNano{Time: time.Time{}},
	}
}

// CreateServiceKey creates a new service key with the given payload and returns the service key and its generated token.
func (s *Store) CreateServiceKey(ctx context.Context, payload driplimit.ServiceKeyCreatePayload) (sk *driplimit.ServiceKey, key *string, err error) {
	model := new(ServiceKeyModel)
	generatedKey := "sk_" + generate.Token()
	model.SKID = generate.IDWithPrefix("sk_")
	model.TokenHash = generate.Hash(generatedKey)
	model.Admin = payload.Admin
	model.Description = payload.Description
	model.CreatedAt = TimeNano{Time: time.Now()}

	_, err = s.db.NamedExecContext(ctx, `
		INSERT INTO service_keys (
			skid, 
			token_hash,
			admin,
			description,
			created_at
		) VALUES (
			:skid,
			:token_hash,
			:admin,
			:description,
			:created_at
		)
	`, model)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to insert service key: %w", err)
	}

	err = s.SetKeyspacesPolicies(ctx, model.SKID, payload.KeyspacesPolicies)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set keyspaces policies: %w", err)
	}

	sk = model.ServiceKey()
	sk.KeyspacesPolicies = payload.KeyspacesPolicies

	return sk, &generatedKey, nil
}

// GetServiceKey returns a service key based on the given payload and populates the keyspace policies.
func (s *Store) GetServiceKey(ctx context.Context, payload driplimit.ServiceKeyGetPayload) (*driplimit.ServiceKey, error) {
	field, value := payload.By()
	sk, err := s.getServiceKeyBy(ctx, field, value)
	if err != nil {
		return nil, fmt.Errorf("failed to get service key: %w", err)
	}

	sk.KeyspacesPolicies, err = s.GetKeyspacesPolicies(ctx, sk.SKID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service key keyspace policies: %w", err)
	}

	return sk, nil
}

func (s *Store) getServiceKeyBy(ctx context.Context, field string, value string) (*driplimit.ServiceKey, error) {
	model := new(ServiceKeyModel)
	err := s.db.GetContext(ctx, model, fmt.Sprintf("SELECT * FROM v_service_keys WHERE %s = $1", field), value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, driplimit.ErrItemNotFound("service key")
		}
		return nil, fmt.Errorf("failed to get service key by %s: %w", field, err)
	}

	return model.ServiceKey(), nil
}

// ListServiceKeys returns a list of service keys based on the given payload.
func (s *Store) ListServiceKeys(ctx context.Context, payload driplimit.ServiceKeyListPayload) (*driplimit.ServiceKeyList, error) {
	totalCount := 0
	models := make([]*ServiceKeyModel, 0)

	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}
	defer conn.Close()

	err = conn.SelectContext(ctx, &models, "SELECT * FROM v_service_keys ORDER BY created_at LIMIT $1 OFFSET $2", payload.List.Limit, payload.List.Offset())
	if err != nil {
		return nil, fmt.Errorf("failed to list service keys: %w", err)
	}

	err = conn.GetContext(ctx, &totalCount, "SELECT COUNT(*) FROM v_service_keys")
	if err != nil {
		return nil, fmt.Errorf("failed to count keys: %w", err)
	}

	list := &driplimit.ServiceKeyList{
		List:        driplimit.NewListMetadata(payload.List, totalCount),
		ServiceKeys: make([]*driplimit.ServiceKey, 0),
	}
	for _, model := range models {
		list.ServiceKeys = append(list.ServiceKeys, model.ServiceKey())
	}

	return list, nil
}

// DeleteServiceKey deletes a service key based on the given payload.
func (s *Store) DeleteServiceKey(ctx context.Context, payload driplimit.ServiceKeyDeletePayload) error {

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, "UPDATE service_keys SET deleted_at=$1 WHERE skid = $2 AND deleted_at = 0", TimeNano{Time: time.Now()}, payload.SKID)
	if err != nil {
		return fmt.Errorf("failed to delete service key: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return driplimit.ErrItemNotFound("service key")
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM keyspaces_policies WHERE skid = $1", payload.SKID)
	if err != nil {
		return fmt.Errorf("failed to delete keyspaces policies: %w", err)
	}

	return tx.Commit()
}
