package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

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
