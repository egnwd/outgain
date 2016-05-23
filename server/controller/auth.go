package controller

import "net/http"

// UserSignIn signs the user in and sets up a session
func UserSignIn(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// Create session
	// Create values with client ID etc. based on dev/production environment
	// Add the state to the session
	// Save the session
	// Redirect to Github
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
}
