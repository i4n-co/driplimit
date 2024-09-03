package authoritative

import (
	"context"
	"fmt"
	"time"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/store"
)

// Authoritative is the Authoritative implementation of the driplimit service.
// It is the source of truth for the rate limits and uses the store directly.
type Authoritative struct {
	store *store.Store
}

// NewService returns a new authoritative driplimit service.
func NewService(store *store.Store) *Authoritative {
	app := &Authoritative{
		store: store,
	}
	return app
}

// KeyCheck checks if the key can be used (not expired, rate limit not exceeded) and returns an error if not.
// In case of success, it decrements the remaining count of the key if the rate limit is set.
func (service *Authoritative) KeyCheck(ctx context.Context, payload driplimit.KeysCheckPayload) (key *driplimit.Key, err error) {
	key, err = service.KeyGet(ctx, driplimit.KeyGetPayload{KSID: payload.KSID, Token: payload.Token})
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	if key.Expired() {
		return nil, driplimit.ErrKeyExpired
	}

	key.LastUsed = time.Now()
	if !key.ConfiguredRatelimit() {
		if err := service.store.UpdateLastUsed(ctx, key); err != nil {
			return nil, fmt.Errorf("failed to update key last used: %w", err)
		}
		return key, nil
	}

	if key.Ratelimit.State.Remaining <= 0 {
		return nil, driplimit.ErrRateLimitExceeded
	}

	if err := service.store.DecrementKeyRemaining(ctx, key); err != nil {
		return nil, fmt.Errorf("failed to decrement key remaining: %w", err)
	}

	return key, nil
}

// KeyCreate creates a new key with the given payload and returns the key information and the token.
func (service *Authoritative) KeyCreate(ctx context.Context, payload driplimit.KeyCreatePayload) (key *driplimit.Key, token *string, err error) {
	key, token, err = service.store.CreateKey(ctx, payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create key: %w", err)
	}
	return key, token, nil
}

// KeyGet returns key based on the given payload. It ensures that the remaining count is up to date if necessary.
func (service *Authoritative) KeyGet(ctx context.Context, payload driplimit.KeyGetPayload) (key *driplimit.Key, err error) {
	key, err = service.store.GetKey(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	remainingUpdated := key.UpdateRemaining()
	if !remainingUpdated {
		return key, nil
	}

	if err := service.store.SetKeyRemaining(ctx, key); err != nil {
		return nil, fmt.Errorf("failed to set key remaining: %w", err)
	}

	return key, nil
}

// KeyList returns a list of keys based on the given payload.
func (service *Authoritative) KeyList(ctx context.Context, payload driplimit.KeyListPayload) (klist *driplimit.KeyList, err error) {
	klist, err = service.store.ListKeys(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}
	return klist, nil
}

// KeyDelete deletes a key based on the given payload.
func (service *Authoritative) KeyDelete(ctx context.Context, payload driplimit.KeyDeletePayload) (err error) {
	if err := service.store.DeleteKey(ctx, payload); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

// KeyspaceGet returns a keyspace based on the given payload.
func (service *Authoritative) KeyspaceGet(ctx context.Context, payload driplimit.KeyspaceGetPayload) (keyspace *driplimit.Keyspace, err error) {
	ks, err := service.store.GetKeyspaceByID(ctx, payload.KSID)
	if err != nil {
		return nil, fmt.Errorf("failed to get keyspace: %w", err)
	}
	return ks, nil
}

// KeyspaceCreate creates a new keyspace with the given payload and returns it.
func (service *Authoritative) KeyspaceCreate(ctx context.Context, payload driplimit.KeyspaceCreatePayload) (keyspace *driplimit.Keyspace, err error) {
	ks, err := service.store.CreateKeyspace(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyspace: %w", err)
	}
	return ks, nil
}

// KeyspaceList returns a list of keyspaces based on the given payload.
func (service *Authoritative) KeyspaceList(ctx context.Context, payload driplimit.KeyspaceListPayload) (kslist *driplimit.KeyspaceList, err error) {
	kslist, err = service.store.ListKeyspaces(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to list keyspaces: %w", err)
	}
	return kslist, nil
}

// KeyspaceDelete deletes a keyspace based on the given payload.
func (service *Authoritative) KeyspaceDelete(ctx context.Context, payload driplimit.KeyspaceDeletePayload) (err error) {
	if err := service.store.DeleteKeyspace(ctx, payload); err != nil {
		return fmt.Errorf("failed to delete keyspace: %w", err)
	}
	return nil
}

// ServiceKeyGet returns a service key based on the given payload.
func (service *Authoritative) ServiceKeyGet(ctx context.Context, payload driplimit.ServiceKeyGetPayload) (sk *driplimit.ServiceKey, err error) {
	sk, err = service.store.GetServiceKey(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get service key: %w", err)
	}
	return sk, nil
}

// ServiceKeyCreate creates a new service key with the given payload and returns the service key information.
func (service *Authoritative) ServiceKeyCreate(ctx context.Context, payload driplimit.ServiceKeyCreatePayload) (sk *driplimit.ServiceKey, err error) {
	var token *string
	sk, token, err = service.store.CreateServiceKey(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create service key: %w", err)
	}
	sk.Token = *token
	return sk, nil
}

func (service *Authoritative) ServiceKeyList(ctx context.Context, payload driplimit.ServiceKeyListPayload) (sklist *driplimit.ServiceKeyList, err error) {
	sklist, err = service.store.ListServiceKeys(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to list service keys: %w", err)
	}
	return sklist, nil
}

func (service *Authoritative) ServiceKeyDelete(ctx context.Context, payload driplimit.ServiceKeyDeletePayload) (err error) {
	if err := service.store.DeleteServiceKey(ctx, payload); err != nil {
		return fmt.Errorf("failed to delete service key: %w", err)
	}
	return nil
}

func (service *Authoritative) ServiceKeySetToken(ctx context.Context, payload driplimit.ServiceKeySetTokenPayload) (err error) {
	if err := service.store.SetServiceKeyToken(ctx, payload); err != nil {
		return fmt.Errorf("failed to set service key token: %w", err)
	}
	return nil
}
