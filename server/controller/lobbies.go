package controller

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/egnwd/outgain/server/lobby"
	"github.com/egnwd/outgain/server/user"
)

func LobbiesView(w http.ResponseWriter, r *http.Request) {
	if !IsUserAuthorised(r) {
		u := fmt.Sprintf("http://%s/", r.Host)
		http.Redirect(w, r, u, http.StatusFound)
		return
	}

	lobby := lobby.NewLobby()
	id := lobby.ID

	form := url.Values{}
	form.Add("id", fmt.Sprintf("%d", id))

	u := fmt.Sprintf("http://%s/lobbies/join", r.Host)
	http.Post(u, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	u = fmt.Sprintf("http://%s/lobbies/%d", r.Host, id)
	http.Redirect(w, r, u, http.StatusFound)
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

	l, ok := lobby.GetLobby(id)
	if !ok {
		http.Error(w, "Lobby doesn't exist", http.StatusBadRequest)
		return
	}

	username, _ := GetUserName(r)
	user := user.NewUser(username)
	l.AddUser(user)

	log.Printf("User: %s Joined Lobby: %d", username, id)
}

func LobbiesGame(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check user is a part of this lobby
		if !IsUserAuthorised(r) {
			u := fmt.Sprintf("http://%s/", r.Host)
			http.Redirect(w, r, u, http.StatusFound)
			return
		}

		http.ServeFile(w, r, staticDir+"/game-view.html")
	})
}
