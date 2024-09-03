package driplimit

import (
	"context"
	"errors"
)

// Service is the main driplimit service interface.
type Service interface {
	KeyCheck(ctx context.Context, payload KeysCheckPayload) (key *Key, err error)
	KeyCreate(ctx context.Context, payload KeyCreatePayload) (key *Key, err error)
	KeyGet(ctx context.Context, payload KeyGetPayload) (key *Key, err error)
	KeyList(ctx context.Context, payload KeyListPayload) (klist *KeyList, err error)
	KeyDelete(ctx context.Context, payload KeyDeletePayload) (err error)

	KeyspaceGet(ctx context.Context, payload KeyspaceGetPayload) (keyspace *Keyspace, err error)
	KeyspaceCreate(ctx context.Context, payload KeyspaceCreatePayload) (keyspace *Keyspace, err error)
	KeyspaceList(ctx context.Context, payload KeyspaceListPayload) (kslist *KeyspaceList, err error)
	KeyspaceDelete(ctx context.Context, payload KeyspaceDeletePayload) (err error)

	ServiceKeyGet(ctx context.Context, payload ServiceKeyGetPayload) (sk *ServiceKey, err error)
	ServiceKeyCreate(ctx context.Context, payload ServiceKeyCreatePayload) (sk *ServiceKey, err error)
	ServiceKeyList(ctx context.Context, payload ServiceKeyListPayload) (sklist *ServiceKeyList, err error)
	ServiceKeyDelete(ctx context.Context, payload ServiceKeyDeletePayload) (err error)
	ServiceKeySetToken(ctx context.Context, payload ServiceKeySetTokenPayload) (err error)
}

var (
	// ErrNotFound is returned when the requested item is not found.
	ErrNotFound = errors.New("not found")
	// ErrInvalidPayload is returned when the payload is invalid.
	ErrInvalidPayload = errors.New("invalid payload")
	// ErrInvalidExpiration is returned when the payload contains an empty or invalid expiration
	ErrInvalidExpiration = errors.New("invalid expiration")
	// ErrRateLimitExceeded is returned when the rate limit is exceeded.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	// ErrKeyExpired is returned when the key is expired.
	ErrKeyExpired = errors.New("key expired")
	// ErrUnauthorized is returned when the request is unauthorized.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrAlreadyExists is returned when the item already exists.
	ErrAlreadyExists = errors.New("already exists")
	// ErrCannotDeleteItself is returned when the item cannot delete itself.
	ErrCannotDeleteItself = errors.New("cannot delete itself")
)

// errHTTPCode is a map of HTTP status codes mapped to known errors.
var errHTTPCode = map[error]int{
	ErrInvalidPayload:     400,
	ErrUnauthorized:       401,
	ErrCannotDeleteItself: 403,
	ErrNotFound:           404,
	ErrAlreadyExists:      409,
	ErrKeyExpired:         419,
	ErrRateLimitExceeded:  429,
	ErrInvalidExpiration:  460,
}

// ErrItemNotFound is returned when the requested item is not found.
// It is a more precise wrapper around ErrNotFound.
type ErrItemNotFound string

// Error returns the error message. It implements the error interface.
func (e ErrItemNotFound) Error() string {
	return string(e) + " not found"
}

// Unwrap returns the wrapped error. It implements the errors.Wrapper interface.
func (e ErrItemNotFound) Unwrap() error {
	return ErrNotFound
}

// ErrItemAlreadyExists is returned when the requested item already exists.
// It is a more precise wrapper around ErrAlreadyExists.
type ErrItemAlreadyExists string

// Error returns the error message. It implements the error interface.
func (e ErrItemAlreadyExists) Error() string {
	return string(e) + " already exists"
}

// Unwrap returns the wrapped error. It implements the errors.Wrapper interface.
func (e ErrItemAlreadyExists) Unwrap() error {
	return ErrAlreadyExists
}

// ErrFromHTTPCode returns an error based on the given HTTP status code.
func ErrFromHTTPCode(code int) error {
	for err, cde := range errHTTPCode {
		if cde == code {
			return err
		}
	}
	return nil
}

// HTTPCodeFromErr returns an HTTP status code based on the given error.
func HTTPCodeFromErr(err error) int {
	for e, code := range errHTTPCode {
		if errors.Is(err, e) {
			return code
		}
	}
	return 500
}
