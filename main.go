package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yuugure-aikouka/kyoto-common/api"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
)

func main() {
	cfg := api.LoadConfig()
	ctx := context.Background()

	connPool, err := pgxpool.New(ctx, cfg.DBAddr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer connPool.Close()

	store := db.NewSQLStore(connPool)

	srv := api.NewServer(cfg, store)

	log.Fatal(srv.Start())
}
