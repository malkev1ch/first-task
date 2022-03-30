package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/malkev1ch/first-task/internal/model"
)

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func (v *Validator) ValidateUpdateCat(input *model.UpdateCat) error {
	if input.Name != nil || input.DateBirth != nil || input.Vaccinated != nil {
		return nil
	}
	return errors.New("there must be at least one field in update method")
}
