package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/egnwd/outgain/server/lobby"
	"github.com/gorilla/mux"
)

func UpdatesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		l, ok := lobby.GetLobby(uint64(id))
		if !ok {
			log.Println("Updates: Lobby doesn't exist")
			http.Error(w, "Lobby doesn't exist", http.StatusInternalServerError)
			return
		}

		l.UpdateRound()

		l.Events.ServeHTTP(w, r)
	})
}
