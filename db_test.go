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

func TestSignupDuplicate(t *testing.T) {
	srv := newTestServer(t)
	req1 := httptest.NewRequest("POST", "/signup", strings.NewReader("vrc_name=teststring&password=teststring"))
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec1 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/signup", strings.NewReader("vrc_name=teststring&password=teststring"))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec2 := httptest.NewRecorder()

	srv.signup(rec1, req1)
	want := http.StatusSeeOther
	if rec1.Code != want {
		t.Fatalf("first signup = %v, want %v", rec1.Code, want)
	}
	srv.signup(rec2, req2)
	want = http.StatusConflict
	if rec2.Code != want {
		t.Errorf("second signup = %v, want %v", rec2.Code, want)
	}
}
