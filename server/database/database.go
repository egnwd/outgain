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
	GamesPlayed  int
	Achievements uint32 // Bitmap corresponding to locked/unlocked achievements
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
	row := instance.QueryRow(
		"SELECT * FROM achievements WHERE username='" + username + "'")
	var (
		totalScore   int
		highScore    int
		gamesPlayed  int
		achievements uint32
	)
	err := row.Scan(&username, &totalScore, &highScore, &gamesPlayed, &achievements)
	if err != nil {
		// Row does not exist, create one using username and default values
		achievements := strconv.FormatUint(uint64(achievements), 10)
		insert := "INSERT INTO achievements "
		insert += "(username, total_score, high_score, games_played, achievements) "
		insert += "VALUES ('" + username + "',"
		insert += "'" + strconv.Itoa(totalScore) + "',"
		insert += "'" + strconv.Itoa(highScore) + "',"
		insert += "'" + strconv.Itoa(gamesPlayed) + "',"
		insert += "'" + achievements + "')"
		_, err = instance.Exec(insert)
		NilCheck(err)
	}
	data := AchievementData{
		Username:     username,
		TotalScore:   totalScore,
		HighScore:    highScore,
		GamesPlayed:  gamesPlayed,
		Achievements: achievements,
	}
	return &data
}

func UpdateAchievements(data *AchievementData) {
	// Display bitmap as base 10 int, cast to uint64 to use function
	// TODO: check that this is being written and read correctly
	achievements := strconv.FormatUint(uint64(data.Achievements), 10)
	// Update row for current user
	update := "UPDATE achievements SET "
	update += "total_score='" + strconv.Itoa(data.TotalScore) + "',"
	update += "high_score='" + strconv.Itoa(data.HighScore) + "',"
	update += "games_played='" + strconv.Itoa(data.GamesPlayed) + "',"
	update += "achievements='" + achievements + "'"
	update += "WHERE username='" + data.Username + "'"
	_, err := instance.Exec(update)
	NilCheck(err)
}
