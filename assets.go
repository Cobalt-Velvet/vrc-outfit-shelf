package main

import (
	"fmt"
	"net/http"
	"os"
)

type Asset struct {
	Name string
}

// -----------get asset name-----------
func (s *server) listAssets(w http.ResponseWriter, req *http.Request) {
	rows, err := s.pool.Query(req.Context(), "select asset_name from assets")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		http.Error(w, "Query fail", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var assetName string
		if err := rows.Scan(&assetName); err != nil {
			fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			http.Error(w, "Scan failed", http.StatusInternalServerError)
			return
		}
		assets = append(assets, Asset{Name: assetName})
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

// -----------post asset name-----------
func (s *server) createAsset(w http.ResponseWriter, req *http.Request) {
	assetName := req.FormValue("asset_name")
	creator := req.FormValue("creator")
	shopURL := req.FormValue("shop_url")
	memo := req.FormValue("memo")
	category := req.FormValue("asset_category")

	if assetName == "" || creator == "" || category == "" {
		http.Error(w, "Non-nullable field is NULL now", http.StatusBadRequest)
		return
	}

	_, err := s.pool.Exec(req.Context(),
		`insert into assets (user_id, asset_name, creator, shop_url, memo, asset_category)
			values (1, $1, $2, $3, $4, $5)`, assetName, creator, shopURL, memo, category)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
		http.Error(w, "db insert fail", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/assets", http.StatusSeeOther)
}
