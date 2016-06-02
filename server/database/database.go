package database

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"

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
