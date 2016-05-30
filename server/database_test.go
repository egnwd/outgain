package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDatabase(t *testing.T) {
	db, err := openDb()
	require.Nil(t, err, "error should be nil")
	_, err = db.Query("CREATE TABLE leaderboardTest(username text, score int)")
	require.Nil(t, err, "error should be nil")
	_, err = db.Query("INSERT INTO leaderboardTest (username, score) VALUES ('plietar', -14)")
	require.Nil(t, err, "error should be nil")
	rows, err := db.Query("SELECT * FROM leaderboardTest")
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
	assert.Equal(t, score, -12, "score is set to -14 in the database setup")

	err = rows.Err()
	assert.Nil(t, err, "error should be nil")

	_, err = db.Query("DROP TABLE leaderboardTest")
	assert.Nil(t, err, "error should be nil")
}
