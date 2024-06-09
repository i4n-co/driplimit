package store

import (
	"context"
	"fmt"

	"github.com/i4n-co/driplimit"
)

// KeyspacesPoliciesModel is the database model for service key keyspaces policies
type KeyspacesPoliciesModel struct {
	SKID  string `db:"skid"`
	KSID  string `db:"ksid"`
	Read  bool   `db:"read"`
	Write bool   `db:"write"`
}

func (s *Store) SetKeyspacesPolicies(ctx context.Context, skid string, policies driplimit.Policies) error {
	if len(policies) == 0 {
		return nil
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = s.db.ExecContext(ctx, `
		DELETE FROM keyspaces_policies WHERE skid = ?
	`, skid)
	if err != nil {
		return fmt.Errorf("failed to delete sk keyspace policies: %w", err)
	}

	for id, policy := range policies {
		ks, err := s.GetKeyspaceByID(ctx, id)
		if err != nil {
			return err
		}
		_, err = s.db.NamedExecContext(ctx, `
			INSERT INTO keyspaces_policies (
				skid,
				ksid,
				read,
				write
			) VALUES (
				:skid,
				:ksid,
				:read,
				:write
			)
		`, KeyspacesPoliciesModel{
			SKID:  skid,
			KSID:  ks.KSID,
			Read:  policy.Read,
			Write: policy.Write,
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) GetKeyspacesPolicies(ctx context.Context, skid string) (driplimit.Policies, error) {
	policies := make([]KeyspacesPoliciesModel, 0)
	err := s.db.SelectContext(ctx, &policies, `
		SELECT * FROM keyspaces_policies WHERE skid = ?
	`, skid)
	if err != nil {
		return nil, fmt.Errorf("failed to get sk keyspace policies: %w", err)
	}

	kps := make(driplimit.Policies, 0)
	for _, kp := range policies {
		kps[kp.KSID] = driplimit.Policy{
			Read:  kp.Read,
			Write: kp.Write,
		}
	}
	return kps, nil
}
