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

	listAssets := func(w http.ResponseWriter, req *http.Request) {
		//---------------get asset name-----------
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
			fmt.Fprintf(w, "%s\n", assetName)
		}
		if rows.Err() != nil {
			http.Error(w, "db get fail", http.StatusInternalServerError)
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", listAssets)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
