package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Asset struct {
	Name     string
	ID       int
	Creator  string
	ShopURL  string
	Memo     string
	Category string
}

// Showing Users' assets on /assets
func (s *server) listAssets(w http.ResponseWriter, req *http.Request) {
	userID, ok := userIDFromContext(req.Context())
	if !ok {
		fmt.Fprintf(os.Stderr, "user ID not found in context\n")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rows, err := s.pool.Query(req.Context(), "select asset_name, asset_id from assets where user_id = $1", userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		http.Error(w, "Query fail", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var assetName string
		var assetID int
		if err := rows.Scan(&assetName, &assetID); err != nil {
			fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			http.Error(w, "Scan failed", http.StatusInternalServerError)
			return
		}
		assets = append(assets, Asset{Name: assetName, ID: assetID})
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "rows failed: %v\n", err)
		http.Error(w, "db get fail", http.StatusInternalServerError)
		return
	}
	if err := s.tmpl.ExecuteTemplate(w, "assets.html", assets); err != nil {
		fmt.Fprintf(os.Stderr, "Template execution failed: %v\n", err)
	}
}

// Create assets using Form
func (s *server) createAsset(w http.ResponseWriter, req *http.Request) {
	userID, ok := userIDFromContext(req.Context())
	if !ok {
		fmt.Fprintf(os.Stderr, "user ID not found in context\n")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	assetName := req.FormValue("asset_name")
	creator := req.FormValue("creator")
	shopURL := req.FormValue("shop_url")
	memo := req.FormValue("memo")
	category := req.FormValue("asset_category")

	if assetName == "" || creator == "" || category == "" {
		http.Error(w, "アセット名・制作者・カテゴリは必須です", http.StatusBadRequest)
		return
	}

	_, err := s.pool.Exec(req.Context(),
		`insert into assets (user_id, asset_name, creator, shop_url, memo, asset_category)
			values ($1, $2, $3, $4, $5, $6)`, userID, assetName, creator, shopURL, memo, category)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
		http.Error(w, "db insert fail", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/assets", http.StatusSeeOther)
}

// Delete assets only owners if not, return 404
func (s *server) deleteAsset(w http.ResponseWriter, req *http.Request) {
	userID, ok := userIDFromContext(req.Context())
	if !ok {
		fmt.Fprintf(os.Stderr, "user ID not found in context\n")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	targetAssetIDTemp := req.PathValue("id")
	targetAssetID, err := strconv.Atoi(targetAssetIDTemp)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	result, err := s.pool.Exec(req.Context(),
		`DELETE FROM ASSETS WHERE asset_id = $1 AND user_id = $2`, targetAssetID, userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DELETE failed: %v\n", err)
		http.Error(w, "db DELETE fail", http.StatusInternalServerError)
		return
	}

	countDel := result.RowsAffected()
	if countDel == 0 {
		http.NotFound(w, req)
		return
	}
	http.Redirect(w, req, "/assets", http.StatusSeeOther)
}

// Update assets only owners if not, return 404
func (s *server) updateAsset(w http.ResponseWriter, req *http.Request) {
	userID, ok := userIDFromContext(req.Context())
	if !ok {
		fmt.Fprintf(os.Stderr, "user ID not found in context\n")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	targetAssetIDTemp := req.PathValue("id")
	targetAssetID, err := strconv.Atoi(targetAssetIDTemp)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	assetName := req.FormValue("asset_name")
	creator := req.FormValue("creator")
	shopURL := req.FormValue("shop_url")
	memo := req.FormValue("memo")
	category := req.FormValue("asset_category")

	if assetName == "" || creator == "" || category == "" {
		http.Error(w, "アセット名・制作者・カテゴリは必須です", http.StatusBadRequest)
		return
	}

	result, err := s.pool.Exec(req.Context(),
		`UPDATE ASSETS SET 
		asset_name = $1, 
		creator = $2, 
		shop_url = $3,
		memo = $4,
		asset_category = $5
		WHERE asset_id = $6 AND user_id = $7`,
		assetName, creator, shopURL, memo, category, targetAssetID, userID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "UPDATE failed: %v\n", err)
		http.Error(w, "db UPDATE fail", http.StatusInternalServerError)
		return
	}

	countUpdate := result.RowsAffected()
	if countUpdate == 0 {
		http.NotFound(w, req)
		return
	}
	http.Redirect(w, req, "/assets", http.StatusSeeOther)
}

func (s *server) editAssetForm(w http.ResponseWriter, req *http.Request) {
	userID, ok := userIDFromContext(req.Context())
	if !ok {
		fmt.Fprintf(os.Stderr, "user ID not found in context\n")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	targetAssetIDTemp := req.PathValue("id")
	targetAssetID, err := strconv.Atoi(targetAssetIDTemp)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	var asset Asset
	err = s.pool.QueryRow(req.Context(),
		`select asset_name, 
				creator, 
				COALESCE(shop_url, ''),
				COALESCE(memo, ''),
				asset_category from assets 
		WHERE asset_id = $1 AND user_id = $2`,
		targetAssetID, userID).Scan(
		&asset.Name, &asset.Creator, &asset.ShopURL, &asset.Memo, &asset.Category)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "自分が所有しているAssetではないか、存在しないAssetです。", http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Select failed: %v\n", err)
		http.Error(w, "db Select fail", http.StatusInternalServerError)
		return
	}

	asset.ID = targetAssetID
	if err := s.tmpl.ExecuteTemplate(w, "edit.html", asset); err != nil {
		fmt.Fprintf(os.Stderr, "Template execution failed: %v\n", err)
	}
}
