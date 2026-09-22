package main

import (
	"context"
	"fmt"
	"html/template"
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
	defer pool.Close()

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read HTML: %v\n", err)
		os.Exit(1)
	}

	srv := &server{pool: pool, tmpl: tmpl}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /assets", srv.listAssets)
	mux.HandleFunc("POST /assets", srv.createAsset)
	mux.HandleFunc("GET /signup", srv.signupForm)
	mux.HandleFunc("POST /signup", srv.signup)
	mux.HandleFunc("GET /signin", srv.signinForm)
	mux.HandleFunc("POST /signin", srv.signin)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
