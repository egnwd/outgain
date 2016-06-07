package database

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"os"
	"testing"
)

var testInstance *sql.DB

func OpenTestDb() error {
	url := os.Getenv("DATABASE_URL")
	url, err := pq.ParseURL(url)
	nilCheck(err)
	url += " sslmode=require"
	instance, err = sql.Open("postgres", url)
	nilCheck(err)
	return nil
}

func TestDatabase(t *testing.T) {
	err := OpenTestDb()
	require.Nil(t, err, "error should be nil")
	_, err = instance.Query("CREATE TABLE leaderboardTest(username text, score int)")
	require.Nil(t, err, "error should be nil")
	_, err = instance.Query("INSERT INTO leaderboardTest (username, score) VALUES ('plietar', -14)")
	require.Nil(t, err, "error should be nil")
	rows, err := instance.Query("SELECT * FROM leaderboardTest")
	require.Nil(t, err, "error should be nil")
	var (
		username string
		score    int
	)
	rows.Next()
	defer rows.Close()
	err = rows.Scan(&username, &score)
	assert.Nil(t, err, "error should be nil")

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, username, "plietar", "username is set to plietar in the database setup")
	assert.Equal(t, score, -14, "score is set to -14 in the database setup")

	err = rows.Err()
	assert.Nil(t, err, "error should be nil")

	_, err = instance.Query("DROP TABLE leaderboardTest")
	assert.Nil(t, err, "error should be nil")
}
