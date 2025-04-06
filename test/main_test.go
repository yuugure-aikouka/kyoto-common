package test

import (
	"context"
	"io"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/yuugure-aikouka/kyoto-common/api"
	cachepkg "github.com/yuugure-aikouka/kyoto-common/cache"
	"github.com/yuugure-aikouka/kyoto-common/config"
	storepkg "github.com/yuugure-aikouka/kyoto-common/db/store"
	"github.com/yuugure-aikouka/kyoto-common/handler"
	helper "github.com/yuugure-aikouka/kyoto-common/test/helper"
)

var dbConn *pgxpool.Pool
var store *storepkg.SQLStore

var rdbConn *redis.Client
var cache *cachepkg.RedisStore

var mailer *helper.MockMailer

var server *api.Server

func TestMain(m *testing.M) {
	cfg := config.LoadConfig()
	ctx := context.Background()

	var err error
	dbConn, err = pgxpool.New(ctx, cfg.DBAddr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	opt, err := redis.ParseURL(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Invalid redis URL: %v\n", err)
	}

	rdbConn = redis.NewClient(opt)
	defer rdbConn.Close()

	store = storepkg.NewSQLStore(dbConn).(*storepkg.SQLStore)
	cache = cachepkg.NewRedisStore(rdbConn).(*cachepkg.RedisStore)
	mailer = helper.NewMockMailer().(*helper.MockMailer)

	handler := handler.NewHandler(cfg, store, cache, mailer)
	server = api.NewServer(cfg, handler)
	server.Route().Logger.SetOutput(io.Discard)

	exitVal := m.Run()

	// clean up
	helper.ResetDB(dbConn)
	helper.ResetCache(rdbConn)

	os.Exit(exitVal)
}

func setupTest() {
	// to be run before each test
	helper.ResetDB(dbConn)
	helper.ResetCache(rdbConn)
	helper.ResetMailer(mailer)
}
