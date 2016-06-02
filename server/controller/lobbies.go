package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/egnwd/outgain/server/guest"
	"github.com/egnwd/outgain/server/lobby"
	"github.com/gorilla/mux"
)

func LobbiesView(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/lobbies.html")
	})
}

// LobbiesPeek peeks at the lobby IDs in use and returns them as a JSON
func LobbiesPeek(w http.ResponseWriter, r *http.Request) {
	// Get IDs of all current lobbies, convert to JSON and return it
	IDs := lobby.GetLobbyIDs()
	js, err := json.Marshal(IDs)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// LobbiesGetUsers gets all the user names from the lobby specified at the end
// of the caught URL, and returns them as a JSON
func LobbiesGetUsers(w http.ResponseWriter, r *http.Request) {
	// Get lobby ID from URL
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)

	l, ok := lobby.GetLobby(uint64(id))
	if !ok {
		// TODO: lobby no longer exists, perhaps refresh page and error popup
		return
	}
	// Get all usernames from lobby
	users := l.Guests.Iterator()
	usernames := make([]string, 0, len(users))
	for _, user := range users {
		usernames = append(usernames, user.GetName())
	}
	// Convert to JSON and return it
	js, err := json.Marshal(usernames)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func LobbiesJoin(w http.ResponseWriter, r *http.Request) {
	if !IsUserAuthorised(r) {
		http.Error(w, "Not logged in.", http.StatusUnauthorized)
	}

	// Get the id of the requested lobby
	id, err := strconv.ParseUint(r.PostFormValue("id"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Get the lobby from the global list
	l, ok := lobby.GetLobby(id)
	if !ok {
		log.Printf("Join: No Lobby (%d)\n", id)
		http.Error(w, "Lobby doesn't exist", http.StatusBadRequest)
		return
	}

	// Get the username of the authenicated user
	username, err := GetUserName(r)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Add the user to the lobby
	user := guest.NewUser(username)
	l.AddUser(user)

	l.Start()

	// Redirect user to the lobby
	log.Printf("User: %s Joined Lobby: %d", username, id)
	rawurl := fmt.Sprintf("http://%s/lobbies/%d", r.Host, id)
	http.Redirect(w, r, rawurl, http.StatusFound)
}

func LobbiesGame(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		l, ok := lobby.GetLobby(id)
		username, _ := GetUserName(r)

		if !ok || !l.ContainsUser(username) {
			u := fmt.Sprintf("http://%s/lobbies", r.Host)
			http.Redirect(w, r, u, http.StatusFound)
			return
		}

		http.ServeFile(w, r, staticDir+"/game-view.html")
	})
}

// LobbiesLeave temporarily logs the user out - this will change in the future
func LobbiesLeave(w http.ResponseWriter, r *http.Request) {
	if !IsUserAuthorised(r) {
		http.Error(w, "Not logged in.", http.StatusUnauthorized)
	}

	// Get the id of the requested lobby
	id, err := strconv.ParseUint(r.PostFormValue("id"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Get the lobby from the global list
	l, ok := lobby.GetLobby(id)
	if !ok {
		log.Printf("Join: No Lobby (%d)\n", id)
		http.Error(w, "Lobby doesn't exist", http.StatusBadRequest)
		return
	}

	// Get the username of the authenicated user
	username, err := GetUserName(r)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Remove the user to the lobby
	user := guest.NewUser(username)
	l.RemoveUser(user)

	// Redirect user to the lobby
	log.Printf("User: %s Left Lobby: %d", username, id)
	http.Redirect(w, r, "/", http.StatusFound)
}
