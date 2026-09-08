package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// db connect
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Passed")
	defer pool.Close()

	srv := &server{pool: pool}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /assets", srv.listAssets)
	mux.HandleFunc("POST /assets", srv.createAsset)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
