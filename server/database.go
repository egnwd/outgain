package main

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"os"
)

func openDb() (*sql.DB, error) {
	url := os.Getenv("DATABASE_URL")
	connection, _ := pq.ParseURL(url)
	connection += " sslmode=require"

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println(err)
	}

	return db, err
}
