package controller

import (
  "encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/egnwd/outgain/server/lobby"
	"github.com/gorilla/mux"
)

func LobbiesView(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsUserAuthorised(r) {
			u := fmt.Sprintf("http://%s/", r.Host)
			http.Redirect(w, r, u, http.StatusFound)
			return
		}
		http.ServeFile(w, r, staticDir+"/lobbies.html")
	})
}

// LobbiesPeek peeks at the lobby IDs in use and returns them as a JSON
func LobbiesPeek(w http.ResponseWriter, r *http.Request) {
	if !IsUserAuthorised(r) {
		u := fmt.Sprintf("http://%s/", r.Host)
		http.Redirect(w, r, u, http.StatusFound)
		return
	}

  IDs := lobby.GetLobbyIDs()
  js, err := json.Marshal(IDs)
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

	id, err := strconv.ParseUint(r.PostFormValue("id"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		return
	}

	l := lobby.GenerateOneLobby()
	log.Printf("Lobby: %d\n", l.ID)

	// if !ok {
	// 	log.Printf("Join: No Lobby (%d)\n", id)
	// 	http.Error(w, "Lobby doesn't exist", http.StatusBadRequest)
	// 	return
	// }

	username, err := GetUserName(r)
	if err != nil {
		log.Println(err.Error())
		return
	}
	user := lobby.NewUser(username)
	l.AddUser(user)

	log.Printf("User: %s Joined Lobby: %d", username, id)

	rawurl := fmt.Sprintf("http://%s/lobbies/%d", r.Host, id)
	http.Redirect(w, r, rawurl, http.StatusFound)
}

func LobbiesGame(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		l, ok := lobby.GetLobby(id)
		username, err := GetUserName(r)

		if err != nil {
			u := fmt.Sprintf("http://%s/", r.Host)
			http.Redirect(w, r, u, http.StatusFound)
			return
		} else if !ok || !l.ContainsUser(username) {
			u := fmt.Sprintf("http://%s/lobbies", r.Host)
			http.Redirect(w, r, u, http.StatusFound)
			return
		}

		http.ServeFile(w, r, staticDir+"/game-view.html")
	})
}
