package driplimit

import (
	"context"
	"errors"
	"fmt"
)

// PolicyAction represents an action that can be performed. eg. read, write
type PolicyAction string

const (
	Read  PolicyAction = "read"
	Write PolicyAction = "write"
)

// Policy gives the read and write permissions.
type Policy struct {
	Read  bool `json:"read" description:"Read permission"`
	Write bool `json:"write" description:"Write permission"`
}

// all is a wildcard for all ids. This can be found in items policies
const all = "*"

// Can checks if the policy allows the action.
func (k Policy) Can(action PolicyAction) bool {
	switch action {
	case "read":
		return k.Read
	case "write":
		return k.Write
	default:
		return false
	}
}

// Policies is a map of policies.
type Policies map[string]Policy

// Can checks if the action can be performed on item id.
func (policies Policies) Can(action PolicyAction, id string) bool {
	policy, found := policies[all]
	if found {
		if policy.Can(action) {
			return true
		}
	}

	policy, found = policies[id]
	if !found {
		return false
	}
	return policy.Can(action)
}

// Authorizer is an authorization wrapper. It implements the Service interface.
type Authorizer struct {
	driplimit Service
}

// NewAuthorizer wraps a Driplimit ServiceWithToken with an authorizer.
func NewAuthorizer(driplimit Service) *Authorizer {
	return &Authorizer{
		driplimit: driplimit,
	}
}

// caller gets the service key from the context.
// -- NO
func (a *Authorizer) caller(ctx context.Context, payload Payload) (sk *ServiceKey, err error) {
	sk, err = a.driplimit.ServiceKeyGet(ctx, ServiceKeyGetPayload{Token: payload.ServiceToken()})
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("failed to get service key: %w", err)
	}
	return sk, nil
}

func (a *Authorizer) KeyCheck(ctx context.Context, payload KeysCheckPayload) (key *Key, err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin || sk.KeyspacesPolicies.Can(Read, payload.KSID) {
		return a.driplimit.KeyCheck(ctx, payload)
	}
	return nil, ErrUnauthorized
}

func (a *Authorizer) KeyCreate(ctx context.Context, payload KeyCreatePayload) (key *Key, err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin || sk.KeyspacesPolicies.Can(Write, payload.KSID) {
		return a.driplimit.KeyCreate(ctx, payload)
	}
	return nil, ErrUnauthorized
}

func (a *Authorizer) KeyGet(ctx context.Context, payload KeyGetPayload) (key *Key, err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin || sk.KeyspacesPolicies.Can(Read, payload.KSID) {
		return a.driplimit.KeyGet(ctx, payload)
	}
	return nil, ErrUnauthorized
}

func (a *Authorizer) KeyList(ctx context.Context, payload KeyListPayload) (klist *KeyList, err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin || sk.KeyspacesPolicies.Can(Read, payload.KSID) {
		return a.driplimit.KeyList(ctx, payload)
	}
	return nil, ErrUnauthorized
}

func (a *Authorizer) KeyDelete(ctx context.Context, payload KeyDeletePayload) (err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return err
	}
	if sk.Admin || sk.KeyspacesPolicies.Can("write", payload.KSID) {
		return a.driplimit.KeyDelete(ctx, payload)
	}
	return ErrUnauthorized
}

func (a *Authorizer) KeyspaceGet(ctx context.Context, payload KeyspaceGetPayload) (keyspace *Keyspace, err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin || sk.KeyspacesPolicies.Can(Read, payload.KSID) {
		return a.driplimit.KeyspaceGet(ctx, payload)
	}
	return nil, ErrUnauthorized
}

func (a *Authorizer) KeyspaceCreate(ctx context.Context, payload KeyspaceCreatePayload) (keyspace *Keyspace, err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin {
		return a.driplimit.KeyspaceCreate(ctx, payload)
	}
	return nil, ErrUnauthorized
}

func (a *Authorizer) KeyspaceList(ctx context.Context, payload KeyspaceListPayload) (kslist *KeyspaceList, err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin || sk.KeyspacesPolicies.Can(Read, all) {
		return a.driplimit.KeyspaceList(ctx, payload)
	}
	// if the service key is not allowed to read all keyspaces, then filter by SKID
	payload.FilterBySKIDKeyspacesPolicies = sk.SKID
	return a.driplimit.KeyspaceList(ctx, payload)
}

func (a *Authorizer) KeyspaceDelete(ctx context.Context, payload KeyspaceDeletePayload) (err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return err
	}
	if sk.Admin {
		return a.driplimit.KeyspaceDelete(ctx, payload)
	}
	return ErrUnauthorized
}

func (a *Authorizer) ServiceKeyGet(ctx context.Context, payload ServiceKeyGetPayload) (sk *ServiceKey, err error) {
	sk, err = a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin {
		return a.driplimit.ServiceKeyGet(ctx, payload)
	}
	requestedServiceKey, err := a.driplimit.ServiceKeyGet(ctx, payload)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("failed to get service key: %w", err)
	}
	if requestedServiceKey.SKID != sk.SKID {
		return nil, ErrUnauthorized
	}
	return requestedServiceKey, nil
}

func (a *Authorizer) ServiceKeyCreate(ctx context.Context, payload ServiceKeyCreatePayload) (sk *ServiceKey, err error) {
	sk, err = a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin {
		return a.driplimit.ServiceKeyCreate(ctx, payload)
	}
	return nil, ErrUnauthorized
}

func (a *Authorizer) ServiceKeyList(ctx context.Context, payload ServiceKeyListPayload) (sklist *ServiceKeyList, err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return nil, err
	}
	if sk.Admin {
		return a.driplimit.ServiceKeyList(ctx, payload)
	}
	return nil, ErrUnauthorized
}

func (a *Authorizer) ServiceKeyDelete(ctx context.Context, payload ServiceKeyDeletePayload) (err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return err
	}
	if sk.SKID == payload.SKID {
		return ErrCannotDeleteItself
	}
	if sk.Admin {
		return a.driplimit.ServiceKeyDelete(ctx, payload)
	}
	return ErrUnauthorized
}

func (a *Authorizer) ServiceKeySetToken(ctx context.Context, payload ServiceKeySetTokenPayload) (err error) {
	sk, err := a.caller(ctx, payload)
	if err != nil {
		return err
	}

	if sk.Admin || sk.SKID == payload.SKID {
		return a.driplimit.ServiceKeySetToken(ctx, payload)
	}

	return ErrUnauthorized
}
