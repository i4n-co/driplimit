package driplimit

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/i4n-co/driplimit/pkg/generate"
)

type ServiceKey struct {
	SKID              string    `json:"skid"`
	Description       string    `json:"description"`
	Admin             bool      `json:"admin"`
	Token             string    `json:"token,omitempty"`
	KeyspacesPolicies Policies  `json:"keyspaces_policies,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type ServiceKeyList struct {
	List        ListMetadata  `json:"list"`
	ServiceKeys []*ServiceKey `json:"service_keys"`
}

type ServiceKeyGetPayload struct {
	SKID  string `json:"skid" description:"The id of the service key to get (skid takes precedence over token)"`
	Token string `json:"token" description:"The token of the service key to get"`
}

func (r *ServiceKeyGetPayload) Validate(validator *validator.Validate) error {
	if r.SKID == "" && r.Token == "" {
		return ErrInvalidPayload
	}
	return nil
}

func (r *ServiceKeyGetPayload) By() (field, value string) {
	if r.SKID != "" {
		return "skid", r.SKID
	}
	return "token_hash", generate.Hash(r.Token)
}

type ServiceKeyCreatePayload struct {
	Description       string   `json:"description" description:"The description of the service key"`
	Admin             bool     `json:"admin" description:"The admin flag of the service key"`
	KeyspacesPolicies Policies `json:"keyspaces_policies" description:"The keyspaces policies of the service key. Map keys are the keyspace ids and the values are the policies for the keyspace"`
}

func (r *ServiceKeyCreatePayload) Validate(validator *validator.Validate) error {
	return validator.Struct(r)
}

type ServiceKeyListPayload struct {
	List ListPayload `json:"list" description:"The list options"`
}

func (r *ServiceKeyListPayload) Validate(validator *validator.Validate) error {
	return r.List.Validate(validator)
}

type ServiceKeyDeletePayload struct {
	SKID string `json:"skid" validate:"required" description:"The id of the service key to delete"`
}

func (r *ServiceKeyDeletePayload) Validate(validator *validator.Validate) error {
	return validator.Struct(r)
}
