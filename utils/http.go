package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yuugure-aikouka/kyoto-common/model"
)

func Response[T any](c echo.Context, code int, data *T, errors ...any) error {
	return c.JSON(code, model.Response[T]{
		Status: http.StatusText(code),
		Data:   data,
		Errors: errors,
	})
}
