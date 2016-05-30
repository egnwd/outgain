package database

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/lib/pq"

	"log"
	"os"
)

func OpenDb() (*sql.DB, error) {
	url := os.Getenv("DATABASE_URL")

	url, err := pq.ParseURL(url)
	if err != nil {
		return nil, err
	}
	url += " sslmode=require"
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
}
