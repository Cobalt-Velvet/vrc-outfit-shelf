package main

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestServer(t *testing.T) *server {

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("dbURL is empty")
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("Can not connect TestDB : %v", err)
	}
	t.Cleanup(pool.Close)

	_, err = pool.Exec(context.Background(), "TRUNCATE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Can not TRUNCATE TestDB : %v", err)
	}

	return &server{pool: pool}
}
