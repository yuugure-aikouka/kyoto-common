package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// https://echo.labstack.com/docs/request#validate-data
type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}

	return nil
}

func ValidateRequest(c echo.Context, i any) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}
