package database

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"

	"fmt"
	"os"
	"strconv"
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

func UpdateLeaderboard(db *sql.DB, username string, score int) {
	// TODO: Query error checking

	deleteSingle := "DELETE FROM leaderboard WHERE ctid "
	deleteSingle += "IN (SELECT ctid FROM leaderboard ORDER BY "
	deleteSingle += "score asc LIMIT 1)"
	db.Query(deleteSingle)
	fmt.Println(deleteSingle)

	insertNew := "INSERT INTO leaderboard (username, score) "
	insertNew += "VALUES ('"
	insertNew += username
	insertNew += "', "
	insertNew += strconv.Itoa(score)
	insertNew += ")"
	db.Query(insertNew)
	fmt.Println(insertNew)
}

func GetMinScore(db *sql.DB) {
	// TODO: Query error checking
	db.Query("SELECT MIN(score) FROM leaderboard")
}
