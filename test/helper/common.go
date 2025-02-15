package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yuugure-aikouka/kyoto-common/model"
)

func ResetDB(dbConn *pgxpool.Pool) {
	ResetUsers(dbConn)
	ResetPartnerships(dbConn)
}

func ResetUsers(dbConn *pgxpool.Pool) {
	_, err := dbConn.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		log.Fatalf("Failed to reset users: %v", err)
	}
}

func ResetPartnerships(dbConn *pgxpool.Pool) {
	_, err := dbConn.Exec(context.Background(), "DELETE FROM partnerships")
	if err != nil {
		log.Fatalf("Failed to reset partnerships: %v", err)
	}
}

func UnmarshalResponseBody[T any](body *bytes.Buffer) (*model.Response[T], error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var responseBody model.Response[T]
	err = json.Unmarshal(bytes, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}
