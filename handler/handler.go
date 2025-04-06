package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yuugure-aikouka/kyoto-common/cache"
	"github.com/yuugure-aikouka/kyoto-common/config"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/mailer"
	"github.com/yuugure-aikouka/kyoto-common/model"
)

type Handler struct {
	config config.Config
	store  db.Store
	cache  cache.Store
	mailer mailer.Mailer
}

func NewHandler(
	config config.Config,
	store db.Store,
	cache cache.Store,
	mailer mailer.Mailer,
) *Handler {
	return &Handler{
		config: config,
		store:  store,
		cache:  cache,
		mailer: mailer,
	}
}

func jsonResponse[T any](c echo.Context, code int, data *T, errors ...any) error {
	return c.JSON(code, model.Response[T]{
		Status: http.StatusText(code),
		Data:   data,
		Errors: errors,
	})
}
