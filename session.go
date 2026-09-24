package main

import (
	"errors"
	"net/http"

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
