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
func TestSigninSuccess(t *testing.T) {
	srv := newTestServer(t)
	req1 := httptest.NewRequest("POST", "/signup", strings.NewReader("vrc_name=teststring&password=teststring"))
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec1 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/signin", strings.NewReader("vrc_name=teststring&password=teststring"))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec2 := httptest.NewRecorder()

	srv.signup(rec1, req1)
	want := http.StatusSeeOther
	if rec1.Code != want {
		t.Fatalf("signup status = %v, want %v", rec1.Code, want)
	}
	srv.signin(rec2, req2)
	want = http.StatusSeeOther
	if rec2.Code != want {
		t.Errorf("signin status = %v, want %v", rec2.Code, want)
	}

	wantLocation := "/assets"
	if rec2.Result().Header.Get("Location") != wantLocation {
		t.Errorf("location = %v, want %v", rec2.Result().Header.Get("Location"), wantLocation)
	}

	flag := false
	for _, cookie := range rec2.Result().Cookies() {
		if cookie.Name == "session_id" {
			flag = true

			if cookie.Value == "" {
				t.Error("session_id cookie value is empty")
			}
			break
		}
	}
	if !flag {
		t.Error("session_id cookie not found")
	}

}

func TestSigninUnknownUser(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest("POST", "/signin", strings.NewReader("vrc_name=teststring&password=teststring"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	srv.signin(rec, req)
	want := http.StatusUnauthorized
	if rec.Code != want {
		t.Fatalf("signin status = %v, want %v", rec.Code, want)
	}
}

func TestSigninWrongPassword(t *testing.T) {
	srv := newTestServer(t)
	req1 := httptest.NewRequest("POST", "/signup", strings.NewReader("vrc_name=teststring&password=teststring"))
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec1 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/signin", strings.NewReader("vrc_name=teststring&password=testwrongstring"))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec2 := httptest.NewRecorder()

	srv.signup(rec1, req1)
	want := http.StatusSeeOther
	if rec1.Code != want {
		t.Fatalf("signup status = %v, want %v", rec1.Code, want)
	}
	srv.signin(rec2, req2)
	want = http.StatusUnauthorized
	if rec2.Code != want {
		t.Errorf("signin status = %v, want %v", rec2.Code, want)
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
