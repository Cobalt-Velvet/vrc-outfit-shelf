package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// db connect
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("unable to connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		slog.Error("unable to read HTML", "err", err)
		os.Exit(1)
	}

	srv := &server{pool: pool, tmpl: tmpl}

	mux := http.NewServeMux()

	mux.Handle("GET /assets", srv.requireLogin(http.HandlerFunc(srv.listAssets)))
	mux.Handle("POST /assets", srv.requireLogin(http.HandlerFunc(srv.createAsset)))
	mux.Handle("POST /assets/{id}/delete", srv.requireLogin(http.HandlerFunc(srv.deleteAsset)))
	mux.Handle("GET /assets/{id}/edit", srv.requireLogin(http.HandlerFunc(srv.editAssetForm)))
	mux.Handle("POST /assets/{id}/edit", srv.requireLogin(http.HandlerFunc(srv.updateAsset)))
	mux.Handle("GET /avatars", srv.requireLogin(http.HandlerFunc(srv.listAvatars)))
	mux.Handle("POST /avatars", srv.requireLogin(http.HandlerFunc(srv.createAvatar)))
	mux.Handle("GET /avatars/{id}", srv.requireLogin(http.HandlerFunc(srv.showAvatar)))
	mux.HandleFunc("GET /signup", srv.signupForm)
	mux.HandleFunc("POST /signup", srv.signup)
	mux.HandleFunc("GET /signin", srv.signinForm)
	mux.HandleFunc("POST /signin", srv.signin)
	mux.HandleFunc("POST /signout", srv.signout)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
