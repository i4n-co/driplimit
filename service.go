package driplimit

import (
	"context"
	"errors"
	"sync"
)

// Service is the main driplimit service interface.
type Service interface {
	KeyCheck(ctx context.Context, payload KeysCheckPayload) (key *Key, err error)
	KeyCreate(ctx context.Context, payload KeyCreatePayload) (key *Key, token *string, err error)
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
}

// ServiceWithToken is the service interface with a token.
type ServiceWithToken interface {
	WithToken(token string) Service
}

var (
	// ErrNotFound is returned when the requested item is not found.
	ErrNotFound = errors.New("not found")
	// ErrInvalidPayload is returned when the payload is invalid.
	ErrInvalidPayload = errors.New("invalid payload")
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
var errHTTPCode = map[int]error{
	400: ErrInvalidPayload,
	401: ErrUnauthorized,
	403: ErrCannotDeleteItself,
	404: ErrNotFound,
	409: ErrAlreadyExists,
	419: ErrKeyExpired,
	429: ErrRateLimitExceeded,
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

// invertedErrHTTPCode is a map of errors mapped to their HTTP status codes.
var invertedErrHTTPCode = new(sync.Map)

// ErrFromHTTPCode returns an error based on the given HTTP status code.
func ErrFromHTTPCode(code int) error {
	if err, ok := errHTTPCode[code]; ok {
		return err
	}
	return nil
}

// HTTPCodeFromErr returns an HTTP status code based on the given error.
func HTTPCodeFromErr(err error) int {
	entry, ok := invertedErrHTTPCode.Load(err)
	if ok {
		if code, ok := entry.(int); ok {
			return code
		}
	}
	for code, e := range errHTTPCode {
		if errors.Is(err, e) {
			invertedErrHTTPCode.Store(e, code)
			return code
		}
	}
	return 500
}
