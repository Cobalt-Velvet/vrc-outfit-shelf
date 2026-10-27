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

	mux.Handle("GET /assets", srv.requireLogin(http.HandlerFunc(srv.listAssets)))
	mux.Handle("POST /assets", srv.requireLogin(http.HandlerFunc(srv.createAsset)))
	mux.Handle("POST /assets/{id}/delete", srv.requireLogin(http.HandlerFunc(srv.deleteAsset)))
	mux.Handle("GET /assets/{id}/edit", srv.requireLogin(http.HandlerFunc(srv.editAssetForm)))
	mux.Handle("POST /assets/{id}/edit", srv.requireLogin(http.HandlerFunc(srv.updateAsset)))
	mux.HandleFunc("GET /signup", srv.signupForm)
	mux.HandleFunc("POST /signup", srv.signup)
	mux.HandleFunc("GET /signin", srv.signinForm)
	mux.HandleFunc("POST /signin", srv.signin)
	mux.HandleFunc("POST /signout", srv.signout)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
