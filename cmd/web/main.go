package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"se07/pkg/models/postgres"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgres.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=1234 host=localhost port=5435 dbname=snippet07 sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}

	defer pool.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &postgres.SnippetModel{Pool: pool},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
