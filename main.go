package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yuugure-aikouka/kyoto-common/api"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/utils"
)

func main() {
	cfg := utils.LoadConfig()
	ctx := context.Background()

	dbConn, err := pgxpool.New(ctx, cfg.DBAddr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbConn.Close()

	store := db.NewSQLStore(dbConn)

	srv := api.NewServer(cfg, store)

	log.Fatal(srv.Start())
}
