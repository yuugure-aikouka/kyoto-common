package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewSQLStore(conn *pgxpool.Pool) Store {
	return SQLStore{
		connPool: conn,
		Queries:  New(conn),
	}
}
