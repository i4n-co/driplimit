package driplimit

import "github.com/go-playground/validator/v10"

// Keyspace represents a driplimit keyspace.
type Keyspace struct {
	KSID       string     `json:"ksid"`
	Name       string     `json:"name"`
	KeysPrefix string     `json:"keys_prefix"`
	Ratelimit  *Ratelimit `json:"ratelimit,omitempty"`
}

// ConfiguredRateLimit returns true if the rate limit is configured for the keyspace.
func (ks *Keyspace) ConfiguredRateLimit() bool {
	return ks.Ratelimit.Configured()
}

// KeyspaceCreatePayload represents the payload for creating a keyspace.
type KeyspaceCreatePayload struct {
	Name       string           `json:"name" validate:"required" description:"The name of the keyspace"`
	KeysPrefix string           `json:"keys_prefix" validate:"required,gte=1,lte=16" description:"The prefix for the keys in the keyspace"`
	Ratelimit  RatelimitPayload `json:"ratelimit" description:"The default rate limit configuration for keys in the keyspace"`
}

// Validate validates the keyspace create payload.
func (ks *KeyspaceCreatePayload) Validate(validator *validator.Validate) error {
	return validator.Struct(ks)
}

type KeyspaceGetPayload struct {
	KSID string `json:"ksid" validate:"required" description:"The id of the keyspace to get"`
}

// Validate validates the keyspace get payload.
func (k *KeyspaceGetPayload) Validate(validator *validator.Validate) error {
	return validator.Struct(k)
}

// KeyspaceList represents a list of keyspaces.
type KeyspaceList struct {
	List      ListMetadata `json:"list"`
	Keyspaces []*Keyspace  `json:"keyspaces"`
}

// KeyspaceListPayload represents the payload for listing keyspaces.
type KeyspaceListPayload struct {
	List                          ListPayload `json:"list" description:"The list options"`
	FilterBySKIDKeyspacesPolicies string      `json:"-"` // filter by sk keyspaces policies
}

// Validate validates the list payload.
func (kl *KeyspaceListPayload) Validate(validator *validator.Validate) error {
	err := kl.List.Validate(validator)
	if err != nil {
		return err
	}
	return validator.Struct(kl)
}

// KeyspaceDeletePayload is the payload for deleting a keyspace.
type KeyspaceDeletePayload struct {
	KSID string `json:"ksid" validate:"required" description:"The id of the keyspace to delete"`
}

// Validate validates the key get payload.
func (k *KeyspaceDeletePayload) Validate(validator *validator.Validate) error {
	return validator.Struct(k)
}
