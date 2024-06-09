package store

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// TimeNano is a wrapper around time.Time that serializes to
// and from int64 nanoseconds.
type TimeNano struct {
	time.Time
}

// Scan implements the sql.Scanner interface.
func (tn *TimeNano) Scan(v interface{}) error {
	tt, ok := v.(int64)
	if !ok {
		if v == nil {
			tn.Time = time.Time{}
			return nil
		}
		return fmt.Errorf("expected int64, got %T", v)
	}
	if tt == 0 {
		tn.Time = time.Time{}
		return nil
	}
	tn.Time = time.Unix(0, tt)
	return nil
}

// Value implements the driver.Valuer interface.
func (tn TimeNano) Value() (driver.Value, error) {
	if tn.IsZero() {
		return int64(0), nil
	}
	return tn.Time.UnixNano(), nil
}
