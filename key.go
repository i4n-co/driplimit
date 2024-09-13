package driplimit

import (
	"encoding/json"
	"time"

	"github.com/i4n-co/driplimit/pkg/generate"

	"github.com/go-playground/validator/v10"
)

// Key represents a driplimit key.
type Key struct {
	KID    string     `json:"kid"`
	KSID      string     `json:"ksid"`
	Token     string     `json:"token,omitempty"`
	LastUsed  time.Time  `json:"last_used"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	Ratelimit *Ratelimit `json:"ratelimit,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
// It is mandatory to implement this method to handle zero time fields.
// With the current encoding/json/v1 go implementation, the omitzero tag is
// not supported yet. Follow the discussion here:
//
//	https://github.com/golang/go/discussions/63397
//
// This "hack" is taken from the following blog post:
//
//	https://choly.ca/post/go-json-marshalling/
func (k Key) MarshalJSON() ([]byte, error) {
	type KeyAlias Key
	lastUsed := ""
	if !k.LastUsed.IsZero() {
		lastUsed = k.LastUsed.Format(time.RFC3339Nano)
	}
	expiresAt := ""
	if !k.ExpiresAt.IsZero() {
		expiresAt = k.ExpiresAt.Format(time.RFC3339Nano)
	}
	return json.Marshal(&struct {
		KeyAlias
		LastUsed  string `json:"last_used,omitempty"`
		ExpiresAt string `json:"expires_at,omitempty"`
	}{
		KeyAlias:  (KeyAlias)(k),
		LastUsed:  lastUsed,
		ExpiresAt: expiresAt,
	})
}

// since is a variable that holds the time.Since function.
// It is used to mock time.Since in tests.
var (
	since = time.Since
	now   = time.Now
)

// UpdateRemaining updates the remaining rate limit state of the key based on the
// rate limit configuration. It returns true if the remaining rate limit was updated.
func (k *Key) UpdateRemaining() (updated bool) {
	if !k.ConfiguredRatelimit() {
		return false
	}
	if k.Ratelimit.State == nil {
		k.Ratelimit.State = &RatelimitState{}
	}
	if k.Ratelimit.RefillInterval.Duration == 0 || k.Ratelimit.RefillRate == 0 {
		return false
	}
	defer func() {
		if updated {
			k.Ratelimit.State.LastRefilled = now()
		}
	}()
	sinceLastRefill := since(k.Ratelimit.State.LastRefilled)
	refills := sinceLastRefill.Nanoseconds() / k.Ratelimit.RefillInterval.Nanoseconds()
	refilled := refills * k.Ratelimit.RefillRate
	remaining := k.Ratelimit.State.Remaining + refilled
	updated = remaining != k.Ratelimit.State.Remaining
	if remaining > k.Ratelimit.Limit {
		k.Ratelimit.State.Remaining = k.Ratelimit.Limit
		return updated
	}
	k.Ratelimit.State.Remaining = remaining
	return updated
}

// ConfiguredRatelimit returns true if the rate limit is configured for the key.
func (k *Key) ConfiguredRatelimit() bool {
	return k.Ratelimit.Configured()
}

// Expired returns true if the key is expired.
func (k *Key) Expired() bool {
	if k.ExpiresAt.IsZero() {
		return false
	}
	return since(k.ExpiresAt) > 0
}

// KeyCreatePayload is the payload for creating a key.
type KeyCreatePayload struct {
	*payload
	KSID      string           `json:"ksid" validate:"required" description:"The id of the keyspace to which the key belongs to"`
	ExpiresIn Milliseconds     `json:"expires_in" description:"The duration in milliseconds after which the key expires"`
	ExpiresAt time.Time        `json:"expires_at" description:"The time at which the key expires (expires_at takes precedence over expires_in)"`
	Ratelimit RatelimitPayload `json:"ratelimit" validate:"required" description:"The rate limit configuration for the key"`
}

// Validate validates the key create payload.
func (k *KeyCreatePayload) Validate(validator *validator.Validate) error {
	if k.ExpiresIn.Duration == 0 && k.ExpiresAt.IsZero() {
		return ErrInvalidExpiration
	}

	if k.ExpiresIn.Duration > 0 && k.ExpiresAt.IsZero() {
		k.ExpiresAt = time.Now().Add(k.ExpiresIn.Duration)
	}

	return validator.Struct(k)
}

// WithServiceToken adds authentication infos to payload
func (k *KeyCreatePayload) WithServiceToken(token string) *KeyCreatePayload {
	k.payload = &payload{
		serviceToken: token,
	}
	return k
}

// KeyCheckPayload is the payload for checking a key.
type KeysCheckPayload struct {
	*payload
	KSID  string `json:"ksid" validate:"required" description:"The id of the keyspace to which the key belongs to"`
	Token string `json:"token" validate:"required" description:"The token to check"`
}

// Validate validates the key check payload.
func (k *KeysCheckPayload) Validate(validator *validator.Validate) error {
	return validator.Struct(k)
}

// WithServiceToken adds authentication infos to payload
func (k *KeysCheckPayload) WithServiceToken(token string) *KeysCheckPayload {
	k.payload = &payload{
		serviceToken: token,
	}
	return k
}

// KeyGetPayload is the payload for getting a key.
type KeyGetPayload struct {
	*payload

	KSID  string `json:"ksid" validate:"required" description:"The id of the keyspace to which the key belongs to"`
	KID   string `json:"kid" description:"The id of the key to get (kid takes precedence over token if both are provided)"`
	Token string `json:"token" description:"The token of the key to get"`
}

// Validate validates the key get payload.
func (k *KeyGetPayload) Validate(validator *validator.Validate) error {
	err := validator.Struct(k)
	if err != nil {
		return err
	}
	if k.KID == "" && k.Token == "" {
		return ErrInvalidPayload
	}
	return nil
}

// WithServiceToken adds authentication infos to payload
func (k *KeyGetPayload) WithServiceToken(token string) *KeyGetPayload {
	k.payload = &payload{
		serviceToken: token,
	}
	return k
}

// GetKeyBy returns the field and the corresponding value by which the key should be retrieved.
func (k *KeyGetPayload) GetKeyBy() (field string, value string, err error) {
	if k.KID != "" {
		return "kid", k.KID, nil
	}
	if k.Token != "" {
		return "token_hash", generate.Hash(k.Token), nil
	}
	return "", "", ErrInvalidPayload
}

// KeyList represents a list of keys.
type KeyList struct {
	List ListMetadata `json:"list"`
	Keys []*Key       `json:"keys"`
}

// KeyspaceListPayload represents the payload for listing keyspaces.
type KeyListPayload struct {
	*payload

	List ListPayload `json:"list" description:"The list options"`
	KSID string      `json:"ksid" validate:"required" description:"The id of the keyspace to which the keys belong to"`
}

// Validate validates the list payload.
func (kl *KeyListPayload) Validate(validator *validator.Validate) error {
	err := kl.List.Validate(validator)
	if err != nil {
		return err
	}
	return validator.Struct(kl)
}

// WithServiceToken adds authentication infos to payload
func (k *KeyListPayload) WithServiceToken(token string) *KeyListPayload {
	k.payload = &payload{
		serviceToken: token,
	}
	return k
}

// KeyDeletePayload is the payload for deleting a key.
type KeyDeletePayload struct {
	*payload

	KSID string `json:"ksid" validate:"required" description:"The id of the keyspace to which the key belongs to"`
	KID  string `json:"kid" validate:"required" description:"The id of the key to delete"`
}

// Validate validates the key get payload.
func (k *KeyDeletePayload) Validate(validator *validator.Validate) error {
	return validator.Struct(k)
}

// WithServiceToken adds authentication infos to payload
func (k *KeyDeletePayload) WithServiceToken(token string) *KeyDeletePayload {
	k.payload = &payload{
		serviceToken: token,
	}
	return k
}
