package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/egnwd/outgain/server/database"
)

// Get the top 10 global values from the database and returns them in a JSON format
func LeaderboardPeek(w http.ResponseWriter, r *http.Request) {
	leaderboard := database.GetAllRows()
	bs, err := json.Marshal(leaderboard)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func Leaderboard(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/leaderboard.html")
	})
}
