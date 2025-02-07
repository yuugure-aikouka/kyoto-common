package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/model"
)

type Handler struct {
	store db.Store
}

func NewHandler(store db.Store) *Handler {
	return &Handler{
		store: store,
	}
}

func jsonResponse[T any](c echo.Context, code int, data *T, errors ...any) error {
	return c.JSON(code, model.Response[T]{
		Status: http.StatusText(code),
		Data:   data,
		Errors: errors,
	})
}
