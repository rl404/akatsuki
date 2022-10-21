// Package playground is a wrapper of the original "github.com/go-playground/validator"
// and "github.com/go-playground/mold" library.
//
// All of the validation and modifier tags from both libraries
// can still be used here.
package playground

import (
	"context"
	"errors"
	"reflect"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
)

// List of errors.
var (
	ErrRequiredName = errors.New("required name")
	ErrRequiredFn   = errors.New("required function")
	ErrUnknownType  = errors.New("unknown type")
)

// Validator is validator and modifier client.
type Validator struct {
	mod        *mold.Transformer
	val        *validator.Validate
	autoMod    bool
	customErrs map[string]func(string, ...string) error
}

// New to create new validator and modifier.
// Pass true if you want to modify the data
// automatically before validate.
func New(autoMod bool) *Validator {
	return &Validator{
		mod:        modifiers.New(),
		val:        validator.New(),
		autoMod:    autoMod,
		customErrs: make(map[string]func(string, ...string) error),
	}
}

// RegisterModifier to register custom modifier.
func (v *Validator) RegisterModifier(name string, fn func(string) string) error {
	if name == "" {
		return ErrRequiredName
	}

	if fn == nil {
		return ErrRequiredFn
	}

	v.mod.Register(name, func(ctx context.Context, fl mold.FieldLevel) error {
		switch fl.Field().Kind() {
		case reflect.String:
			fl.Field().SetString(fn(fl.Field().String()))
		}
		return nil
	})

	return nil
}

// Modify to modify/set struct field value according
// to modifier tag. Param `data` should be a pointer.
func (v *Validator) Modify(data interface{}) error {
	return v.mod.Struct(context.Background(), data)
}

// RegisterValidator to register custom validator.
func (v *Validator) RegisterValidator(name string, fn func(interface{}, ...string) bool) error {
	if name == "" {
		return ErrRequiredName
	}

	if fn == nil {
		return ErrRequiredFn
	}

	return v.val.RegisterValidation(name, func(fl validator.FieldLevel) bool {
		var param []string
		if fl.Param() != "" {
			param = append(param, fl.Param())
		}
		return fn(fl.Field().Interface(), param...)
	})
}

// RegisterValidatorError to register custom error message handling.
func (v *Validator) RegisterValidatorError(name string, fn func(string, ...string) error) error {
	if name == "" {
		return ErrRequiredName
	}

	if fn == nil {
		return ErrRequiredFn
	}

	v.customErrs[name] = fn

	return nil
}

// Validate to validate struct field value according
// to validator tag. Will modify/set the data if you
// turn on the autoMod. Param `data` should be a pointer.
func (v *Validator) Validate(data interface{}) error {
	if v.autoMod {
		if err := v.Modify(data); err != nil {
			return err
		}
	}

	if err := v.val.Struct(data); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		if len(errs) > 0 {
			fn, ok := v.customErrs[errs[0].Tag()]
			if !ok {
				return err
			}

			var param []string
			if errs[0].Param() != "" {
				param = append(param, errs[0].Param())
			}

			return fn(errs[0].Field(), param...)
		}

		return err
	}

	return nil
}
