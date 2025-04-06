package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/yuugure-aikouka/kyoto-common/api"
	"github.com/yuugure-aikouka/kyoto-common/cache"
	"github.com/yuugure-aikouka/kyoto-common/config"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/handler"
	"github.com/yuugure-aikouka/kyoto-common/mailer"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	dbConn, err := pgxpool.New(ctx, cfg.DBAddr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbConn.Close()

	opt, err := redis.ParseURL(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Invalid redis URL: %v\n", err)
	}

	rdb := redis.NewClient(opt)
	defer rdb.Close()

	store := db.NewSQLStore(dbConn)
	cache := cache.NewRedisStore(rdb)
	mailer := mailer.NewGoMailer(cfg.EmailSender, mailer.SMTPConfig{
		Host:     cfg.SMTPHost,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
	})

	handler := handler.NewHandler(cfg, store, cache, mailer)

	srv := api.NewServer(cfg, handler)

	log.Fatal(srv.Start())
}
