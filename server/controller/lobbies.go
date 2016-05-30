package controller

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func LobbiesView(w http.ResponseWriter, r *http.Request) {
	if !IsUserAuthorised(r) {
		u := fmt.Sprintf("http://%s/", r.Host)
		http.Redirect(w, r, u, http.StatusFound)
		return
	}

	id := strconv.FormatUint(uint64(rand.Uint32()), 10)
	form := url.Values{}
	form.Add("id", id)

	u := fmt.Sprintf("http://%s/lobbies/join", r.Host)

	http.Post(u, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	u = fmt.Sprintf("http://%s/lobbies/%s", r.Host, id)
	http.Redirect(w, r, u, http.StatusFound)
}

func LobbiesJoin(w http.ResponseWriter, r *http.Request) {
	if !IsUserAuthorised(r) {
		http.Error(w, "Not logged in.", http.StatusUnauthorized)
	}

	username, _ := GetUserName(r)
	id := r.PostFormValue("id")
	log.Printf("User: %s Joined Lobby: %s", username, id)
}

func LobbiesGame(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsUserAuthorised(r) {
			u := fmt.Sprintf("http://%s/", r.Host)
			http.Redirect(w, r, u, http.StatusFound)
			return
		}

		http.ServeFile(w, r, staticDir+"/game-view.html")
	})
}
