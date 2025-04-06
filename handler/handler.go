package handler

import (
	"github.com/yuugure-aikouka/kyoto-common/cache"
	"github.com/yuugure-aikouka/kyoto-common/config"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/mailer"
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
