package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/lib/pq"

	"log"
	"os"
)

func openDb() *sql.DB {
	url := os.Getenv("DATABASE_URL")

	//url := DATABASE_URL
	url, _ = pq.ParseURL(url)
	url += " sslmode=require"

	log.Println(url)

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Println(err)
	}

	return db
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
}
