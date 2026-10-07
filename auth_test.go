package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestSignupEmptyName(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/signup", strings.NewReader("password=teststring"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	srv.signup(rec, req)
	want := http.StatusBadRequest
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
}

func TestSignupEmptyPassword(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/signup", strings.NewReader("vrc_name=teststring"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	srv.signup(rec, req)
	want := http.StatusBadRequest
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
}

func TestSigninEmptyName(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/signin", strings.NewReader("password=teststring"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	srv.signin(rec, req)
	want := http.StatusBadRequest
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
}

func TestSigninEmptyPassword(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/signin", strings.NewReader("vrc_name=teststring"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	srv.signin(rec, req)
	want := http.StatusBadRequest
	if rec.Code != want {
		t.Errorf("status = %v, want %v", rec.Code, want)
	}
}

func TestRequireLogin(t *testing.T) {
	srv := &server{}
	req := httptest.NewRequest("POST", "/signin", nil)
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
