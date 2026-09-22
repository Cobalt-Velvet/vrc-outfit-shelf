package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) signupForm(w http.ResponseWriter, req *http.Request) {
	if err := s.tmpl.ExecuteTemplate(w, "signup.html", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Template execution failed: %v\n", err)
	}
}

func (s *server) signup(w http.ResponseWriter, req *http.Request) {
	vrcName := req.FormValue("vrc_name")
	password := req.FormValue("password")

	if vrcName == "" || password == "" {
		http.Error(w, "Non-nullable field is NULL now", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hashing failed: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
		http.Error(w, "db insert fail", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/assets", http.StatusSeeOther)
}

func (s *server) signinForm(w http.ResponseWriter, req *http.Request) {
	if err := s.tmpl.ExecuteTemplate(w, "signin.html", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Template execution failed: %v\n", err)
	}
}

func (s *server) signin(w http.ResponseWriter, req *http.Request) {

	vrcName := req.FormValue("vrc_name")
	password := req.FormValue("password")

	if vrcName == "" || password == "" {
		http.Error(w, "Non-nullable field is NULL now", http.StatusBadRequest)
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
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashCheck), []byte(password)); err != nil {
		http.Error(w, "ユーザー名またはパスワードが正しくありません", http.StatusUnauthorized)
		return
	}

	//session
	sessionID := rand.Text()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.pool.Exec(req.Context(),
		"insert into sessions (session_id, user_id, expires_at) values ($1, $2, $3)",
		sessionID, idCheck, expiresAt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Session insert failed: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, req, "/assets", http.StatusSeeOther)
}
