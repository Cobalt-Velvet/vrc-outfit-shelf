package main

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Avatar struct {
	ID      int
	Name    string
	Creator string
	ShopURL string
	Memo    string
}

type avatarPage struct {
	Avatar Avatar
	Assets []Asset
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

func (s *server) showAvatar(w http.ResponseWriter, req *http.Request) {
	userID, ok := userIDFromContext(req.Context())
	if !ok {
		slog.Error("user ID not found in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	targetAvatarTemp := req.PathValue("id")
	targetAvatarID, err := strconv.Atoi(targetAvatarTemp)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	var avatar Avatar
	err = s.pool.QueryRow(req.Context(),
		`select avatar_name, 
				creator, 
				COALESCE(shop_url, ''),
				COALESCE(memo, '')
		from avatars WHERE avatar_id = $1 AND user_id = $2`,
		targetAvatarID, userID).Scan(
		&avatar.Name, &avatar.Creator, &avatar.ShopURL, &avatar.Memo)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "自分が所有しているAvatarではないか、存在しないAvatarです。", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("select failed", "err", err, "user_id", userID, "avatar_id", targetAvatarID)
		http.Error(w, "db Select fail", http.StatusInternalServerError)
		return
	}
	avatar.ID = targetAvatarID

	rows, err := s.pool.Query(req.Context(),
		`SELECT assets.asset_name, assets.asset_id
		FROM compatibility
		JOIN assets ON assets.asset_id = compatibility.asset_id
		WHERE compatibility.avatar_id = $1`,
		targetAvatarID)
	if err != nil {
		slog.Error("query failed", "err", err, "user_id", userID, "avatar_id", targetAvatarID)
		http.Error(w, "Query fail", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var assetName string
		var assetID int
		if err := rows.Scan(&assetName, &assetID); err != nil {
			slog.Error("scan failed", "err", err, "user_id", userID, "avatar_id", targetAvatarID)
			http.Error(w, "Scan failed", http.StatusInternalServerError)
			return
		}
		assets = append(assets, Asset{Name: assetName, ID: assetID})
	}
	if err := rows.Err(); err != nil {
		slog.Error("rows failed", "err", err, "user_id", userID, "avatar_id", targetAvatarID)
		http.Error(w, "db get fail", http.StatusInternalServerError)
		return
	}

	page := avatarPage{Avatar: avatar, Assets: assets}
	if err := s.tmpl.ExecuteTemplate(w, "avatar.html", page); err != nil {
		slog.Error("template execution failed", "err", err, "user_id", userID, "avatar_id", targetAvatarID)
	}
}
