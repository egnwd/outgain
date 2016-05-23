package routes

import (
	"net/http"

	"github.com/egnwd/outgain/server/controller"
	"github.com/gorilla/mux"
)

//GetHandler returns a mux that mapps routes to controller actions
func GetHandler(static string) http.Handler {
	mux := mux.NewRouter()

	get := mux.Methods(http.MethodGet).Subrouter()
	// post := mux.Methods(http.MethodPost).Subrouter()

	get.HandleFunc("/ping", controller.PingHandler)
	get.HandleFunc("/signin", controller.UserSignIn)
	get.HandleFunc("/oauthSignInCallback", controller.OAuthSignInCallback)
	get.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir(static))))

	return mux
}
