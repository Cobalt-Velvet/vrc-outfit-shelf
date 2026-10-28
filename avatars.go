package main

import (
	"log/slog"
	"net/http"
)

type Avatar struct {
	ID      int
	Name    string
	Creator string
	ShopURL string
	Memo    string
}

func (s *server) listAvatars(w http.ResponseWriter, req *http.Request) {
	userID, ok := userIDFromContext(req.Context())
	if !ok {
		slog.Error("user ID not found in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rows, err := s.pool.Query(req.Context(), "select avatar_name, avatar_id from avatars where user_id = $1", userID)
	if err != nil {
		slog.Error("query failed", "err", err, "user_id", userID)
		http.Error(w, "Query fail", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var avatars []Avatar
	for rows.Next() {
		var avatarName string
		var avatarID int
		if err := rows.Scan(&avatarName, &avatarID); err != nil {
			slog.Error("scan failed", "err", err, "user_id", userID)
			http.Error(w, "Scan failed", http.StatusInternalServerError)
			return
		}
		avatars = append(avatars, Avatar{Name: avatarName, ID: avatarID})
	}
	if err := rows.Err(); err != nil {
		slog.Error("rows failed", "err", err, "user_id", userID)
		http.Error(w, "db get fail", http.StatusInternalServerError)
		return
	}
	if err := s.tmpl.ExecuteTemplate(w, "avatars.html", avatars); err != nil {
		slog.Error("template execution failed", "err", err, "user_id", userID)
	}
}

func (s *server) createAvatar(w http.ResponseWriter, req *http.Request) {
	userID, ok := userIDFromContext(req.Context())
	if !ok {
		slog.Error("user ID not found in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	avatarName := req.FormValue("avatar_name")
	creator := req.FormValue("creator")
	shopURL := req.FormValue("shop_url")
	memo := req.FormValue("memo")

	if avatarName == "" || creator == "" {
		http.Error(w, "アバター名・制作者は必須です", http.StatusBadRequest)
		return
	}

	_, err := s.pool.Exec(req.Context(),
		`insert into avatars (user_id, avatar_name, creator, shop_url, memo)
			values ($1, $2, $3, $4, $5)`, userID, avatarName, creator, shopURL, memo)
	if err != nil {
		slog.Error("insert failed", "err", err, "user_id", userID)
		http.Error(w, "db insert fail", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/avatars", http.StatusSeeOther)
}
