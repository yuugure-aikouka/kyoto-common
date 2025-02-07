package test

import (
	"context"
	"io"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yuugure-aikouka/kyoto-common/api"
	"github.com/yuugure-aikouka/kyoto-common/config"
	db "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/handler"
	helper "github.com/yuugure-aikouka/kyoto-common/test/helper"
)

var dbConn *pgxpool.Pool

var server *api.Server

var store *db.SQLStore

func TestMain(m *testing.M) {
	cfg := config.LoadConfig()

	var err error
	dbConn, err = pgxpool.New(context.Background(), cfg.DBAddr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	store = db.NewSQLStore(dbConn).(*db.SQLStore)
	server = api.NewServer(cfg, handler.NewHandler(store))
	server.Route().Logger.SetOutput(io.Discard)

	exitVal := m.Run()

	helper.ResetDB(dbConn)

	os.Exit(exitVal)
}

func setupTest() {
	// to be run before each test
	helper.ResetDB(dbConn)
}
