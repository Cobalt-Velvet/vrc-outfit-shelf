package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

var errNotLoggedIn = errors.New("not logged in")

func (s *server) currentUserID(req *http.Request) (int, error) {
	cookie, err := req.Cookie("session_id")
	if err != nil {
		return 0, errNotLoggedIn
	}
	var userID int
	err = s.pool.QueryRow(req.Context(),
		"select user_id from sessions where session_id = $1 and expires_at > now()",
		cookie.Value,
	).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errNotLoggedIn
	}
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *server) requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, err := s.currentUserID(req)
		if errors.Is(err, errNotLoggedIn) {
			http.Redirect(w, req, "/signin", http.StatusSeeOther)
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Session check failed: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, req)
	})
}
