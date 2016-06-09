package database

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"

	"fmt"
	"log"
	"os"
	"strconv"
)

var instance *sql.DB

type Leaderboard struct {
	Usernames []string
	Scores    []int
}

func OpenDb() error {
	url := os.Getenv("DATABASE_URL")
	url, err := pq.ParseURL(url)
	NilCheck(err)
	url += " sslmode=require"
	instance, err = sql.Open("postgres", url)
	NilCheck(err)
	return nil
}

func UpdateLeaderboard(username string, score int) {
	// TODO: Query error checking
	deleteSingle := `DELETE FROM leaderboard WHERE ctid 
	                 IN (SELECT ctid FROM leaderboard ORDER BY 
                         score asc LIMIT 1)`
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
	NilCheck(err)
	defer rows.Close()
	rows.Next()
	var score int
	err = rows.Scan(&score)
	NilCheck(err)
	return score
}

func NilCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllRows() *Leaderboard {
	usernames := make([]string, 10)
	scores := make([]int, 10)
	rows, err := instance.Query("SELECT * FROM leaderboard ORDER BY score desc")
	NilCheck(err)
	defer rows.Close()
	i := 0
	var (
		username string
		score    int
	)
	for rows.Next() && i < 10 {
		err = rows.Scan(&username, &score)
		NilCheck(err)
		usernames[i] = username
		scores[i] = score
		i += 1
	}
	leaderboard := Leaderboard{
		Usernames: usernames,
		Scores:    scores,
	}
	fmt.Println(leaderboard)
	return &leaderboard
}
