package controller

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/egnwd/outgain/server/lobby"
	"github.com/gorilla/mux"
)

func GetAISource() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		lobbyId, _ := strconv.ParseUint(vars["id"], 10, 64)

		username, err := GetUserName(r)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "No username", http.StatusUnauthorized)
			return
		}

		lobby, ok := lobby.GetLobby(uint64(lobbyId))
		if !ok {
			http.Error(w, "Lobby doesn't exist", http.StatusNotFound)
			return
		}

		user := lobby.FindGuest(username)
		if user == nil {
			http.Error(w, "User not in lobby", http.StatusNotFound)
			return
		}

		io.WriteString(w, user.Source)
	})
}

func PostAISource() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		lobbyId, _ := strconv.ParseUint(vars["id"], 10, 64)

		username, err := GetUserName(r)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "No username", http.StatusUnauthorized)
			return
		}

		lobby, ok := lobby.GetLobby(uint64(lobbyId))
		if !ok {
			http.Error(w, "Lobby doesn't exist", http.StatusNotFound)
			return
		}

		user := lobby.FindGuest(username)
		if user == nil {
			http.Error(w, "User not in lobby", http.StatusNotFound)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Read failed", http.StatusBadRequest)
		}

		user.Source = string(data)
	})
}
