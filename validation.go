package driplimit

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// driplimitValidator is a validator for the Driplimit service. It implements the Driplimit interface.
type driplimitValidator struct {
	validator *validator.Validate
	driplimit ServiceWithToken
}

func NewServiceValidator(driplimit ServiceWithToken) ServiceWithToken {
	validator := validator.New()
	// https://github.com/go-playground/validator/issues/258#issuecomment-257281334
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return &driplimitValidator{
		validator: validator,
		driplimit: driplimit,
	}
}

type driplimitValidatorWithToken struct {
	*driplimitValidator
	token string
}

func (v *driplimitValidator) WithToken(token string) Service {
	return &driplimitValidatorWithToken{
		driplimitValidator: v,
		token:              token,
	}
}

// KeyCheck validates the payload and calls the KeyCheck method of the wrapped Driplimit service.
func (v *driplimitValidatorWithToken) KeyCheck(ctx context.Context, payload KeysCheckPayload) (key *Key, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).KeyCheck(ctx, payload)
}

// KeyCreate validates the payload and calls the KeyCreate method of the wrapped Driplimit service.
func (v *driplimitValidatorWithToken) KeyCreate(ctx context.Context, payload KeyCreatePayload) (key *Key, token *string, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, nil, err
	}
	return v.driplimit.WithToken(v.token).KeyCreate(ctx, payload)
}

// KeyGet validates the payload and calls the KeyGet method of the wrapped Driplimit service.
func (v *driplimitValidatorWithToken) KeyGet(ctx context.Context, payload KeyGetPayload) (key *Key, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).KeyGet(ctx, payload)
}

// KeyList validates the payload and calls the KeyList method of the wrapped Driplimit service.
func (v *driplimitValidatorWithToken) KeyList(ctx context.Context, payload KeyListPayload) (klist *KeyList, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).KeyList(ctx, payload)
}

// KeyDelete validates the payload and calls the KeyDelete method of the wrapped Driplimit service.
func (v *driplimitValidatorWithToken) KeyDelete(ctx context.Context, payload KeyDeletePayload) (err error) {
	if err := payload.Validate(v.validator); err != nil {
		return err
	}
	return v.driplimit.WithToken(v.token).KeyDelete(ctx, payload)
}

// KeyspaceGet validates the payload and calls the KeyspaceGet method of the wrapped Driplimit service.
func (v *driplimitValidatorWithToken) KeyspaceGet(ctx context.Context, payload KeyspaceGetPayload) (keyspace *Keyspace, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).KeyspaceGet(ctx, payload)
}

// KeyspaceCreate validates the payload and calls the KeyspaceCreate method of the wrapped Driplimit service.
func (v *driplimitValidatorWithToken) KeyspaceCreate(ctx context.Context, payload KeyspaceCreatePayload) (keyspace *Keyspace, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).KeyspaceCreate(ctx, payload)
}

func (v *driplimitValidatorWithToken) KeyspaceList(ctx context.Context, payload KeyspaceListPayload) (kslist *KeyspaceList, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).KeyspaceList(ctx, payload)
}

func (v *driplimitValidatorWithToken) KeyspaceDelete(ctx context.Context, payload KeyspaceDeletePayload) (err error) {
	if err := payload.Validate(v.validator); err != nil {
		return err
	}
	return v.driplimit.WithToken(v.token).KeyspaceDelete(ctx, payload)
}

func (v *driplimitValidatorWithToken) ServiceKeyGet(ctx context.Context, payload ServiceKeyGetPayload) (sk *ServiceKey, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).ServiceKeyGet(ctx, payload)
}

func (v *driplimitValidatorWithToken) ServiceKeyCreate(ctx context.Context, payload ServiceKeyCreatePayload) (sk *ServiceKey, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).ServiceKeyCreate(ctx, payload)
}

func (v *driplimitValidatorWithToken) ServiceKeyList(ctx context.Context, payload ServiceKeyListPayload) (sklist *ServiceKeyList, err error) {
	if err := payload.Validate(v.validator); err != nil {
		return nil, err
	}
	return v.driplimit.WithToken(v.token).ServiceKeyList(ctx, payload)
}

func (v *driplimitValidatorWithToken) ServiceKeyDelete(ctx context.Context, payload ServiceKeyDeletePayload) (err error) {
	if err := payload.Validate(v.validator); err != nil {
		return err
	}
	return v.driplimit.WithToken(v.token).ServiceKeyDelete(ctx, payload)
}
