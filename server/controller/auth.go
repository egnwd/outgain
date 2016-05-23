package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/oauth2/github"
)

type config struct {
	ClientID, ClientSecret, RedirectURI string
}

var c config
var store = sessions.NewCookieStore([]byte(os.Getenv("USER_STORE_SECRET")))

const (
	githubScope = "user"
	sessionName = "session-login-user"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	c.ClientID = os.Getenv("CLIENT_ID")
	c.ClientSecret = os.Getenv("CLIENT_SECRET")
	c.RedirectURI = os.Getenv("REDIRECT_URI")
}

// UserSignIn signs the user in and sets up a session
func UserSignIn(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := uuid.NewV4()
	state := id.String()

	values := make(url.Values)
	values.Add("client_id", c.ClientID)
	values.Add("redirect_uri", c.RedirectURI)
	values.Add("scope", githubScope)
	values.Add("state", state)

	session.Values["State"] = state
	sessions.Save(r, w)

	githubOAuth := fmt.Sprintf("%s?%s", github.Endpoint.AuthURL, values.Encode())
	http.Redirect(w, r, githubOAuth, http.StatusFound)
}

// OAuthSignInCallback gets the access token from the Github API and uses it
// to get/create a user
func OAuthSignInCallback(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// Get session
	// Make Github API
	// Get Access Token
	// Get username
	// Create user if they do not exist
	// Redirect to main/lobby page

	// state := r.FormValue("state")
	code := r.FormValue("code")

	fmt.Fprintln(w, code)
}
