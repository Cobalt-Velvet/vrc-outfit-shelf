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

func TestLinkAssetOwner(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO ASSETS (user_id, asset_name, creator, asset_category) 
	VALUES (1, 'testAsset', 'Anonymous Kipfel', 'other')`)
	if err != nil {
		t.Fatalf("INSERT asset FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO AVATARS (user_id, avatar_name, creator) 
	VALUES (1, 'testAvatar', 'Anonymous Kipfel')`)
	if err != nil {
		t.Fatalf("INSERT avatar FAIL : %v", err)
	}
	req := httptest.NewRequest("POST", "/avatars/1/assets", strings.NewReader("asset_id=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), userIDKey{}, 1)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	srv.linkAsset(rec, req)
	want := http.StatusSeeOther
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}

	checkRow := 0
	err = srv.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM compatibility WHERE avatar_id = 1 AND asset_id = 1`).Scan(&checkRow)
	if err != nil {
		t.Fatalf("COUNT Query Fail : %v", err)
	}
	if checkRow != 1 {
		t.Errorf("comp row = %v, want %v", checkRow, 1)
	}

	wantLocation := "/avatars/1"
	if rec.Result().Header.Get("Location") != wantLocation {
		t.Errorf("location = %v, want %v", rec.Result().Header.Get("Location"), wantLocation)
	}
}

func TestLinkAssetNotOwnerAsset(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1'), ('testUser2', 'hoge2')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO ASSETS (user_id, asset_name, creator, asset_category) 
	VALUES (2, 'testAsset', 'Anonymous Kipfel', 'other')`)
	if err != nil {
		t.Fatalf("INSERT asset FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO AVATARS (user_id, avatar_name, creator) 
	VALUES (1, 'testAvatar', 'Anonymous Kipfel')`)
	if err != nil {
		t.Fatalf("INSERT avatar FAIL : %v", err)
	}
	req := httptest.NewRequest("POST", "/avatars/1/assets", strings.NewReader("asset_id=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), userIDKey{}, 1)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	srv.linkAsset(rec, req)
	want := http.StatusNotFound
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}

	checkRow := 1
	err = srv.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM compatibility WHERE avatar_id = 1 AND asset_id = 1`).Scan(&checkRow)
	if err != nil {
		t.Fatalf("COUNT Query Fail : %v", err)
	}
	if checkRow != 0 {
		t.Errorf("comp row = %v, want %v", checkRow, 0)
	}
}

func TestLinkAssetNotOwnerAvatar(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1'), ('testUser2', 'hoge2')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO ASSETS (user_id, asset_name, creator, asset_category) 
	VALUES (1, 'testAsset', 'Anonymous Kipfel', 'other')`)
	if err != nil {
		t.Fatalf("INSERT asset FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO AVATARS (user_id, avatar_name, creator) 
	VALUES (2, 'testAvatar', 'Anonymous Kipfel')`)
	if err != nil {
		t.Fatalf("INSERT avatar FAIL : %v", err)
	}
	req := httptest.NewRequest("POST", "/avatars/1/assets", strings.NewReader("asset_id=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), userIDKey{}, 1)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	srv.linkAsset(rec, req)
	want := http.StatusNotFound
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}

	checkRow := 1
	err = srv.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM compatibility WHERE avatar_id = 1 AND asset_id = 1`).Scan(&checkRow)
	if err != nil {
		t.Fatalf("COUNT Query Fail : %v", err)
	}
	if checkRow != 0 {
		t.Errorf("comp row = %v, want %v", checkRow, 0)
	}
}

func TestLinkAssetDuplicate(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO ASSETS (user_id, asset_name, creator, asset_category) 
	VALUES (1, 'testAsset', 'Anonymous Kipfel', 'other')`)
	if err != nil {
		t.Fatalf("INSERT asset FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO AVATARS (user_id, avatar_name, creator) 
	VALUES (1, 'testAvatar', 'Anonymous Kipfel')`)
	if err != nil {
		t.Fatalf("INSERT avatar FAIL : %v", err)
	}
	req1 := httptest.NewRequest("POST", "/avatars/1/assets", strings.NewReader("asset_id=1"))
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx1 := context.WithValue(req1.Context(), userIDKey{}, 1)
	req1 = req1.WithContext(ctx1)
	req1.SetPathValue("id", "1")
	rec1 := httptest.NewRecorder()

	req2 := httptest.NewRequest("POST", "/avatars/1/assets", strings.NewReader("asset_id=1"))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx2 := context.WithValue(req2.Context(), userIDKey{}, 1)
	req2 = req2.WithContext(ctx2)
	req2.SetPathValue("id", "1")
	rec2 := httptest.NewRecorder()

	srv.linkAsset(rec1, req1)
	want := http.StatusSeeOther
	if rec1.Code != want {
		t.Errorf("status = %v, want %v", rec1.Code, want)
	}

	checkRow := 0
	err = srv.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM compatibility WHERE avatar_id = 1 AND asset_id = 1`).Scan(&checkRow)
	if err != nil {
		t.Fatalf("COUNT Query Fail : %v", err)
	}
	if checkRow != 1 {
		t.Errorf("comp row = %v, want %v", checkRow, 1)
	}

	wantLocation := "/avatars/1"
	if rec1.Result().Header.Get("Location") != wantLocation {
		t.Errorf("location = %v, want %v", rec1.Result().Header.Get("Location"), wantLocation)
	}

	srv.linkAsset(rec2, req2)
	want = http.StatusSeeOther
	if rec2.Code != want {
		t.Errorf("status = %v, want %v", rec2.Code, want)
	}

	checkRow = 0
	err = srv.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM compatibility WHERE avatar_id = 1 AND asset_id = 1`).Scan(&checkRow)
	if err != nil {
		t.Fatalf("COUNT Query Fail : %v", err)
	}
	if checkRow != 1 {
		t.Errorf("comp row = %v, want %v", checkRow, 1)
	}

	wantLocation = "/avatars/1"
	if rec2.Result().Header.Get("Location") != wantLocation {
		t.Errorf("location = %v, want %v", rec2.Result().Header.Get("Location"), wantLocation)
	}
}
