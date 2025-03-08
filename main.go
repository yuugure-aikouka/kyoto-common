package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yuugure-aikouka/kyoto-common/api"
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

	store := db.NewSQLStore(dbConn)
	mailer := mailer.NewMailer(cfg.EmailSender, mailer.SMTPConfig{
		Host:     cfg.SMTPHost,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
	})

	handler := handler.NewHandler(cfg, store, mailer)

	srv := api.NewServer(cfg, handler)

	log.Fatal(srv.Start())
}
