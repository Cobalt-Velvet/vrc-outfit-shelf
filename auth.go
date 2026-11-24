package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) index(w http.ResponseWriter, req *http.Request) {
	if err := s.tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		slog.Error("template execution failed", "err", err)
	}
}

func (s *server) signupForm(w http.ResponseWriter, req *http.Request) {
	if err := s.tmpl.ExecuteTemplate(w, "signup.html", nil); err != nil {
		slog.Error("template execution failed", "err", err)
	}
}

func (s *server) signup(w http.ResponseWriter, req *http.Request) {
	vrcName := req.FormValue("vrc_name")
	password := req.FormValue("password")

	if vrcName == "" || password == "" {
		http.Error(w, "ユーザー名とパスワードを入力してください", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("hashing failed", "err", err)
		http.Error(w, "signup failed", http.StatusInternalServerError)
		return
	}

	_, err = s.pool.Exec(req.Context(),
		"insert into users (vrc_name, password_hash) values ($1, $2)", vrcName, string(hashedPassword))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, "その名前は既に使われています", http.StatusConflict) //popup later
			return
		}
		slog.Error("insert failed", "err", err)
		http.Error(w, "db insert fail", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/signin", http.StatusSeeOther)
}

func (s *server) signinForm(w http.ResponseWriter, req *http.Request) {
	if err := s.tmpl.ExecuteTemplate(w, "signin.html", nil); err != nil {
		slog.Error("template execution failed", "err", err)
	}
}

func (s *server) signin(w http.ResponseWriter, req *http.Request) {

	vrcName := req.FormValue("vrc_name")
	password := req.FormValue("password")

	if vrcName == "" || password == "" {
		http.Error(w, "ユーザー名とパスワードを入力してください", http.StatusBadRequest)
		return
	}

	var idCheck int
	var hashCheck string

	err := s.pool.QueryRow(req.Context(),
		"select user_id, password_hash from users where vrc_name = $1",
		vrcName,
	).Scan(&idCheck, &hashCheck)

	// TODO: run bcrypt against a dummy hash when the user is not found,
	// so response time doesn't reveal whether the name is registered

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "ユーザー名またはパスワードが正しくありません", http.StatusUnauthorized)
		return
	}

	if err != nil {
		slog.Error("query failed", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashCheck), []byte(password)); err != nil {
		http.Error(w, "ユーザー名またはパスワードが正しくありません", http.StatusUnauthorized)
		return
	}

	//session
	err = s.createSession(w, req, idCheck)
	if err != nil {
		slog.Error("session create failed", "err", err, "user_id", idCheck)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/avatars", http.StatusSeeOther)
}

func (s *server) signout(w http.ResponseWriter, req *http.Request) {
	err := s.deleteSession(w, req)
	if errors.Is(err, errNotLoggedIn) {
		http.Redirect(w, req, "/signin", http.StatusSeeOther)
		return
	}
	if err != nil {
		slog.Error("session delete failed", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/signin", http.StatusSeeOther)
}
