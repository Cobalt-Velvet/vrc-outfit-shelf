package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var count int64
	err = conn.QueryRow(context.Background(), "select count(*) from assets").Scan(&count)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(count)

	httphandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "asdf\n")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", httphandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
