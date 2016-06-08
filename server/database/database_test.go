package database

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"os"
	"strconv"
	"testing"
)

var testInstance *sql.DB

func OpenTestDb() error {
	url := os.Getenv("DATABASE_URL")
	url, err := pq.ParseURL(url)
	if err != nil {
		return err
	}
	url += " sslmode=require"
	testInstance, err = sql.Open("postgres", url)
	nilCheck(err)
	if err != nil {
		return err
	}
	return nil
}

func setUpGeneric() error {
	err := OpenTestDb()
	if err != nil {
		return err
	}
	_, err = testInstance.Query("CREATE TABLE leaderboardTest(username text, score int)")
	if err != nil {
		return err
	}
	_, err = testInstance.Query("INSERT INTO leaderboardTest (username, score) VALUES ('plietar', -14)")
	return err
}

func populate() error {
	_, err := testInstance.Query("INSERT INTO leaderboardTest (username, score) VALUES ('michaelRadigan', 100)")
	if err != nil {
		return err
	}
	_, err = testInstance.Query("INSERT INTO leaderboardTest (username, score) VALUES ('egnwd', -14)")
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	_, err = testInstance.Query("INSERT INTO leaderboardTest (username, score) VALUES ('lanejr', 3)")
	return err
}

func dropTable() error {
	_, err := testInstance.Query("DROP TABLE leaderboardTest")
	return err
}

func TestDatabase(t *testing.T) {
	err := setUpGeneric()
	defer dropTable()
	require.Nil(t, err, "error should be nil")
	rows, err := testInstance.Query("SELECT * FROM leaderboardTest")
	require.Nil(t, err, "error should be nil")
	var (
		username string
		score    int
	)
	rows.Next()
	defer rows.Close()
	err = rows.Scan(&username, &score)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, username, "plietar", "username is set to plietar in the database setup")
	assert.Equal(t, score, -14, "score is set to -14 in the database setup")

	err = rows.Err()
	assert.Nil(t, err, "error should be nil")
}

func TestGetMinScore(t *testing.T) {
	err := setUpGeneric()
	require.Nil(t, err, "error should be nil")
	defer dropTable()
	err = populate()
	require.Nil(t, err, "error should be nil")

	rows, err := testInstance.Query("SELECT MIN(score) FROM leaderboardTest")
	require.Nil(t, err, "error should be nil")
	defer rows.Close()
	rows.Next()

	var score int
	err = rows.Scan(&score)

	assert.Equal(t, score, -14, "-14 is the minimum score")
	err = rows.Err()
	assert.Nil(t, err, "error should be nil")
}

func TestDeleteSingle(t *testing.T) {
	err := setUpGeneric()
	require.Nil(t, err, "error should be nil")
	defer dropTable()
	err = populate()
	require.Nil(t, err, "error should be nil")

	deleteSingle := "DELETE FROM leaderboardTest WHERE ctid "
	deleteSingle += "IN (SELECT ctid FROM leaderboardTest ORDER BY "
	deleteSingle += "score asc LIMIT 1)"

	_, err = testInstance.Query(deleteSingle)
	require.Nil(t, err, "error should be nil")
	rows, err := testInstance.Query("SELECT MIN(score) FROM leaderboardTest")
	require.Nil(t, err, "error should be nil")
	rows.Next()
	var score int
	err = rows.Scan(&score)
	assert.Equal(t, score, -14, "-14 is the minimum score if correct value is deleted")
	rows.Close()

	// Count should be 3 (even though multiple entries with min)
	rows, err = testInstance.Query("select COUNT(*) from leaderboardTest")
	require.Nil(t, err, "error should be nil")
	rows.Next()
	var count int
	err = rows.Scan(&count)
	assert.Equal(t, count, 3, "There should now be 3 entries in the table")
	rows.Close()

	_, err = testInstance.Query(deleteSingle)
	require.Nil(t, err, "error should be nil")
	rows, err = testInstance.Query("SELECT MIN(score) FROM leaderboardTest")
	require.Nil(t, err, "error should be nil")
	rows.Next()
	err = rows.Scan(&score)
	assert.Equal(t, score, 3, "3 is the minimum score if correct 2 are deleted")
	rows.Close()

	_, err = testInstance.Query(deleteSingle)
	require.Nil(t, err, "error should be nil")
	rows, err = testInstance.Query("SELECT MIN(score) FROM leaderboardTest")
	require.Nil(t, err, "error should be nil")
	rows.Next()
	err = rows.Scan(&score)
	assert.Equal(t, score, 100, "100 is the minimum score if correct 3 are deleted")
	rows.Close()

	err = rows.Err()
	assert.Nil(t, err, "error should be nil")
}

func TestInsertNew(t *testing.T) {
	err := setUpGeneric()
	require.Nil(t, err, "error should be nil")
	defer dropTable()

	insertNew := "INSERT INTO leaderboardTest (username, score) "
	insertNew += "VALUES ('"
	insertNew += "michaelRadigan"
	insertNew += "', "
	insertNew += strconv.Itoa(9000)
	insertNew += ")"
	_, err = testInstance.Query(insertNew)
	require.Nil(t, err, "error should be nil")

	// Count should be 2
	rows, err := testInstance.Query("select COUNT(*) from leaderboardTest")
	require.Nil(t, err, "error should be nil")
	rows.Next()
	defer rows.Close()
	var count int
	err = rows.Scan(&count)
	assert.Equal(t, count, 2, "There should now be 2 entries in the table")
	rows.Close()

	err = rows.Err()
	assert.Nil(t, err, "error should be nil")
}

func TestReplace(t *testing.T) {
	err := setUpGeneric()
	require.Nil(t, err, "error should be nil")
	defer dropTable()

	deleteSingle := "DELETE FROM leaderboardTest WHERE ctid "
	deleteSingle += "IN (SELECT ctid FROM leaderboardTest ORDER BY "
	deleteSingle += "score asc LIMIT 1)"
	_, err = testInstance.Query(deleteSingle)
	require.Nil(t, err, "error should be nil")

	insertNew := "INSERT INTO leaderboardTest (username, score) "
	insertNew += "VALUES ('"
	insertNew += "michaelRadigan"
	insertNew += "', "
	insertNew += strconv.Itoa(9000)
	insertNew += ")"
	_, err = testInstance.Query(insertNew)
	require.Nil(t, err, "error should be nil")

	// Count should be 1
	rows, err := testInstance.Query("select COUNT(*) from leaderboardTest")
	require.Nil(t, err, "error should be nil")
	rows.Next()
	var count int
	err = rows.Scan(&count)
	assert.Equal(t, count, 1, "There should now be 1 entry in the table")
	rows.Close()

	rows, err = testInstance.Query("select * from leaderboardTest")
	require.Nil(t, err, "error should be nil")
	rows.Next()
	var (
		score int
		name  string
	)
	err = rows.Scan(&name, &score)
	require.Nil(t, err, "error should be nil")
	assert.Equal(t, name, "michaelRadigan", "The username should be \"michaelRadigan\"")
	assert.Equal(t, score, 9000, "The score should be 9000")
	rows.Close()

	err = rows.Err()
	assert.Nil(t, err, "error should be nil")
}
