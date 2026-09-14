package main

import (
	"html/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type server struct {
	pool *pgxpool.Pool
	tmpl *template.Template
}
