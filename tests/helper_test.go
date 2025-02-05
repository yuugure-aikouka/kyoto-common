package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/yuugure-aikouka/kyoto-common/api"
)

func resetDB() {
	resetUsers()
	resetPartnerships()
}

func resetUsers() {
	_, err := dbConn.Exec(context.Background(), "DELETE FROM users")
	if err != nil {
		log.Fatalf("Failed to reset users: %v", err)
	}
}

func resetPartnerships() {
	_, err := dbConn.Exec(context.Background(), "DELETE FROM partnerships")
	if err != nil {
		log.Fatalf("Failed to reset partnerships: %v", err)
	}
}

func unmarshalResponseBody[T any](body *bytes.Buffer) (*api.APIResponse[T], error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var responseBody api.APIResponse[T]
	err = json.Unmarshal(bytes, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}
