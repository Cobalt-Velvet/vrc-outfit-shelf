package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestSignupSuccess(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest("POST", "/signup", strings.NewReader("vrc_name=teststring&password=teststring"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	srv.signup(rec, req)
	want := http.StatusSeeOther
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}

	wantLocation := "/signin"
	if rec.Result().Header.Get("Location") != wantLocation {
		t.Errorf("location = %v, want %v", rec.Result().Header.Get("Location"), wantLocation)
	}
}
