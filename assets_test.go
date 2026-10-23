package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteAssetNotOwner(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1'), ('testUser2', 'hoge')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO ASSETS (user_id, asset_name, creator, asset_category) 
	VALUES (1, 'testAsset', 'Anonymous Kipfel', 'other')`)
	if err != nil {
		t.Fatalf("INSERT asset FAIL : %v", err)
	}

	req := httptest.NewRequest("POST", "/assets/1/delete", nil)
	req.SetPathValue("id", "1")
	ctx := context.WithValue(req.Context(), userIDKey{}, 2)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	srv.deleteAsset(rec, req)
	want := http.StatusNotFound
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
	countRow := 0
	err = srv.pool.QueryRow(context.Background(), `SELECT count(*) FROM assets WHERE asset_id = 1`).Scan(&countRow)
	if err != nil {
		t.Fatalf("SCAN asset FAIL : %v", err)
	}
	if countRow != 1 {
		t.Errorf("count = %v, want %v", countRow, 1)
	}

}

func TestDeleteAssetOwner(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1'), ('testUser2', 'hoge')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO ASSETS (user_id, asset_name, creator, asset_category) 
	VALUES (1, 'testAsset', 'Anonymous Kipfel', 'other')`)
	if err != nil {
		t.Fatalf("INSERT asset FAIL : %v", err)
	}

	req := httptest.NewRequest("POST", "/assets/1/delete", nil)
	req.SetPathValue("id", "1")
	ctx := context.WithValue(req.Context(), userIDKey{}, 1)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	srv.deleteAsset(rec, req)
	want := http.StatusSeeOther
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
	countRow := 0
	err = srv.pool.QueryRow(context.Background(), `SELECT count(*) FROM assets WHERE asset_id = 1`).Scan(&countRow)
	if err != nil {
		t.Fatalf("SCAN asset FAIL : %v", err)
	}
	if countRow != 0 {
		t.Errorf("count = %v, want %v", countRow, 0)
	}

}

func TestUpdateAssetOwner(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1'), ('testUser2', 'hoge')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO ASSETS (user_id, asset_name, creator, asset_category) 
	VALUES (1, 'testAsset', 'Anonymous Kipfel', 'other')`)
	if err != nil {
		t.Fatalf("INSERT asset FAIL : %v", err)
	}

	req := httptest.NewRequest("POST", "/assets/1/edit", strings.NewReader("asset_name=NewKipfel!&creator=Anonymous Kipfel&asset_category=other"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "1")
	ctx := context.WithValue(req.Context(), userIDKey{}, 1)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	srv.updateAsset(rec, req)
	want := http.StatusSeeOther
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
	checkRow := ""
	err = srv.pool.QueryRow(context.Background(), `SELECT asset_name FROM assets WHERE asset_id = 1`).Scan(&checkRow)
	if err != nil {
		t.Fatalf("SCAN asset FAIL : %v", err)
	}
	if checkRow != "NewKipfel!" {
		t.Errorf("asset_name = %v, want %v", checkRow, "NewKipfel!")
	}
}

func TestUpdateAssetNotOwner(t *testing.T) {
	srv := newTestServer(t)
	_, err := srv.pool.Exec(context.Background(), `
	INSERT INTO USERS (vrc_name, password_hash) 
	VALUES ('testUser1', 'hoge1'), ('testUser2', 'hoge')`)
	if err != nil {
		t.Fatalf("INSERT users FAIL : %v", err)
	}
	_, err = srv.pool.Exec(context.Background(), `
	INSERT INTO ASSETS (user_id, asset_name, creator, asset_category) 
	VALUES (1, 'testAsset', 'Anonymous Kipfel', 'other')`)
	if err != nil {
		t.Fatalf("INSERT asset FAIL : %v", err)
	}

	req := httptest.NewRequest("POST", "/assets/1/edit", strings.NewReader("asset_name=NewKipfel!&creator=Anonymous Kipfel&asset_category=other"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "1")
	ctx := context.WithValue(req.Context(), userIDKey{}, 2)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	srv.updateAsset(rec, req)
	want := http.StatusNotFound
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
	checkRow := ""
	err = srv.pool.QueryRow(context.Background(), `SELECT asset_name FROM assets WHERE asset_id = 1`).Scan(&checkRow)
	if err != nil {
		t.Fatalf("SCAN asset FAIL : %v", err)
	}
	if checkRow != "testAsset" {
		t.Errorf("asset_name = %v, want %v", checkRow, "testAsset")
	}
}
