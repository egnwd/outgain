package routes

import (
	"net/http"
	"os"

	"github.com/egnwd/outgain/server/controller"
	"github.com/gorilla/handlers"
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

	// FIXME: Wrap the FileServer in a Handler that hooks w upon writing
	// 404 to the Header
	mux.NotFoundHandler = http.HandlerFunc(controller.NotFound)

	get.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(http.Dir(static))))

	return handlers.LoggingHandler(os.Stdout, mux)
}
