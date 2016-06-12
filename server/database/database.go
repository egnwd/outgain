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

type Leaderboard struct {
	Usernames []string
	Scores    []int
}

type AchievementData struct {
	Username     string
	TotalScore   int
	HighScore    int
	RoundsPlayed int
	Achievements uint64 // Bitmap corresponding to locked/unlocked achievements
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
	_, err := instance.Exec(deleteSingle)
	NilCheck(err)
	insertNew := "INSERT INTO leaderboard (username, score) "
	insertNew += "VALUES ('"
	insertNew += username
	insertNew += "', "
	insertNew += strconv.Itoa(score)
	insertNew += ")"
	_, err = instance.Exec(insertNew)
	NilCheck(err)
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
	return &leaderboard
}

func GetAchievements(username string) *AchievementData {
	row, err := instance.Query(
		"SELECT * FROM achievements WHERE username='" + username + "'")
	NilCheck(err)
	defer row.Close()
	var (
		totalScore   int
		highScore    int
		roundsPlayed int
		achievements uint64
	)
	err = row.Scan(&username, &totalScore, &highScore, &roundsPlayed, &achievements)
	NilCheck(err)
	data := AchievementData{
		Username:     username,
		TotalScore:   totalScore,
		HighScore:    highScore,
		RoundsPlayed: roundsPlayed,
		Achievements: achievements,
	}
	return &data
}

func UpdateAchievements(data *AchievementData) {
	// TODO: Query error checking
	// TODO: update database row with new data
}
