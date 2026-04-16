// Package validator validate config
package validator

import (
	"net/url"

	"github.com/go-playground/validator/v10"
)

const validationAMPQURI = "amqpuri"

// Validator is a wrapper around the validator.Validate library for struct and field validation in Go.
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new Validator instance.
func NewValidator() *Validator {
	validate := validator.New()
	return &Validator{
		validator: validate,
	}
}

// Validate validates the given struct.
func (v *Validator) Validate(i any) error {
	err := v.registerValidationAMPQURI()
	if err != nil {
		return err
	}
	return v.validator.Struct(i)
}

func (v *Validator) registerValidationAMPQURI() error {
	err := v.validator.RegisterValidation(validationAMPQURI, func(fl validator.FieldLevel) bool {
		u, err := url.Parse(fl.Field().String())
		if err != nil {
			return false
		}

		if u.Scheme != "amqp" && u.Scheme != "amqps" {
			return false
		}

		if u.Host == "" {
			return false
		}

		if u.User == nil {
			return false
		}

		return true
	})
	if err != nil {
		return err
	}

	return nil
}
