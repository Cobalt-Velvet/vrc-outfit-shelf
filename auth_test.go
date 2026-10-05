package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignupEmpty(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/signup", nil)
	rec := httptest.NewRecorder()

	srv.signup(rec, req)
	want := http.StatusBadRequest
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
}
