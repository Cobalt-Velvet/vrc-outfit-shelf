package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// db connect
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Passed")
	defer conn.Close(context.Background())

	//--------------------------http server-------------------------
	//-----------get asset name-----------
	listAssets := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		rows, err := conn.Query(context.Background(), "select asset_name from assets")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
			http.Error(w, "db connect fail", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var assetName string
			if err := rows.Scan(&assetName); err != nil {
				fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
				fmt.Fprint(w, "Failed\n")
				continue
			}
			fmt.Fprintf(w, "<div>%s</div>", assetName)
		}
		if rows.Err() != nil {
			http.Error(w, "db get fail", http.StatusInternalServerError)
		}

		fmt.Fprint(w,
			`<form method="post" action="/assets">
category : <select name="asset_category">
<option>clothing</option>
<option>hair</option>
<option>accessory</option>
<option>prop</option>
<option>texture</option>
<option>tool</option>
<option>animation</option>
<option>other</option>
</select>
<label>asset_name: <input name="asset_name"></label><br>
<label>creator: <input name="creator"></label><br>
<label>shop_url: <input name="shop_url"></label><br>
<label>memo: <input name="memo"></label><br>
<input type="submit" value="submit"><br>
</form>`)
	}
	//-----------post asset name-----------
	createAsset := func(w http.ResponseWriter, req *http.Request) {
		assetName := req.FormValue("asset_name")
		creator := req.FormValue("creator")
		shopURL := req.FormValue("shop_url")
		memo := req.FormValue("memo")
		category := req.FormValue("asset_category")

		if assetName == "" || creator == "" || category == "" {
			http.Error(w, "Non nullable field is NULL now", http.StatusBadRequest)
			return
		}

		_, err := conn.Exec(context.Background(),
			"insert into assets (user_id, asset_name, creator, shop_url, memo, asset_category) values (1, $1, $2, $3, $4, $5)", assetName, creator, shopURL, memo, category)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
			http.Error(w, "db insert fail", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/assets", http.StatusSeeOther)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /assets", listAssets)
	mux.HandleFunc("POST /assets", createAsset)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
