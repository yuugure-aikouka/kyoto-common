package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	dbConn *pgxpool.Pool
}

func NewSQLStore(dbConn *pgxpool.Pool) Store {
	return &SQLStore{
		dbConn:  dbConn,
		Queries: New(dbConn),
	}
}
