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
	ListMetadata
	ServiceKeys []*ServiceKey `json:"service_keys"`
}

type ServiceKeyGetPayload struct {
	SKID  string `json:"skid"`
	Token string `json:"token"`
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
	Description       string   `json:"description"`
	Admin             bool     `json:"admin"`
	KeyspacesPolicies Policies `json:"keyspaces_policies"`
}

func (r *ServiceKeyCreatePayload) Validate(validator *validator.Validate) error {
	return validator.Struct(r)
}

type ServiceKeyListPayload struct {
	ListPayload
}

func (r *ServiceKeyListPayload) Validate(validator *validator.Validate) error {
	return r.ListPayload.Validate(validator)
}

type ServiceKeyDeletePayload struct {
	SKID string `json:"skid" validate:"required"`
}

func (r *ServiceKeyDeletePayload) Validate(validator *validator.Validate) error {
	return validator.Struct(r)
}
