package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
)

const GetById = "SELECT id, name FROM person where id=$1"

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=1234 host=localhost port=5432 dbname=snippet07 sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	var id int
	var name string
	err = pool.QueryRow(context.Background(), GetById, 1).Scan(&id, &name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(id, name)

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
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
