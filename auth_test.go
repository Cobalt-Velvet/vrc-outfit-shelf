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

func TestSigninEmpty(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/signin", nil)
	rec := httptest.NewRecorder()

	srv.signin(rec, req)
	want := http.StatusBadRequest
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
}

func TestSignoutNoCookie(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/signout", nil)
	rec := httptest.NewRecorder()

	srv.signout(rec, req)
	want := http.StatusSeeOther
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}

	wantLocation := "/signin"
	if rec.Result().Header.Get("Location") != wantLocation {
		t.Errorf("location = %v, want %v", rec.Result().Header.Get("Location"), wantLocation)
	}
}
