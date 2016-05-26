package controller

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/egnwd/outgain/server/github"
	"github.com/gorilla/sessions"
	"github.com/nu7hatch/gouuid"
)

var store *sessions.CookieStore

const (
	stateKey       = "state"
	usernameKey    = "username"
	accessTokenKey = "access_token"
	sessionName    = "session"
)

func init() {
	authKey, _ := hex.DecodeString(os.Getenv("USER_STORE_SECRET_AUTH"))
	encKey, _ := hex.DecodeString(os.Getenv("USER_STORE_SECRET_ENC"))

	store = sessions.NewCookieStore(authKey, encKey)

	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 3600 * 8, // 8 hours
	}
}

// UserLogIn signs the user in and sets up a session
func UserLogIn(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := uuid.NewV4()
	state := id.String()

	session.Values[stateKey] = state
	log.Printf("Session: %v\n", session.Values)
	if err := sessions.Save(r, w); err != nil {
		log.Println(err.Error())
	}

	http.Redirect(w, r, github.GetOAuthURL(state), http.StatusFound)
}

// OAuthSignInCallback gets the access token from the Github API and uses it
// to get/create a user
func OAuthSignInCallback(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// Create user if they do not exist
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := r.FormValue("state")
	code := r.FormValue("code")

	log.Printf("Session: %v\n", session.Values)

	if state != session.Values[stateKey] {
		errorMessage := fmt.Sprintf("%d: Invalid state,\n\texpected: %s\n\tactual:%s",
			http.StatusUnauthorized, session.Values[stateKey], state)
		http.Error(w, errorMessage, http.StatusUnauthorized)
		return
	}

	accessToken, err := github.GetAccessToken(state, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username, err := github.GetUsername(accessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values[usernameKey] = username
	session.Values[accessTokenKey] = accessToken
	if err := sessions.Save(r, w); err != nil {
		log.Println(err.Error())
	}

	u := fmt.Sprintf("http://%s/", r.Host)

	http.Redirect(w, r, u, http.StatusFound)
}

// Logout deletes the user session
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	sessions.Save(r, w)

	u := fmt.Sprintf("http://%s/", r.Host)

	http.Redirect(w, r, u, http.StatusFound)
}

// CurrentUser returns the username of the session's user
func CurrentUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if username, ok := session.Values[usernameKey]; ok {
		fmt.Fprint(w, username)
	} else {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
	}
}
