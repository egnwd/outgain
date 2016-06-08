package database

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"

	"log"
	"os"
	"strconv"
)

var instance *sql.DB

func OpenDb() error {
	url := os.Getenv("DATABASE_URL")
	url, err := pq.ParseURL(url)
	nilCheck(err)
	url += " sslmode=require"
	instance, err = sql.Open("postgres", url)
	nilCheck(err)
	return nil
}

func UpdateLeaderboard(username string, score int) {
	// TODO: Query error checking
	deleteSingle := "DELETE FROM leaderboard WHERE ctid "
	deleteSingle += "IN (SELECT ctid FROM leaderboard ORDER BY "
	deleteSingle += "score asc LIMIT 1)"
	instance.Query(deleteSingle)

	insertNew := "INSERT INTO leaderboard (username, score) "
	insertNew += "VALUES ('"
	insertNew += username
	insertNew += "', "
	insertNew += strconv.Itoa(score)
	insertNew += ")"
	instance.Query(insertNew)
}

func GetMinScore() int {
	// TODO: Query error checking
	rows, err := instance.Query("SELECT MIN(score) FROM leaderboard")
	nilCheck(err)
	defer rows.Close()
	rows.Next()
	var score int
	err = rows.Scan(&score)
	nilCheck(err)
	return score
}

func nilCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
