package main

import "github.com/jackc/pgx/v5/pgxpool"

type server struct {
	pool *pgxpool.Pool
}
