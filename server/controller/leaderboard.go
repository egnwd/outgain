package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/egnwd/outgain/server/database"
)

// Get the top 10 global values from the database and returns them in a JSON format
func LeaderboardPeek(w http.ResponseWriter, r *http.Request) {
	leaderboard := database.GetAllRows()
	js, err := json.Marshal(leaderboard)
	database.NilCheck(err)
	fmt.Println(string(js))
	_ = w
	_ = r
}
