package driplimit

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is a validator for the Driplimit service. It implements the Driplimit interface.
type Validator struct {
	validator *validator.Validate
	driplimit Service
}

func NewServiceValidator(driplimit Service) *Validator {
	validator := validator.New()
	// https://github.com/go-playground/validator/issues/258#issuecomment-257281334
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return &Validator{
		validator: validator,
		driplimit: driplimit,
	}
}

// KeyCheck validates the payload and calls the KeyCheck method of the wrapped Driplimit service.
func (v *Validator) KeyCheck(ctx context.Context, payload KeysCheckPayload) (key *Key, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.KeyCheck(ctx, payload)
}

// KeyCreate validates the payload and calls the KeyCreate method of the wrapped Driplimit service.
func (v *Validator) KeyCreate(ctx context.Context, payload KeyCreatePayload) (key *Key, token *string, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, nil, err
	}
	return v.driplimit.KeyCreate(ctx, payload)
}

// KeyGet validates the payload and calls the KeyGet method of the wrapped Driplimit service.
func (v *Validator) KeyGet(ctx context.Context, payload KeyGetPayload) (key *Key, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.KeyGet(ctx, payload)
}

// KeyList validates the payload and calls the KeyList method of the wrapped Driplimit service.
func (v *Validator) KeyList(ctx context.Context, payload KeyListPayload) (klist *KeyList, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.KeyList(ctx, payload)
}

// KeyDelete validates the payload and calls the KeyDelete method of the wrapped Driplimit service.
func (v *Validator) KeyDelete(ctx context.Context, payload KeyDeletePayload) (err error) {
	if err := payload.Validate(v.validator); err != nil {
		return err
	}
	return v.driplimit.KeyDelete(ctx, payload)
}

// KeyspaceGet validates the payload and calls the KeyspaceGet method of the wrapped Driplimit service.
func (v *Validator) KeyspaceGet(ctx context.Context, payload KeyspaceGetPayload) (keyspace *Keyspace, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.KeyspaceGet(ctx, payload)
}

// KeyspaceCreate validates the payload and calls the KeyspaceCreate method of the wrapped Driplimit service.
func (v *Validator) KeyspaceCreate(ctx context.Context, payload KeyspaceCreatePayload) (keyspace *Keyspace, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.KeyspaceCreate(ctx, payload)
}

func (v *Validator) KeyspaceList(ctx context.Context, payload KeyspaceListPayload) (kslist *KeyspaceList, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.KeyspaceList(ctx, payload)
}

func (v *Validator) KeyspaceDelete(ctx context.Context, payload KeyspaceDeletePayload) (err error) {
	if err := payload.Validate(v.validator); err != nil {
		return err
	}
	return v.driplimit.KeyspaceDelete(ctx, payload)
}

func (v *Validator) ServiceKeyGet(ctx context.Context, payload ServiceKeyGetPayload) (sk *ServiceKey, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.ServiceKeyGet(ctx, payload)
}

func (v *Validator) ServiceKeyCreate(ctx context.Context, payload ServiceKeyCreatePayload) (sk *ServiceKey, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.ServiceKeyCreate(ctx, payload)
}

func (v *Validator) ServiceKeyList(ctx context.Context, payload ServiceKeyListPayload) (sklist *ServiceKeyList, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.ServiceKeyList(ctx, payload)
}

func (v *Validator) ServiceKeyDelete(ctx context.Context, payload ServiceKeyDeletePayload) (err error) {
	if err := payload.Validate(v.validator); err != nil {
		return err
	}
	return v.driplimit.ServiceKeyDelete(ctx, payload)
}

func (v *Validator) ServiceKeySetToken(ctx context.Context, payload ServiceKeySetTokenPayload) (err error) {
	if err := payload.Validate(v.validator); err != nil {
		return err
	}
	return v.driplimit.ServiceKeySetToken(ctx, payload)
}
