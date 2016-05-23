package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/egnwd/outgain/server/github"
	"github.com/gorilla/sessions"
	"github.com/nu7hatch/gouuid"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("USER_STORE_SECRET_AUTH")),
	[]byte(os.Getenv("USER_STORE_SECRET_ENC")))

const (
	stateKey    = "state"
	sessionName = "session"
)

// UserSignIn signs the user in and sets up a session
func UserSignIn(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sesh")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := uuid.NewV4()
	state := id.String()

	session.Values[stateKey] = state
	fmt.Printf("Session: %#v\n", session.Values)
	sessions.Save(r, w)

	http.Redirect(w, r, github.GetOAuthURL(state), http.StatusFound)
}

// OAuthSignInCallback gets the access token from the Github API and uses it
// to get/create a user
func OAuthSignInCallback(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// Get username
	// Create user if they do not exist
	// Redirect to main/lobby page

	session, err := store.Get(r, "sesh")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := r.FormValue("state")
	code := r.FormValue("code")

	fmt.Printf("Session: %#v\n", session.Values)

	// HACK: We should not set this it should already be in the cookie store
	session.Values[stateKey] = state

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
		fmt.Println("Error on username")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Username: %s", username)
	fmt.Fprintln(w, message)
}
