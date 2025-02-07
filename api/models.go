package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIResponse[T any] struct {
	Status string `json:"status"`
	Data   *T     `json:"data,omitempty"`
	Errors []any  `json:"errors,omitempty"`
}

func jsonResponse[T any](c echo.Context, code int, data *T, errors ...any) error {
	return c.JSON(code, APIResponse[T]{
		Status: http.StatusText(code),
		Data:   data,
		Errors: errors,
	})
}
