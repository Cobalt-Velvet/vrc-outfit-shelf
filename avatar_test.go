package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAvatarSuccess(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	req := httptest.NewRequest("POST", "/avatars", strings.NewReader("avatar_name=testkipfel&creator=coolkipfel"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), userIDKey{}, 1)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	srv.createAvatar(rec, req)
	want := http.StatusSeeOther
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}

	wantLocation := "/avatars"
	if rec.Result().Header.Get("Location") != wantLocation {
		t.Errorf("location = %v, want %v", rec.Result().Header.Get("Location"), wantLocation)
	}

	var checkUserID int
	err = srv.pool.QueryRow(context.Background(), `
	SELECT user_id FROM avatars WHERE avatar_name = 'testkipfel'`).Scan(&checkUserID)
	wantID := 1
	if err != nil {
		t.Fatalf("SELECT avatar FAIL : %v", err)
	}

	if checkUserID != wantID {
		t.Errorf("user_id = %v, want %v", checkUserID, wantID)
	}
}

func TestCreateAvatarEmptyName(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/avatars", strings.NewReader("creator=coolkipfel"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), userIDKey{}, 1)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	srv.createAvatar(rec, req)
	want := http.StatusBadRequest
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
}
