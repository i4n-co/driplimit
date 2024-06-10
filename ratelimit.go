package driplimit

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Ratelimit struct {
	State          *RatelimitState `json:"state,omitempty"`
	Limit          int64           `json:"limit" db:"rate_limit"`
	RefillRate     int64           `json:"refill_rate" db:"rate_limit_refill_rate"`
	RefillInterval Milliseconds    `json:"refill_interval" db:"rate_limit_refill_interval"`
}

// Milliseconds is a duration that is serialized as milliseconds.
type Milliseconds struct {
	time.Duration
}

// MarshalJSON marshals the duration as milliseconds.
func (d Milliseconds) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", d.Milliseconds())), nil
}

// UnmarshalJSON unmarshals the duration as milliseconds.
func (d *Milliseconds) UnmarshalJSON(data []byte) error {
	var ms int64
	if _, err := fmt.Sscanf(string(data), "%d", &ms); err != nil {
		return err
	}
	d.Duration = time.Duration(ms) * time.Millisecond
	return nil
}

// RatelimitState represents the state of a rate limit in a key.
type RatelimitState struct {
	Remaining    int64     `json:"remaining" db:"remaining"`
	LastRefilled time.Time `json:"last_refilled" db:"last_refilled"`
}

// Configured returns true if the rate limit is configured.
func (r *Ratelimit) Configured() bool {
	if r == nil {
		return false
	}
	return r.Limit > 0
}

// RatelimitPayload represents the payload for configuring a rate limit.
type RatelimitPayload struct {
	Limit          int64          `json:"limit" validate:"gte=0"`
	RefillRate     int64          `json:"refill_rate" validate:"gte=0"`
	RefillInterval Milliseconds `json:"refill_interval" validate:"gte=0"`
}

// Configured returns true if the rate limit is configured.
func (r *RatelimitPayload) Configured() bool {
	return r.Limit > 0
}

// Validate validates the rate limit payload.
func (r *RatelimitPayload) Validate(validator *validator.Validate) error {
	return validator.Struct(r)
}
