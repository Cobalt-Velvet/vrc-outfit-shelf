package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireLogin(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/assets", nil)
	rec := httptest.NewRecorder()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		called = true
	})

	testhandler := srv.requireLogin(next)
	testhandler.ServeHTTP(rec, req)

	want := http.StatusSeeOther
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}

	wantLocation := "/signin"
	if rec.Result().Header.Get("Location") != wantLocation {
		t.Errorf("location = %v, want %v", rec.Result().Header.Get("Location"), wantLocation)
	}

	if called {
		t.Errorf("next called")
	}

}
